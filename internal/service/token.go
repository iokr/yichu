package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/iokr/yichu/internal/model"
)

// Token 类型。
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// Claims JWT 载荷。
type Claims struct {
	UserID   uint64 `json:"uid"`
	Username string `json:"username"`
	Type     string `json:"typ"`
	jwt.RegisteredClaims
}

// TokenManager 负责签发与校验 JWT。
type TokenManager struct {
	secret     []byte
	ttl        time.Duration
	refreshTTL time.Duration
}

// NewTokenManager 构造管理器，刷新 Token 固定 30 天。
func NewTokenManager(secret string, ttlHours int) *TokenManager {
	if ttlHours <= 0 {
		ttlHours = 168
	}
	return &TokenManager{
		secret:     []byte(secret),
		ttl:        time.Duration(ttlHours) * time.Hour,
		refreshTTL: 30 * 24 * time.Hour,
	}
}

// TTL 访问 Token 有效期。
func (m *TokenManager) TTL() time.Duration { return m.ttl }

func (m *TokenManager) sign(u *model.User, typ string, dur time.Duration) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(dur)
	claims := &Claims{
		UserID:   u.ID,
		Username: u.Username,
		Type:     typ,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "yichu",
			Subject:   fmt.Sprintf("%d", u.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	return signed, exp, err
}

// Issue 同时签发访问 Token 与刷新 Token。
func (m *TokenManager) Issue(u *model.User) (access string, refresh string, expiresAt time.Time, err error) {
	access, expiresAt, err = m.sign(u, TokenTypeAccess, m.ttl)
	if err != nil {
		return "", "", time.Time{}, err
	}
	refresh, _, err = m.sign(u, TokenTypeRefresh, m.refreshTTL)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return access, refresh, expiresAt, nil
}

// Parse 校验 Token 并校验类型。
func (m *TokenManager) Parse(token, expectType string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("签名算法不匹配")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("无效的 Token")
	}
	if expectType != "" && claims.Type != expectType {
		return nil, errors.New("Token 类型不匹配")
	}
	return claims, nil
}
