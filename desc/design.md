# 衣橱（yichu）· 最终设计文档

> 本文档为项目**唯一开发依据**，可直接用于代码生成。所有决策已确认，无待讨论项。

---

## 0. 项目标识

| 项 | 值 |
| --- | --- |
| 产品中文名 | **衣橱** |
| 代码服务名 | **`yichu`** |
| Go module | **`github.com/iokr/yichu`** |
| 入口目录 | `cmd/yichu/` |
| 可执行文件 | `yichu` / `yichu.exe` |
| 前端包名 | `yichu-web` |
| Docker 镜像 | `yichu/yichu:1.5.0` |
| Docker 容器 / Compose 服务名 | `yichu` |
| systemd 单元 | `yichu.service` |
| 环境变量前缀 | `YICHU_` |
| MySQL 数据库 / 用户 | `yichu` |
| 日志文件 | `yichu.log` |
| 部署目录 | `/opt/yichu` |
| 默认 `app.name` | `"衣橱"` |

**命名原则**：用户可见一律中文「衣橱」，代码与运维一律 `yichu`。

---

## 1. 技术栈

| 层 | 选型 |
| --- | --- |
| 后端语言 | Go **1.25.0** |
| Web 框架 | gin（最新版，`go get @latest`，保持 v1.x） |
| ORM | gorm（最新版，保持 v1.x） + `gorm.io/driver/mysql` |
| 数据库 | **MySQL 8**（utf8mb4 / utf8mb4_0900_ai_ci） |
| 鉴权 | golang-jwt v5 + bcrypt |
| 配置 | viper（`YICHU_` 前缀环境变量覆盖） |
| 图片处理 | disintegration/imaging |
| 前端 | **Vue 3.4** + Vite 5 + TypeScript |
| 状态 / 路由 | Pinia + Vue Router 4 |
| UI | TailwindCSS + Element Plus（有限使用） |
| 图表 | ECharts |
| PWA | vite-plugin-pwa |
| 部署 | Docker Compose（mysql + yichu + caddy），Caddy 自动 HTTPS |
| 仓库 | 单仓库 Monorepo，go.mod 在根 |

---

## 2. 产品范围

**定位**：家庭单用户衣物（衣服/鞋子/帽子）数字化管理，多端 Web 共用一套代码，多端同域。

**目标**：单件录入 ≤ 60 秒；任意衣物 3 次点击内定位；完整记录购买时间/价格/渠道/穿着次数；按分类/品牌/年份统计；`docker compose up -d` 一键启动；PWA 可安装可离线；价格服务端脱敏。

**非目标**：多用户共享、电商比价、AI 搭配、二手交易、社交分享、原生 App、微信小程序。

**角色**：单账号（管理员），首次启动自动初始化，无注册页。

---

## 3. 功能需求

| 模块 | 功能 | 优先级 |
| --- | --- | --- |
| 登录 | 账密登录、Token 续期、退出 | P0 |
| 仪表盘 | 总数、分类占比、年度消费、最近添加、快捷入口 | P0 |
| 衣物管理 | CRUD、图片、组合筛选、搜索、排序、分页、批量 | P0 |
| 分类管理 | 一级（衣服/鞋子/帽子）+ 二级维护 | P0 |
| 图片管理 | 上传、压缩、缩略图、预览、设封面 | P0 |
| 价格脱敏 | 服务端抹除 + 开关 | P0 |
| 枚举字典 | 颜色/季节/场合/状态 | P0 |
| 深色模式 | 跟随系统 + 手动切换 | P0 |
| PWA | 可安装、离线浏览、更新提示 | P0 |
| 产品名可配置 | 改配置重启即生效 | P0 |
| 统计报表 | 分类/品牌/年份/季节/颜色/单次成本 | P1 |
| 穿着记录 | 打卡、次数、最近穿着、单次成本、日历 | P1 |
| 标签 | 自定义标签、按标签筛选 | P1 |
| 导出 | CSV / Excel（受脱敏约束） | P1 |
| 存放位置 | 位置记录 | P1 |
| 状态流转 | 在穿/闲置/已捐赠/已丢弃 | P1 |
| 系统设置 | 主题、隐藏价格、备份、版本 | P1 |
| 心愿单 / 搭配套装 | — | P2 |

### 3.1 衣物字段

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| name | string(100) | ✅ | 名称 |
| category_id | uint64 | ✅ | 一级分类 |
| sub_category_id | uint64 | ❌ | 二级分类 |
| brand | string(50) | ❌ | 品牌 |
| color | JSON[] | ❌ | 颜色 code 多选 |
| size | string(20) | ❌ | 尺码 |
| season | JSON[] | ❌ | 季节 code 多选 |
| occasion | JSON[] | ❌ | 场合 code 多选 |
| purchase_date | date | ❌ | 购买时间（核心） |
| purchase_price | decimal(10,2) | ❌ | 价格（受脱敏） |
| purchase_channel | string(50) | ❌ | 渠道 |
| purchase_url | string(500) | ❌ | 订单链接 |
| status | string(16) | ✅ | 状态 code，默认 `wearing` |
| storage_location | string(50) | ❌ | 存放位置 |
| images | array | ❌ | 1~5 张，首张封面 |
| remark | string(500) | ❌ | 备注 |
| tags | array | ❌ | 标签 |
| wear_count / last_worn_at | 派生 | 自动 | 穿着统计 |

### 3.2 分类种子数据

| 一级 | 二级 |
| --- | --- |
| 衣服 | 上衣、T恤、衬衫、卫衣、毛衣、外套、羽绒服、连衣裙、半身裙、裤子、牛仔裤、内衣、睡衣 |
| 鞋子 | 运动鞋、板鞋、帆布鞋、高跟鞋、平底鞋、靴子、凉鞋、拖鞋 |
| 帽子 | 棒球帽、渔夫帽、贝雷帽、毛线帽、遮阳帽 |

一级分类不可删可改名；二级可增删改排序，有关联衣物时禁止删除。

### 3.3 图片规则

- jpg/png/webp，单张 ≤ 10MB，单件 ≤ 5 张
- 压缩：长边 ≤ 1600px、质量 80%；缩略图 300px
- 存储 `uploads/original` + `uploads/thumb`，文件名服务端 UUID 生成
- 预留 OSS/MinIO 适配接口

### 3.4 登录规则

- bcrypt 校验；JWT 默认 7 天 + 刷新 Token
- 失败统一提示「用户名或密码错误」
- Token 存 localStorage，Axios 拦截器注入 `Authorization: Bearer <token>`
- 401 跳登录页并清空本地状态
- `/api/v1/auth/*` 与 `/api/v1/settings` 响应 `Cache-Control: no-store`

---

## 4. 多端响应式

| 端 | 断点 | 定位 |
| --- | --- | --- |
| 手机 | `< 768px` | 录入主场景（拍照） |
| 平板 | `768 ~ 1199px` | 浏览为主，兼顾录入 |
| 电脑 | `≥ 1200px` | 批量管理、统计 |

实现：Tailwind 断点 + `useBreakpoint()`，禁止 UA 判断。多端同域。

| 区域 | 手机 | 平板 | 电脑 |
| --- | --- | --- | --- |
| 导航 | 底部 Tab（首页/衣物/统计/我的） | 左侧图标栏可展开 | 左侧侧边栏可折叠 |
| 顶栏 | 标题 + 搜索图标 | 标题 + 搜索框 | 标题 + 面包屑 + 全局搜索 |
| 内容 | 单列边距 12px | 2 列 / 侧栏+内容 | 侧栏 + 主区，最大宽 1440 居中 |
| 新增 | 右下 FAB | 悬浮 FAB | 右上主按钮 |
| 弹窗 | 底部抽屉/全屏 | 居中弹窗 | 居中弹窗 |
| 筛选 | 全屏抽屉 | 顶部弹层 | 左侧分类树 + 顶部筛选条 |
| 表格 | 不用 | 简化列表 | 完整表格 + 列配置 |

**页面适配**：登录（手机全屏居中 / 平板居中卡片 420 / 电脑左图右表单）；仪表盘（2×2 卡 / 4 列 / 4 列 + 3 列图表）；列表（1 列大卡 / 3 列 / 4~6 列或表格）；新增编辑（手机全屏底部固定保存 / 单列 / 左右两栏）；详情（上图下文 / 双栏 / 左图 40% 右信息 60%）；统计（单列 / 2 列 / 2~3 列）；分类管理（列表+抽屉 / 双栏 / 双栏）。

**移动端专项**：拍照上传 `accept="image/*" capture="environment"`；可点击元素 ≥ 44px；`inputmode` 正确且聚焦自动上滚；安全区 `env(safe-area-inset-*)`；图片左右滑、下拉刷新、列表左滑删除；PWA 全屏打开；上传进度可取消可重试；图片懒加载；横屏不错乱。

**电脑端专项**：表格视图列宽拖拽/列显隐/批量操作；快捷键 `/` 搜索、`N` 新增、`Esc` 关闭；多图拖拽 + 粘贴板上传；侧边栏折叠状态本地记忆。

---

## 5. 视觉设计系统

**关键词**：温柔 · 克制 · 留白 · 高级感。参考 Apple 留白 + 无印良品克制 + Notion 轻量组件。禁止：玻璃拟态、拟物化、大面积色块、重投影、霓虹色、花哨插画、过多分割线。

### 5.1 色彩

**亮色**

| Token | 值 | 用途 |
| --- | --- | --- |
| `--color-primary` | `#A67C89` | 莫兰迪玫瑰，主按钮/选中态 |
| `--color-primary-hover` | `#94697A` | 悬停 |
| `--color-primary-soft` | `#F5EDEF` | 主色浅底 |
| `--color-accent` | `#7E9A93` | 雾青，图表第二色 |
| `--color-bg` | `#FAFAF8` | 页面背景 |
| `--color-surface` | `#FFFFFF` | 卡片/弹窗/侧栏 |
| `--color-surface-2` | `#F4F3F1` | 次级容器/输入底 |
| `--color-border` | `#EAE8E5` | 唯一分割线 |
| `--color-text` | `#1F1D1B` | 主文字 |
| `--color-text-2` | `#6B6661` | 次要文字 |
| `--color-text-3` | `#A6A19B` | 占位/禁用 |
| success / warning / danger / info | `#6E9C7C` / `#D1A055` / `#C06C6C` / `#7C93AE` | 语义 |

**暗色**

| Token | 值 |
| --- | --- |
| primary | `#C79AA8` |
| bg | `#141312` |
| surface | `#1C1B1A` |
| surface-2 | `#252322` |
| border | `#2E2C2A` |
| text / 2 / 3 | `#F2F0EE` / `#A8A39D` / `#6E6963` |

跟随 `prefers-color-scheme` + 手动切换（设置页），选择存 localStorage。

**约束**：单页彩色面积 ≤ 10%；主色不用于大面积背景；分割线只用 `--color-border`；图表固定 6 色序列 `#A67C89 / #7E9A93 / #C4A484 / #8D9DB6 / #B5A8A0 / #9AA88E`。

### 5.2 排版

| 层级 | 字号/行高 | 字重 |
| --- | --- | --- |
| Display | 32/40 | 600 |
| H1 | 24/32 | 600 |
| H2 | 18/26 | 600 |
| Body | 14/22 | 400 |
| Caption | 12/18 | 400 |

字体栈：`"Inter", "PingFang SC", "HarmonyOS Sans SC", "Microsoft YaHei", system-ui`。数字 `tabular-nums`；中文标题 `letter-spacing: -0.01em`；段落最大宽 68ch。

### 5.3 间距 / 圆角 / 阴影

```css
--space-1:4px; --space-2:8px; --space-3:12px; --space-4:16px;
--space-5:24px; --space-6:32px; --space-7:48px; --space-8:64px;
--radius-sm:8px; --radius-md:12px; --radius-lg:16px; --radius-xl:24px; --radius-full:999px;
--shadow-sm:0 1px 2px rgba(31,29,27,.04);
--shadow-md:0 4px 16px rgba(31,29,27,.06);
--shadow-lg:0 12px 32px rgba(31,29,27,.10);
```

卡片默认**无阴影 + 1px 边框**；仅 hover、弹窗、FAB 加阴影。

### 5.4 组件规范

| 组件 | 规范 |
| --- | --- |
| 主按钮 | 高 40（手机 44），radius-md，主色填充白字，hover 加深，active 缩放 0.98 |
| 次按钮 | 白底 + 1px 边框，文字主色 |
| 输入框 | 高 40，底 surface-2 无边框，focus 2px 主色描边 + 浅外发光 |
| 卡片 | 白底 radius-lg，1px border，内边距 16，hover 上移 2px + shadow-sm |
| Chip | 胶囊高 28，浅底，选中主色底白字（颜色/季节/场合多选） |
| 图片容器 | 3:4，cover，radius-lg，加载前 surface-2 占位 |
| 空状态 | 居中线性插画（1.5px 描边，主色 40% 透明）+ 温柔文案 + 1 主按钮 |
| Toast | 顶部居中白底 shadow-md，左侧 4px 语义色竖条 |
| 弹窗 | 白底 radius-xl 无边框 shadow-lg，遮罩 `rgba(31,29,27,.32)` + 4px 模糊 |
| 骨架屏 | surface-2 + 1.2s 微光，形状与真实内容一致 |
| Money | 全站金额唯一渲染组件，脱敏态显示 `¥ •••` |

图标统一 Lucide，`stroke-width:1.5`，尺寸 16/20/24。

### 5.5 动效

| 场景 | 时长 | 曲线 |
| --- | --- | --- |
| hover/按钮 | 150ms | `cubic-bezier(.4,0,.2,1)` |
| 弹窗/抽屉 | 220ms | `cubic-bezier(.32,.72,0,1)` |
| 页面切换 | 240ms | opacity + translateY(8px) |
| 列表加载 | 逐项 30ms 递增淡入，最多 8 项 |
| 数字变化 | 600ms 滚动计数 |

只动 transform/opacity；尊重 `prefers-reduced-motion`；禁弹跳/旋转/闪烁。

### 5.6 页面要点

- **登录**：电脑端左 45% 主色 6% 透明色块 + 居中 logo 文案，右侧纯白无边框表单；手机全屏白底顶部留白 64px
- **仪表盘**：4 数据卡，数字大标签小无边框；图表卡白底，去网格线仅留极浅横向参考线
- **列表**：卡片网格，图片占 75% 高，下方「名称 + 分类·日期」，价格小字右对齐；筛选胶囊 Chip 横滑；颜色以色点展示
- **详情**：左圆角大图 + 缩略条；右信息「字段名 text-3 小字 + 值 text 大字」两行一组，组距 20px，不用表格线；底部操作栏白底 + 顶部 1px 边框
- **新增/编辑**：分区（基本信息/购买信息/收纳信息），区标题 H2 + 上方 32px；图片区为首屏最大元素；底部满宽固定保存
- **统计**：图表容器高 280，卡内标题 + 右上胶囊切换；数值直接标注图形上，去图例
- **分类管理**：双栏，左树选中行浅主色底 + 左侧 3px 主色竖条，右编辑区大留白

### 5.7 技术落地

```css
/* frontend/src/styles/tokens.css */
:root{
  --color-primary:#A67C89; --color-primary-hover:#94697A; --color-primary-soft:#F5EDEF;
  --color-accent:#7E9A93; --color-bg:#FAFAF8; --color-surface:#FFFFFF;
  --color-surface-2:#F4F3F1; --color-border:#EAE8E5;
  --color-text:#1F1D1B; --color-text-2:#6B6661; --color-text-3:#A6A19B;
  --radius-md:12px; --radius-lg:16px; --space-4:16px;
  --shadow-md:0 4px 16px rgba(31,29,27,.06);
  --font-sans:"Inter","PingFang SC","Microsoft YaHei",system-ui,sans-serif;
}
@media (prefers-color-scheme: dark){ :root{ /* 见 5.1 暗色 */ } }
```

```ts
// frontend/tailwind.config.ts
export default {
  darkMode:'class',
  theme:{ extend:{
    colors:{ primary:{DEFAULT:'#A67C89',hover:'#94697A',soft:'#F5EDEF'},
      accent:'#7E9A93', bg:'#FAFAF8', surface:{DEFAULT:'#FFFFFF',2:'#F4F3F1'},
      line:'#EAE8E5', ink:{DEFAULT:'#1F1D1B',2:'#6B6661',3:'#A6A19B'} },
    borderRadius:{ md:'12px', lg:'16px', xl:'24px' },
    boxShadow:{ sm:'0 1px 2px rgba(31,29,27,.04)',
      md:'0 4px 16px rgba(31,29,27,.06)', lg:'0 12px 32px rgba(31,29,27,.10)' },
    aspectRatio:{ photo:'3 / 4' },
  }},
}
```

Element Plus 仅用于 Form 校验、DatePicker、Dialog、Message；通过 `--el-color-primary` 等变量覆盖主题；优先用 Tailwind + 自研原子组件：`AppButton / AppCard / AppChip / AppField / AppEmpty / AppSkeleton / AppImage / AppPageHeader / Money`。

---

## 6. 枚举字典

**颜色**（存 code，展示中文 + 色点）

| code | label | hex |
| --- | --- | --- |
| white | 白 | #FFFFFF |
| black | 黑 | #1F1D1B |
| gray | 灰 | #9A9A9A |
| beige | 米/驼 | #D8C3A5 |
| brown | 棕 | #8B6B4A |
| red | 红 | #C0504D |
| pink | 粉 | #E6A8B8 |
| orange | 橙 | #D98E5A |
| yellow | 黄 | #D9C46A |
| green | 绿 | #7E9A7B |
| blue | 蓝 | #7C93AE |
| purple | 紫 | #9A8AA8 |
| multi | 多色 | 前端渐变渲染 |

**季节**：spring 春 / summer 夏 / autumn 秋 / winter 冬 / all 四季
**场合**：daily 日常 / work 通勤 / formal 正式 / sport 运动 / home 居家
**状态**：wearing 在穿 / idle 闲置 / donated 已捐赠 / discarded 已丢弃

- 存储：MySQL JSON 存 code 数组，如 `["white","beige"]`
- 字典接口免登录：`GET /api/v1/dicts` 返回上述四组（colors 含 hex）
- 后端 handler 对提交 code 白名单校验，非法返回 `1001`
- 前端 `useDict()`：Pinia 缓存 + localStorage 持久化 7 天

---

## 7. 价格脱敏

**原则**：脱敏在服务端完成，开关打开时接口响应**不出现真实金额**。

**配置与存储**：`app.hide_price`（默认 false）；运行时开关存 `settings` 表（key=`hide_price`），优先级 **settings 表 > config.yaml > 默认 false**；`PUT /api/v1/settings` body `{ "hide_price": true }`；预留 `app.price_toggle_require_password`（V1 默认关闭）。

**响应契约**

列表/详情未脱敏：`{ "purchase_price": 399.00, "purchase_price_display": "¥399.00", "price_masked": false }`
脱敏：`{ "purchase_price": null, "purchase_price_display": "¥ ***", "price_masked": true }`
统计脱敏：`{ "total_count":128, "total_amount":null, "total_amount_display":"***", "price_masked":true }`
导出脱敏：价格列 `***`，响应头 `X-Price-Masked: true`。

**实现**：

```go
// internal/service/mask.go
type PriceMasker struct{ hidden bool }
func (m PriceMasker) Price(v *float64) (*float64, string) {
    if v == nil { return nil, "" }
    if m.hidden { return nil, "¥ ***" }
    return v, fmt.Sprintf("¥%.2f", *v)
}
```

- DTO 层统一调用 masker 填充 `*_display`，**禁止** handler 直接返回 model
- 日志中间件不记录响应体；错误日志禁打印价格
- `price_masked` 前端用于渲染占位与禁用「性价比排序」

**前端**：全站金额只允许通过 `<Money>` 渲染（评审硬性检查）。脱敏态下性价比榜隐藏、年度消费柱状图改件数维度、导出旁提示「价格已隐藏」。切换开关后清空 PWA `api-data-*` 缓存（见第 10 节）。

---

## 8. 后端架构

### 8.1 目录

```
cmd/yichu/main.go
internal/
├── config/          # viper，YICHU_ 环境变量
├── database/        # MySQL 连接、连接池、重试、AutoMigrate、种子
├── router/
├── middleware/      # JWT、CORS、日志、Recover
├── handler/         # HTTP 层（含 meta.go）
├── service/         # 业务逻辑（含 mask.go）
├── repository/      # GORM 数据访问
├── model/
├── dto/
└── pkg/             # dict 字典、图片压缩、响应封装、错误码、IP
```

分层依赖：router → handler → service → repository → model。handler 不直接访问 repository。

### 8.2 配置

```yaml
# config.yaml
app:
  name: "衣橱"
  hide_price: false
server:
  addr: ":8080"
  mode: "release"            # debug | release
  tls:
    enabled: false
    cert_file: ""
    key_file: ""
database:
  driver: "mysql"
  dsn: "yichu:${DB_PASSWORD}@tcp(mysql:3306)/yichu?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s&readTimeout=30s&writeTimeout=30s"
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime_minutes: 30
  auto_migrate: true
auth:
  jwt_secret: "change-me-in-production"
  token_ttl_hours: 168
  init_username: "admin"
  init_password: "yichu123"
storage:
  upload_dir: "./uploads"
  max_size_mb: 10
  max_count: 5
log:
  level: "info"
  file: "./logs/yichu.log"
```

环境变量前缀 `YICHU_`，如 `YICHU_SERVER_ADDR=:9000`、`YICHU_DATABASE_DSN=...`。

### 8.3 数据库连接

- MySQL 8，`utf8mb4` / `utf8mb4_0900_ai_ci`，时区 `Asia/Shanghai`
- 启动 Ping 重试（最多 10 次，间隔 2s），适配 Compose 中 MySQL 启动慢
- 账号最小权限 `SELECT, INSERT, UPDATE, DELETE`
- V1.0 用 GORM `AutoMigrate` + `settings` 表记录 `schema_version`；V1.2 引入 `migrations/` 版本化 SQL
- 备份：`mysqldump --single-transaction` 每日 03:00，保留 14 天

### 8.4 错误码

`0` 成功、`1001` 参数错误、`1002` 未登录、`1003` 无权限、`1004` 资源不存在、`5000` 服务器错误。

---

## 9. 数据模型（MySQL）

```
users         (id, username, password_hash, nickname, created_at)
categories    (id, parent_id, name, sort, is_system, created_at)
items         (id, name, category_id, sub_category_id, brand, color(json), size,
               season(json), occasion(json), purchase_date, purchase_price,
               purchase_channel, purchase_url, status, storage_location,
               remark, wear_count, last_worn_at, created_at, updated_at, deleted_at)
item_images   (id, item_id, url, thumb_url, sort, created_at)
tags          (id, name, created_at)
item_tags     (item_id, tag_id)
wear_logs     (id, item_id, worn_date, occasion, remark, created_at)
settings      (key, value)
```

```sql
CREATE TABLE `items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `category_id` BIGINT UNSIGNED NOT NULL,
  `sub_category_id` BIGINT UNSIGNED DEFAULT NULL,
  `brand` VARCHAR(50) DEFAULT NULL,
  `color` JSON DEFAULT NULL,
  `size` VARCHAR(20) DEFAULT NULL,
  `season` JSON DEFAULT NULL,
  `occasion` JSON DEFAULT NULL,
  `purchase_date` DATE DEFAULT NULL,
  `purchase_price` DECIMAL(10,2) DEFAULT NULL,
  `purchase_channel` VARCHAR(50) DEFAULT NULL,
  `purchase_url` VARCHAR(500) DEFAULT NULL,
  `status` VARCHAR(16) NOT NULL DEFAULT 'wearing',
  `storage_location` VARCHAR(50) DEFAULT NULL,
  `remark` VARCHAR(500) DEFAULT NULL,
  `wear_count` INT UNSIGNED NOT NULL DEFAULT 0,
  `last_worn_at` DATETIME DEFAULT NULL,
  `created_at` DATETIME(3) NOT NULL,
  `updated_at` DATETIME(3) NOT NULL,
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category` (`category_id`),
  KEY `idx_sub_category` (`sub_category_id`),
  KEY `idx_purchase_date` (`purchase_date`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_name` (`name`(32))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

- 主键统一 `BIGINT UNSIGNED`（GORM `uint64`）
- 软删除 `deleted_at` 必须建索引
- JSON 字段 V1 用应用层 `JSON_CONTAINS` 筛选

---

## 10. PWA 离线

```ts
// frontend/vite.config.ts（节选）
VitePWA({
  registerType: 'prompt',
  manifest: false,                 // manifest 由后端动态生成
  workbox: {
    globPatterns: ['**/*.{js,css,html,svg,png,woff2}'],
    navigateFallback: '/index.html',
    navigateFallbackDenylist: [/^\/api\//, /^\/uploads\//, /^\/manifest\.webmanifest$/],
    runtimeCaching: [
      { urlPattern: /^\/api\/v1\/(dicts|meta|version)/,
        handler: 'StaleWhileRevalidate',
        options: { cacheName: 'api-static', expiration: { maxAgeSeconds: 604800 } } },
      { urlPattern: /^\/api\/v1\/(items|categories|dashboard|stats)(\?.*)?$/,
        handler: 'NetworkFirst',
        options: { cacheName: 'api-data-v1', networkTimeoutSeconds: 3,
          expiration: { maxEntries: 200, maxAgeSeconds: 86400 },
          cacheableResponse: { statuses: [200] } } },
      { urlPattern: /^\/uploads\//, handler: 'CacheFirst',
        options: { cacheName: 'images-v1',
          expiration: { maxEntries: 500, maxAgeSeconds: 2592000 } } },
    ],
  },
})
```

| 资源 | 策略 |
| --- | --- |
| app shell（js/css/html） | Precache |
| `/assets/*`（hash） | CacheFirst + 永久 |
| 字典 / meta / version | StaleWhileRevalidate |
| 列表/详情/统计 GET | NetworkFirst（3s 超时回退） |
| 图片 | CacheFirst + 500 条上限 |
| 写操作（POST/PUT/DELETE） | 不缓存 |
| `/api/*` 非 200 | 不缓存 |
| settings / auth | `no-store` |

**离线交互**：离线时顶部细条「当前离线，数据可能不是最新」；写操作 Toast「网络不可用，请联网后重试」，不做离线队列；`offline.html` 兜底。

**脱敏联动**：切换价格可见性时清空缓存：
```ts
const names = await caches.keys()
await Promise.all(names.filter(n => n.startsWith('api-data')).map(n => caches.delete(n)))
```

**更新提示**：`useRegisterSW` 的 `needRefresh` 为 true 时底部弹「衣橱有新版本，点击刷新」。

---

## 11. 前端嵌入与单文件部署

### 11.1 仓库结构

```
yichu/                              # 仓库根 = Go module 根
├── go.mod                          # module github.com/iokr/yichu; go 1.25.0
├── go.sum
├── cmd/yichu/main.go
├── internal/                       # 见 8.1
├── web/
│   ├── embed.go                    # //go:embed all:dist
│   └── dist/.gitkeep
├── frontend/                       # Vue 工程
├── migrations/                     # V1.2 起
├── scripts/backup.sh
├── docs/
├── Caddyfile
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── .env.example
├── .gitignore
└── README.md
```

**关键约束**：
1. go.mod 在仓库根，`go build ./cmd/yichu` 直接出二进制
2. `web/dist` 必须在 Go 包目录内（embed 不能 `../`），Vite 产物直出 `../web/dist`
3. 提交 `web/dist/.gitkeep` 保证未构建时可编译

### 11.2 embed 实现

```go
// web/embed.go
package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
)

//go:embed all:dist
var distFS embed.FS

func FS() (fs.FS, error) { return fs.Sub(distFS, "dist") }

const notBuiltHTML = `<!doctype html><html><head><meta charset="utf-8"><title>衣橱</title></head>
<body style="font-family:system-ui;padding:64px;color:#6B6661">
<h2 style="color:#1F1D1B">前端尚未构建</h2>
<p>请执行 <code>make front</code> 或 <code>cd frontend &amp;&amp; npm run build</code> 后重新编译。</p>
<p>API 已正常运行，可访问 <code>/api/v1/version</code>。</p></body></html>`

func Register(r *gin.Engine) error {
	sub, err := FS()
	if err != nil { return err }
	fileServer := http.FileServer(http.FS(sub))

	r.NoRoute(func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": 1004, "message": "resource not found"})
			return
		}
		if f, err := sub.Open(strings.TrimPrefix(p, "/")); err == nil {
			_ = f.Close()
			if strings.HasPrefix(p, "/assets/") {
				c.Header("Cache-Control", "public, max-age=31536000, immutable")
			}
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}
		c.Header("Cache-Control", "no-cache")
		index, err := fs.ReadFile(sub, "index.html")
		if err != nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(notBuiltHTML))
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
	return nil
}
```

### 11.3 main.go

```go
package main

import (
	"context"; "errors"; "log"; "net"; "net/http"
	"os"; "os/signal"; "syscall"; "time"
	"github.com/gin-gonic/gin"
	"github.com/iokr/yichu/internal/config"
	"github.com/iokr/yichu/internal/database"
	"github.com/iokr/yichu/internal/router"
	"github.com/iokr/yichu/web"
)

var ( Version = "dev"; BuildTime = "unknown" )

func main() {
	cfg := config.Load()
	if cfg.Server.Mode == "release" { gin.SetMode(gin.ReleaseMode) }

	db := database.Init(cfg)          // 连接 + 重试 + AutoMigrate + 种子 + 首次建账号
	r := router.New(cfg, db)
	if err := web.Register(r); err != nil { log.Fatalf("挂载前端失败: %v", err) }

	srv := &http.Server{Addr: cfg.Server.Addr, Handler: r}
	go func() {
		log.Printf("🌿 衣橱 %s (build %s) 已启动: http://%s", Version, BuildTime, cfg.Server.Addr)
		if ip := localIP(); ip != "" { log.Printf("📱 局域网访问: http://%s%s", ip, cfg.Server.Addr) }
		var err error
		if cfg.Server.TLS.Enabled {
			err = srv.ListenAndServeTLS(cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) { log.Fatalf("服务异常: %v", err) }
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("👋 衣橱已退出")
}

func localIP() string {
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}
```

---

## 12. 产品名可配置

改 `config.yaml` 的 `app.name` 并重启，无需重建前端，全站名称同步生效。

| 位置 | 实现 |
| --- | --- |
| Meta 接口 | `GET /api/v1/meta` → `{ name, version, buildTime, hidePrice, pwaEnabled }` |
| 动态 manifest | `GET /manifest.webmanifest` 由后端读配置生成 |
| 页面标题 | 前端启动拉 meta 后 `document.title = meta.name` |
| Logo 文字 | 布局组件读 `meta.name` |
| index.html | 兜底 `<title>衣橱</title>` + `<link rel="manifest" href="/manifest.webmanifest">` |

```go
func (h *MetaHandler) Manifest(c *gin.Context) {
    c.Header("Content-Type", "application/manifest+json; charset=utf-8")
    c.Header("Cache-Control", "public, max-age=3600")
    c.JSON(http.StatusOK, gin.H{
        "name": h.cfg.App.Name, "short_name": h.cfg.App.Name,
        "start_url": "/", "display": "standalone",
        "background_color": "#FAFAF8", "theme_color": "#A67C89",
        "icons": []gin.H{
            {"src": "/icons/icon-192.png", "sizes": "192x192", "type": "image/png"},
            {"src": "/icons/icon-512.png", "sizes": "512x512", "type": "image/png"},
            {"src": "/icons/icon-maskable-512.png", "sizes": "512x512", "type": "image/png", "purpose": "maskable"},
        },
    })
}
```

> Gin 显式路由优先于 NoRoute，`/manifest.webmanifest` 不会被 SPA 回退拦截。图标放 `frontend/public/icons/`。

---

## 13. API 设计

统一响应 `{ "code": 0, "message": "ok", "data": {} }`。除标注「否」外均需 `Authorization: Bearer <token>`。

| 方法 | 路径 | 鉴权 | 说明 |
| --- | --- | --- | --- |
| POST | `/api/v1/auth/login` | 否 | 登录 |
| POST | `/api/v1/auth/refresh` | 否 | 刷新 Token |
| POST | `/api/v1/auth/logout` | 是 | 退出 |
| GET | `/api/v1/meta` | 否 | 产品名/版本/脱敏态/PWA 开关 |
| GET | `/api/v1/dicts` | 否 | 枚举字典 |
| GET | `/manifest.webmanifest` | 否 | 动态 PWA manifest |
| GET | `/api/v1/version` | 否 | 版本 + 构建时间 |
| GET | `/api/v1/dashboard` | 是 | 仪表盘数据 |
| GET | `/api/v1/items` | 是 | 列表（筛选/搜索/排序/分页） |
| POST | `/api/v1/items` | 是 | 新增 |
| GET | `/api/v1/items/:id` | 是 | 详情 |
| PUT | `/api/v1/items/:id` | 是 | 更新 |
| DELETE | `/api/v1/items/:id` | 是 | 删除 |
| POST | `/api/v1/items/batch` | 是 | 批量删除/改状态 |
| GET | `/api/v1/items/export` | 是 | 导出 CSV/Excel |
| POST | `/api/v1/upload` | 是 | 上传图片 |
| GET | `/api/v1/categories` | 是 | 分类树 |
| POST/PUT/DELETE | `/api/v1/categories/:id` | 是 | 分类维护 |
| GET | `/api/v1/stats/overview` | 是 | 统计概览 |
| GET | `/api/v1/stats/by-category` | 是 | 分类统计 |
| GET | `/api/v1/stats/by-brand` | 是 | 品牌统计 |
| GET | `/api/v1/stats/by-year` | 是 | 年度统计 |
| POST | `/api/v1/items/:id/wear-logs` | 是 | 新增穿着记录 |
| GET | `/api/v1/items/:id/wear-logs` | 是 | 穿着记录列表 |
| GET/PUT | `/api/v1/settings` | 是 | 设置读写 |
| GET | `/api/v1/backup` | 是 | 下载备份 |

**列表请求示例**：
```
GET /api/v1/items?category_id=1&status=wearing&season=winter&color=white
    &keyword=针织&min_price=100&max_price=500&sort=purchase_date&order=desc&page=1&size=20
```

**响应示例**：
```json
{
  "code": 0, "message": "ok",
  "data": {
    "list": [{
      "id": 1, "name": "白色针织开衫", "category": "衣服",
      "sub_category": "毛衣", "brand": "COS",
      "color": [{"code":"white","label":"白","hex":"#FFFFFF"}],
      "size": "M", "purchase_date": "2025-11-03",
      "purchase_price": null, "purchase_price_display": "¥ ***",
      "status": {"code":"wearing","label":"在穿"},
      "cover": "/uploads/thumb/xxx.jpg", "wear_count": 12,
      "price_masked": true
    }],
    "total": 128, "page": 1, "size": 20
  }
}
```

---

## 14. 部署

### 14.1 拓扑

```
Internet ──443/80──▶ Caddy（自动 Let's Encrypt）
                       │ 反代
                       ▼
                    yichu:8080（容器内网，不映射宿主机端口）
                       ├─ MySQL 8（内网）
                       └─ /data/uploads
```

多端同域：`/` 前端、`/api` 接口、`/uploads` 图片、`/manifest.webmanifest` 动态。

### 14.2 Caddyfile

```caddyfile
yichu.example.com {
    encode zstd gzip
    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains"
        X-Content-Type-Options "nosniff"
        Referrer-Policy "strict-origin-when-cross-origin"
    }
    reverse_proxy yichu:8080
    handle_errors { respond "{err.status_code} {err.status_text}" }
}
```

### 14.3 docker-compose.yml

```yaml
services:
  mysql:
    image: mysql:8.4
    restart: always
    command: --character-set-server=utf8mb4 --collation-server=utf8mb4_0900_ai_ci
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: yichu
      MYSQL_USER: yichu
      MYSQL_PASSWORD: ${MYSQL_PASSWORD}
      TZ: Asia/Shanghai
    volumes: [mysql_data:/var/lib/mysql]
    healthcheck:
      test: ["CMD","mysqladmin","ping","-h","localhost","-p${MYSQL_ROOT_PASSWORD}"]
      interval: 5s
      retries: 20

  yichu:
    image: yichu/yichu:1.5.0
    restart: always
    depends_on:
      mysql: { condition: service_healthy }
    environment:
      YICHU_APP_NAME: 衣橱
      YICHU_SERVER_ADDR: ":8080"
      YICHU_DATABASE_DSN: "yichu:${MYSQL_PASSWORD}@tcp(mysql:3306)/yichu?charset=utf8mb4&parseTime=True&loc=Local"
      YICHU_AUTH_JWT_SECRET: ${JWT_SECRET}
      YICHU_STORAGE_UPLOAD_DIR: /data/uploads
      YICHU_LOG_FILE: /data/logs/yichu.log
      TZ: Asia/Shanghai
    volumes: [app_data:/data]

  caddy:
    image: caddy:2-alpine
    restart: always
    depends_on: [yichu]
    ports: ["80:80","443:443"]
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy_data:/data
      - caddy_config:/config

volumes: { mysql_data: {}, app_data: {}, caddy_data: {}, caddy_config: {} }
```

密码走 `.env`（写入 `.gitignore`）；`app_data` 含 `uploads/` 与 `logs/`，升级不丢数据。

### 14.4 备选内置 TLS

```yaml
server:
  addr: ":8443"
  tls:
    enabled: true
    cert_file: "/data/certs/fullchain.pem"
    key_file: "/data/certs/privkey.pem"
```

与 Caddy 互斥，配置校验失败直接退出并打印原因。

### 14.5 systemd（非 Docker）

```ini
[Unit]
Description=Yichu Wardrobe Manager
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/yichu
ExecStart=/opt/yichu/yichu -c /opt/yichu/config.yaml
Restart=always
RestartSec=3
User=yichu

[Install]
WantedBy=multi-user.target
```

---

## 15. 构建与工程

### 15.1 Dockerfile

```dockerfile
FROM node:20-alpine AS fe
WORKDIR /src
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm ci
COPY frontend/ ./frontend/
COPY web/ ./web/
RUN cd frontend && npm run build

FROM golang:1.25-alpine AS be
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=fe /src/web/dist ./web/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/yichu ./cmd/yichu

FROM alpine:3.20
WORKDIR /data
RUN adduser -D -H yichu && apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=be /out/yichu /usr/local/bin/yichu
USER yichu
EXPOSE 8080
VOLUME ["/data"]
ENV TZ=Asia/Shanghai YICHU_SERVER_ADDR=:8080 \
    YICHU_DATABASE_DSN="yichu:yichu@tcp(mysql:3306)/yichu?charset=utf8mb4&parseTime=True&loc=Local" \
    YICHU_STORAGE_UPLOAD_DIR=/data/uploads YICHU_LOG_FILE=/data/logs/yichu.log
ENTRYPOINT ["yichu"]
```

> MySQL 驱动纯 Go，`CGO_ENABLED=0` 交叉编译成立。

### 15.2 Makefile

```makefile
APP     := yichu
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD   := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD)'

.PHONY: all dev front verify build run cross up down logs backup clean

all: front verify build

front:
	cd frontend && npm ci && npm run build

verify:
	@test -f web/dist/index.html || (echo "❌ web/dist/index.html 缺失" && exit 1)

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP) ./cmd/yichu

run: all
	./bin/$(APP)

dev:
	@echo "后端: go run ./cmd/yichu   前端: cd frontend && npm run dev"

cross: front verify
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP)-linux-amd64     ./cmd/yichu
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP)-linux-arm64     ./cmd/yichu
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP)-darwin-arm64    ./cmd/yichu
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP)-windows-amd64.exe ./cmd/yichu

up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f yichu

backup:
	docker compose exec -T mysql mysqldump --single-transaction -uyichu -p$$MYSQL_PASSWORD yichu | gzip > backup/yichu_$$(date +%Y%m%d_%H%M).sql.gz

clean:
	rm -rf bin && rm -rf web/dist/* && touch web/dist/.gitkeep
```

### 15.3 Vite 配置

```ts
export default defineConfig({
  base: '/',
  plugins: [vue(), VitePWA({ /* 见第 10 节 */ })],
  resolve: { alias: { '@': path.resolve(__dirname, 'src') } },
  build: {
    outDir: '../web/dist',
    emptyOutDir: true,
    assetsDir: 'assets',
    rollupOptions: { output: { manualChunks: {
      vendor: ['vue','vue-router','pinia'],
      ui: ['element-plus'], chart: ['echarts'],
    }}},
  },
  server: { port: 5173, proxy: { '/api': { target: 'http://localhost:8080', changeOrigin: true } } },
})
```

### 15.4 .gitignore

```gitignore
/bin/
*.exe
frontend/node_modules/
frontend/dist/
web/dist/*
!web/dist/.gitkeep
*.db
/uploads/
/logs/
/data/
.env
/config.yaml
.idea/
.vscode/
.DS_Store
```

### 15.5 工作流

| 场景 | 命令 |
| --- | --- |
| 日常开发 | 终端1 `go run ./cmd/yichu`；终端2 `cd frontend && npm run dev`（:5173，`/api` 代理） |
| 本地验证生产形态 | `make all && ./bin/yichu` |
| 全栈启动 | `docker compose up -d --build` |
| 发版 | `git tag v1.5.0 && make cross` |

---

## 16. 前端工程

```
frontend/
├── src/
│   ├── api/             # 按模块封装请求
│   ├── components/      # 原子组件 + 业务组件（含 Money.vue）
│   ├── composables/     # useBreakpoint / useTheme / useUpload / useDict
│   ├── layouts/         # DesktopLayout / TabletLayout / MobileLayout
│   ├── pages/           # login/dashboard/items/categories/stats/settings
│   ├── router/
│   ├── stores/
│   ├── styles/          # tokens.css / element-override.scss
│   ├── types/
│   ├── App.vue
│   └── main.ts
├── public/icons/        # PWA 图标（192 / 512 / maskable-512）
├── vite.config.ts
├── tailwind.config.ts
└── package.json         # name: yichu-web
```

路由：`/login`、`/`（仪表盘）、`/items`、`/items/new`、`/items/:id`、`/items/:id/edit`、`/categories`、`/stats`、`/settings`。

---

## 17. 非功能需求

| 类别 | 要求 |
| --- | --- |
| 性能 | 列表接口 P95 < 300ms（1 万条）；首屏 < 2s；CLS < 0.1 |
| 并发 | 单用户，10 并发足够 |
| 安全 | bcrypt；JWT 7 天可刷新；接口鉴权；上传白名单；路径遍历防护；HSTS；敏感接口 no-store；上传文件名服务端生成 |
| 数据 | MySQL 每日备份保留 14 天；支持手动备份下载；容器重启数据不丢 |
| 兼容 | Chrome/Safari/Edge 近两版；iOS Safari；PWA 安装 |
| 可维护 | 后端分层；前端按模块组织 |
| 日志 | 请求日志 + 错误日志按天切割；不记录响应体与价格 |
| 部署 | 单二进制；上传目录外置；启动自初始化；交叉编译；优雅退出；版本注入；HTTPS |
| 可配置 | 产品名、脱敏开关、主题、PWA 均配置驱动 |

---

## 18. 迭代计划

| 版本 | 内容 | 预估 |
| --- | --- | --- |
| V1.0（MVP） | 登录、分类、衣物 CRUD、图片、筛选搜索、仪表盘、三端响应式、embed 单文件、MySQL、价格脱敏、枚举字典、深色模式、PWA 基础、动态产品名、Docker Compose + Caddy | 3.5 周 |
| V1.1 | 统计报表、穿着记录、标签、导出、存放位置、状态流转、电脑端表格与快捷键、PWA 更新提示优化 | 1.8 周 |
| V1.2 | 心愿单、搭配套装、版本化 SQL 迁移、备份恢复、CI 交叉编译与镜像发布 | 1 周 |

**MVP 完成标准**：可登录保持会话；衣物增删改查全流程可用；图片可上传压缩预览；手机端可完成一次完整录入；首页有数量统计；前后端均有错误提示；`docker compose up -d` 全栈可跑通且 HTTPS 访问正常；`/items/1` 刷新不 404。

---

## 19. 验收标准

**功能**
- [ ] 衣物 CRUD、筛选、搜索、排序、分页全部通过
- [ ] 图片上传/压缩/缩略图/封面设置正常
- [ ] 分类二级维护，有关联时禁止删除并提示
- [ ] 仪表盘数据与列表一致
- [ ] 颜色/季节/场合/状态非法 code 被拒绝

**价格脱敏**
- [ ] `hide_price=true` 时抓包确认所有接口响应不含真实金额
- [ ] 导出价格列为 `***`，响应头含 `X-Price-Masked: true`
- [ ] 切换可见性后 `api-data-*` 缓存被清空
- [ ] 全站无裸金额渲染（`<Money>` 覆盖率 100%）

**多端**
- [ ] 375 / 768 / 1024 / 1440 / 1920 五档无横向滚动、无溢出
- [ ] 手机端完成「拍照→填名→选分类→保存」≤ 60s（iOS + Android 真机各一台）
- [ ] 平板横竖屏切换布局正确
- [ ] 电脑端表格、快捷键、拖拽/粘贴上传可用
- [ ] 手机「添加到主屏幕」全屏打开，图标与名称正确
- [ ] 抽查 20 个可点击元素均 ≥ 44px
- [ ] 深色模式无白底白字，正文对比度 ≥ 4.5:1

**PWA**
- [ ] 手机 Chrome / iOS Safari 均可「添加到主屏幕」
- [ ] 离线可打开首屏并浏览已缓存列表
- [ ] 离线写操作提示「请联网」
- [ ] 服务端更新后提示「有新版本」，点击刷新生效
- [ ] HTTPS 下无混合内容警告，Lighthouse PWA 项通过

**部署**
- [ ] `make all` 产出单个二进制，可独立运行
- [ ] `docker compose up -d` 三服务正常启动
- [ ] 直接访问 `/items/1`、`/stats` 刷新不 404
- [ ] `/api/v1/not-exist` 返回 JSON 404
- [ ] `assets/*.js` 含 `immutable`，`index.html` 为 `no-cache`
- [ ] 上传图片后重启容器，图片仍可访问
- [ ] 删除 `web/dist` 与 `frontend/dist` 后二进制不受影响
- [ ] Linux amd64/arm64、macOS、Windows 四平台产物均可启动
- [ ] `SIGTERM` 后 5s 内正常退出
- [ ] `/api/v1/version` 返回注入版本与构建时间
- [ ] 改 `app.name` 重启后标题/Logo/PWA 名称同步变化，前端未重建
- [ ] `docker compose down && up` 数据不丢
- [ ] Go 1.25.0 编译通过，`go vet` 与 `golangci-lint` 无 error

**视觉**
- [ ] 每页彩色面积 ≤ 10%，主色仅出现在主按钮/选中态/关键数值
- [ ] 分割线仅用 `--color-border`
- [ ] 卡片默认无阴影，仅 1px 边框，hover 才 `shadow-sm`
- [ ] 所有衣物图片 3:4、圆角一致、无拉伸
- [ ] 图片未加载有骨架占位，无白闪/跳动
- [ ] 长名称（30 字）单行省略不撑破卡片
- [ ] 动效时长与曲线符合 5.5，「减少动态效果」下动画关闭且功能正常
- [ ] 空/加载/错误/无权限四态均有设计且文案温柔
- [ ] 所有按钮高度、圆角、字重一致
- [ ] 首屏给非技术用户看 5 秒，能说出「这是记录衣服的」且不觉得复杂

---

## 20. 交付物清单

**后端**
1. `go.mod`（`github.com/iokr/yichu`，`go 1.25.0`）
2. `internal/config`（app.name / hide_price / mysql / tls / 环境变量）
3. `internal/database`（MySQL 连接池 + 重试 + AutoMigrate + 种子 + 首次建账号）
4. `internal/pkg/dict`（四组字典 + 校验）
5. `internal/model` / `dto` / `repository` / `service`（含 mask.go）/ `handler` / `router` / `middleware`
6. `cmd/yichu/main.go` + `web/embed.go` + `web/dist/.gitkeep`

**前端**
7. `frontend/` 完整工程：tokens、Tailwind、Element 主题覆盖、三端布局、深色模式、`Money.vue`、`useDict`/`useTheme`/`useBreakpoint`/`useUpload`、PWA（动态 manifest + SW）、登录/仪表盘/列表/表单/详情/统计/设置页

**部署与工程**
8. `Dockerfile` + `docker-compose.yml` + `Caddyfile` + `.env.example`
9. `Makefile` + `scripts/backup.sh`
10. `systemd` 示例
11. `.gitignore` + `README.md`
12. 可选：GitHub Actions（tag 触发交叉编译 + 镜像发布）

---

**文档终稿，无待确认项，可直接用于代码生成。**