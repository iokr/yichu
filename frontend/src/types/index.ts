export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface DictOption {
  code: string
  label: string
  hex?: string
}

export interface Meta {
  name: string
  version: string
  buildTime: string
  hidePrice: boolean
}

export interface Dicts {
  colors: DictOption[]
  seasons: DictOption[]
  occasions: DictOption[]
  statuses: DictOption[]
  defaults: { status: string }
}

export interface UserInfo {
  id: number
  username: string
  nickname: string
}

export interface TokenResponse {
  token: string
  refresh_token: string
  expires_at: string
  expires_in: number
  user: UserInfo
}

export interface ItemImage {
  id: number
  url: string
  thumb_url: string
  sort: number
}

export interface TagItem {
  id: number
  name: string
  count?: number
}

export interface ItemListItem {
  id: number
  name: string
  category_id: number
  category: string
  sub_category_id: number | null
  sub_category: string
  brand: string
  color: DictOption[]
  size: string
  season: DictOption[]
  status: DictOption
  purchase_date: string | null
  purchase_price: number | null
  purchase_price_display: string
  cover: string
  wear_count: number
  last_worn_at: string | null
  price_masked: boolean
  created_at: string
}

export interface ItemDetail extends ItemListItem {
  occasion: DictOption[]
  purchase_channel: string
  purchase_url: string
  storage_location: string
  remark: string
  images: ItemImage[]
  tags: TagItem[]
  cost_per_wear: number | null
  cost_per_wear_display: string
  updated_at: string
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  size: number
}

export interface CategoryNode {
  id: number
  parent_id: number
  name: string
  sort: number
  is_system: boolean
  item_count: number
  children?: CategoryNode[]
}

export interface KeyStat {
  key: string
  label: string
  count: number
  percent: number
  amount: number | null
  amount_display: string
}

export interface CodeStat {
  code: string
  label: string
  hex?: string
  count: number
  percent: number
}

export interface DashboardData {
  total_count: number
  wearing_count: number
  idle_count: number
  total_wear_count: number
  category_stats: KeyStat[]
  status_stats: CodeStat[]
  year: number
  year_amount: number | null
  year_amount_display: string
  total_amount: number | null
  total_amount_display: string
  price_masked: boolean
  recent: ItemListItem[]
}

export interface StatsOverview {
  total_count: number
  wearing_count: number
  idle_count: number
  total_wear_count: number
  total_amount: number | null
  total_amount_display: string
  avg_price: number | null
  avg_price_display: string
  price_masked: boolean
}

export interface CostPerWearItem {
  id: number
  name: string
  cover: string
  wear_count: number
  purchase_price: number | null
  purchase_price_display: string
  cost_per_wear: number | null
  cost_per_wear_display: string
  price_masked: boolean
}

export interface WearLogItem {
  id: number
  item_id: number
  worn_date: string
  occasion: DictOption | null
  remark: string
  created_at: string
}

export interface SettingsData {
  hide_price: boolean
  price_toggle_require_password: boolean
  app_name: string
  version: string
  build_time: string
}

export interface UploadResult {
  url: string
  thumb_url: string
  width: number
  height: number
  size: number
}

export interface ItemImagePayload {
  url: string
  thumb_url: string
}

export interface ItemPayload {
  name: string
  category_id: number | null
  sub_category_id: number | null
  brand: string
  color: string[]
  size: string
  season: string[]
  occasion: string[]
  purchase_date: string
  purchase_price: number | null
  purchase_channel: string
  purchase_url: string
  status: string
  storage_location: string
  remark: string
  images: ItemImagePayload[]
  tags: string[]
}

export interface ItemQueryParams {
  category_id?: number
  sub_category_id?: number
  tag_id?: number
  status?: string
  keyword?: string
  brand?: string
  color?: string
  season?: string
  occasion?: string
  min_price?: number
  max_price?: number
  year?: number
  purchase_from?: string
  purchase_to?: string
  sort?: string
  order?: 'asc' | 'desc'
  page?: number
  size?: number
}
