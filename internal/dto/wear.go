package dto

import "github.com/iokr/yichu/internal/pkg/dict"

// WearLogRequest 穿着打卡请求。
type WearLogRequest struct {
	WornDate string `json:"worn_date"`
	Occasion string `json:"occasion"`
	Remark   string `json:"remark"`
}

// WearLogItem 穿着记录。
type WearLogItem struct {
	ID        uint64       `json:"id"`
	ItemID    uint64       `json:"item_id"`
	WornDate  string       `json:"worn_date"`
	Occasion  *dict.Option `json:"occasion"`
	Remark    string       `json:"remark"`
	CreatedAt string       `json:"created_at"`
}
