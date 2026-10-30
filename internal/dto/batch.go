package dto

// BatchRequest 批量删除/改状态。
type BatchRequest struct {
	IDs    []uint64 `json:"ids" binding:"required"`
	Action string   `json:"action" binding:"required"`
	Status string   `json:"status"`
}

// 批量动作。
const (
	BatchActionDelete = "delete"
	BatchActionStatus = "status"
)
