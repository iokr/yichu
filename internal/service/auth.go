package service

import (
	"context"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/iokr/yichu/internal/dto"
	"github.com/iokr/yichu/internal/model"
	"github.com/iokr/yichu/internal/pkg/errcode"
	"github.com/iokr/yichu/internal/repository"
)

// AuthService 登录与 Token 续期。
type AuthService struct {
	users  *repository.UserRepository
	tokens *TokenManager
}

// NewAuthService 构造服务。
func NewAuthService(users *repository.UserRepository, tokens *TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

// Login 校验账密并签发 Token，失败提示统一。
func (s *AuthService) Login(_ context.Context, req dto.LoginRequest) (*dto.TokenResponse, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		return nil, errcode.Param("用户名或密码错误")
	}
	user, err := s.users.ByUsername(username)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	if user == nil {
		return nil, errcode.Param("用户名或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errcode.Param("用户名或密码错误")
	}
	return s.issue(user)
}

// Refresh 使用刷新 Token 换取新的 Token 对。
func (s *AuthService) Refresh(_ context.Context, refreshToken string) (*dto.TokenResponse, error) {
	claims, err := s.tokens.Parse(strings.TrimSpace(refreshToken), TokenTypeRefresh)
	if err != nil {
		return nil, errcode.Unauthorized("登录状态已失效，请重新登录")
	}
	user, err := s.users.ByID(claims.UserID)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	if user == nil {
		return nil, errcode.Unauthorized("登录状态已失效，请重新登录")
	}
	return s.issue(user)
}

func (s *AuthService) issue(user *model.User) (*dto.TokenResponse, error) {
	access, refresh, expiresAt, err := s.tokens.Issue(user)
	if err != nil {
		return nil, errcode.Internal(err)
	}
	return &dto.TokenResponse{
		Token:        access,
		RefreshToken: refresh,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
		ExpiresIn:    int64(time.Until(expiresAt).Seconds()),
		User: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
		},
	}, nil
}
