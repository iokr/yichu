package dict

import "strings"

// Option 枚举项：颜色带色值，其余仅 code + label。
type Option struct {
	Code  string `json:"code"`
	Label string `json:"label"`
	Hex   string `json:"hex,omitempty"`
}

// 状态 code。
const (
	StatusWearing   = "wearing"
	StatusIdle      = "idle"
	StatusDonated   = "donated"
	StatusDiscarded = "discarded"
)

// Colors 颜色字典（展示中文 + 色点）。
var Colors = []Option{
	{Code: "white", Label: "白", Hex: "#FFFFFF"},
	{Code: "black", Label: "黑", Hex: "#1F1D1B"},
	{Code: "gray", Label: "灰", Hex: "#9A9A9A"},
	{Code: "beige", Label: "米/驼", Hex: "#D8C3A5"},
	{Code: "brown", Label: "棕", Hex: "#8B6B4A"},
	{Code: "red", Label: "红", Hex: "#C0504D"},
	{Code: "pink", Label: "粉", Hex: "#E6A8B8"},
	{Code: "orange", Label: "橙", Hex: "#D98E5A"},
	{Code: "yellow", Label: "黄", Hex: "#D9C46A"},
	{Code: "green", Label: "绿", Hex: "#7E9A7B"},
	{Code: "blue", Label: "蓝", Hex: "#7C93AE"},
	{Code: "purple", Label: "紫", Hex: "#9A8AA8"},
	{Code: "multi", Label: "多色", Hex: ""},
}

// Seasons 季节字典。
var Seasons = []Option{
	{Code: "spring", Label: "春"},
	{Code: "summer", Label: "夏"},
	{Code: "autumn", Label: "秋"},
	{Code: "winter", Label: "冬"},
	{Code: "all", Label: "四季"},
}

// Occasions 场合字典。
var Occasions = []Option{
	{Code: "daily", Label: "日常"},
	{Code: "work", Label: "通勤"},
	{Code: "formal", Label: "正式"},
	{Code: "sport", Label: "运动"},
	{Code: "home", Label: "居家"},
}

// Statuses 状态字典。
var Statuses = []Option{
	{Code: StatusWearing, Label: "在穿"},
	{Code: StatusIdle, Label: "闲置"},
	{Code: StatusDonated, Label: "已捐赠"},
	{Code: StatusDiscarded, Label: "已丢弃"},
}

// Dicts 字典接口响应体，见设计文档第 6 节。
type Dicts struct {
	Colors    []Option `json:"colors"`
	Seasons   []Option `json:"seasons"`
	Occasions []Option `json:"occasions"`
	Statuses  []Option `json:"statuses"`
	// Defaults 供前端初始化表单使用。
	Defaults Defaults `json:"defaults"`
}

// Defaults 默认值。
type Defaults struct {
	Status string `json:"status"`
}

// All 返回全部字典。
func All() Dicts {
	return Dicts{
		Colors:    Colors,
		Seasons:   Seasons,
		Occasions: Occasions,
		Statuses:  Statuses,
		Defaults:  Defaults{Status: StatusWearing},
	}
}

func buildSet(list []Option) map[string]Option {
	set := make(map[string]Option, len(list))
	for _, o := range list {
		set[o.Code] = o
	}
	return set
}

var (
	colorSet    = buildSet(Colors)
	seasonSet   = buildSet(Seasons)
	occasionSet = buildSet(Occasions)
	statusSet   = buildSet(Statuses)
)

// IsColor 校验颜色 code。
func IsColor(code string) bool { _, ok := colorSet[code]; return ok }

// IsSeason 校验季节 code。
func IsSeason(code string) bool { _, ok := seasonSet[code]; return ok }

// IsOccasion 校验场合 code。
func IsOccasion(code string) bool { _, ok := occasionSet[code]; return ok }

// IsStatus 校验状态 code。
func IsStatus(code string) bool { _, ok := statusSet[code]; return ok }

// FindColor 查找颜色。
func FindColor(code string) (Option, bool) { o, ok := colorSet[code]; return o, ok }

// FindSeason 查找季节。
func FindSeason(code string) (Option, bool) { o, ok := seasonSet[code]; return o, ok }

// FindOccasion 查找场合。
func FindOccasion(code string) (Option, bool) { o, ok := occasionSet[code]; return o, ok }

// FindStatus 查找状态。
func FindStatus(code string) (Option, bool) { o, ok := statusSet[code]; return o, ok }

// StatusLabel 返回状态中文名；未知 code 原样返回。
func StatusLabel(code string) string {
	if o, ok := statusSet[code]; ok {
		return o.Label
	}
	return code
}

// Normalize 去空白、转小写、去重并保持顺序。
func Normalize(codes []string) []string {
	if len(codes) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(codes))
	out := make([]string, 0, len(codes))
	for _, raw := range codes {
		code := strings.ToLower(strings.TrimSpace(raw))
		if code == "" {
			continue
		}
		if _, dup := seen[code]; dup {
			continue
		}
		seen[code] = struct{}{}
		out = append(out, code)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// Validate 校验 code 列表，返回第一个非法 code（空字符串表示全部合法）。
func Validate(codes []string, valid func(string) bool) string {
	for _, code := range codes {
		if !valid(code) {
			return code
		}
	}
	return ""
}

// Options 把 code 列表转换成字典项列表，忽略非法 code。
func Options(codes []string, valid func(string) bool, find func(string) (Option, bool)) []Option {
	if len(codes) == 0 {
		return []Option{}
	}
	out := make([]Option, 0, len(codes))
	for _, code := range codes {
		if !valid(code) {
			continue
		}
		if o, ok := find(code); ok {
			out = append(out, o)
		}
	}
	return out
}

// ColorOptions 颜色 code → 字典项。
func ColorOptions(codes []string) []Option {
	return Options(codes, IsColor, FindColor)
}

// SeasonOptions 季节 code → 字典项。
func SeasonOptions(codes []string) []Option {
	return Options(codes, IsSeason, FindSeason)
}

// OccasionOptions 场合 code → 字典项。
func OccasionOptions(codes []string) []Option {
	return Options(codes, IsOccasion, FindOccasion)
}

// StatusOption 状态 code → 字典项（未知 code 也返回可展示对象）。
func StatusOption(code string) Option {
	if o, ok := FindStatus(code); ok {
		return o
	}
	return Option{Code: code, Label: code}
}
