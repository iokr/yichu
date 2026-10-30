package dto

// MetaResponse 产品元信息，字段名与前端约定一致（camelCase）。
type MetaResponse struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	BuildTime  string `json:"buildTime"`
	HidePrice  bool   `json:"hidePrice"`
	PWAEnabled bool   `json:"pwaEnabled"`
}

// VersionResponse 版本信息。
type VersionResponse struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
}
