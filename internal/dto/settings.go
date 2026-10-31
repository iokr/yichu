package dto

// SettingsResponse 运行时设置。
type SettingsResponse struct {
	HidePrice                  bool   `json:"hide_price"`
	PriceToggleRequirePassword bool   `json:"price_toggle_require_password"`
	AppName                    string `json:"app_name"`
	Version                    string `json:"version"`
	BuildTime                  string `json:"build_time"`
}

// SettingsUpdateRequest 设置更新，只更新显式传入的字段。
type SettingsUpdateRequest struct {
	HidePrice *bool `json:"hide_price"`
}
