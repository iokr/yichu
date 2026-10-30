package dto

// CategoryNode 分类树节点。
type CategoryNode struct {
	ID        uint64         `json:"id"`
	ParentID  uint64         `json:"parent_id"`
	Name      string         `json:"name"`
	Sort      int            `json:"sort"`
	IsSystem  bool           `json:"is_system"`
	ItemCount int64          `json:"item_count"`
	Children  []CategoryNode `json:"children,omitempty"`
}

// CategoryRequest 分类创建/更新请求。
type CategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID uint64 `json:"parent_id"`
	Sort     int    `json:"sort"`
}

// CategorySortItem 排序项。
type CategorySortItem struct {
	ID   uint64 `json:"id"`
	Sort int    `json:"sort"`
}

// CategorySortRequest 分类排序请求。
type CategorySortRequest struct {
	Items []CategorySortItem `json:"items" binding:"required"`
}
