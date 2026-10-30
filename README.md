# 衣橱 · yichu

记录与管理家庭衣物的自托管应用：Go 单二进制后端 + 内嵌 Vue 3 前端，支持手机 / 平板 / 电脑三端，可按需隐藏价格。

## 技术栈

| 层 | 选型 |
| --- | --- |
| 后端 | Go 1.25.0、gin、gorm、MySQL 8 |
| 认证 | JWT（7 天）+ 刷新 Token、bcrypt |
| 前端 | Vue 3 + TypeScript + Vite + Pinia + Vue Router + Tailwind + ECharts + PWA |
| 部署 | 单二进制（`embed` 前端）、Docker Compose + Caddy 自动 HTTPS |

## 快速开始（本地开发）

```bash
# 1. 准备 MySQL 8 并建库
CREATE DATABASE yichu DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

# 2. 配置（首次运行会读取仓库根 config.yaml，可直接改或使用环境变量）
cp config.yaml config.local.yaml   # 按需修改 database.dsn / auth.jwt_secret

# 3. 启动后端（自动迁移 + 种子 + 首次建管理员）
go run ./cmd/yichu -c config.yaml

# 如需用 SQL 维护表结构（关闭 database.auto_migrate 后执行）：
#   mysql -u yichu -p yichu < migrations/0001_init.up.sql
#   mysql -u yichu -p yichu < migrations/0002_seed.up.sql

# 4. 启动前端开发服务器（:5173，/api 代理到 :8080）
cd frontend && npm install && npm run dev
```

默认管理员：`admin / yichu123`（来自 `auth.init_username` / `auth.init_password`，首次启动时创建，登录后请尽快修改）。

## 生产构建

```bash
make all          # 前端构建 → web/dist，再编译 bin/yichu（前端已内嵌）
./bin/yichu       # 单文件运行，默认 :8080
```

交叉编译：

```bash
make cross        # bin/yichu-{linux-amd64,linux-arm64,darwin-arm64,windows-amd64.exe}
```

## Docker Compose

```bash
cp .env.example .env      # 修改 MYSQL_ROOT_PASSWORD / MYSQL_PASSWORD / JWT_SECRET
# 将 Caddyfile 中的域名改成你的域名
docker compose up -d --build
```

拓扑：`Caddy(80/443) → yichu:8080 → MySQL`，图片与日志落在 `app_data` 卷（`/data`），升级不丢数据。

备份：`make backup` 或 `scripts/backup.sh`（默认保留 14 天）。

## 配置

优先级：`config.yaml` < 环境变量（前缀 `YICHU_`，下划线分隔层级）。常用项：

| 环境变量 | 说明 |
| --- | --- |
| `YICHU_APP_NAME` | 产品名，全站标题 / Logo / PWA 名称同步 |
| `YICHU_APP_HIDE_PRICE` | 价格脱敏默认值（运行时开关存 `settings` 表，优先级更高） |
| `YICHU_SERVER_ADDR` | 监听地址，如 `:9000` |
| `YICHU_DATABASE_DSN` | MySQL DSN，支持 `${DB_PASSWORD}` 变量展开 |
| `YICHU_AUTH_JWT_SECRET` | JWT 密钥，生产必须修改 |
| `YICHU_STORAGE_UPLOAD_DIR` | 图片根目录 |
| `YICHU_LOG_FILE` | 日志文件路径（按天切割） |

内置 TLS（与 Caddy 互斥）：

```yaml
server:
  addr: ":8443"
  tls:
    enabled: true
    cert_file: "/data/certs/fullchain.pem"
    key_file: "/data/certs/privkey.pem"
```

## 目录结构

```
cmd/yichu/           入口
internal/
  config/            viper 配置
  database/          连接、迁移、种子
  model/ dto/        数据模型与传输对象
  repository/        GORM 数据访问
  service/           业务逻辑（含 mask.go 价格脱敏）
  handler/ middleware/ router/   HTTP 层
  pkg/               dict / errcode / response / logx / imgproc / netutil
frontend/            Vue 工程（构建产物输出到 web/dist）
web/                 embed.go + dist（前端内嵌）
migrations/          MySQL 表结构迁移（0001 建表 8 张 / 0002 种子数据）
deploy/              systemd 示例
scripts/backup.sh    数据库备份
```

## API

统一响应 `{"code":0,"message":"ok","data":{}}`，错误码：`1001` 参数错误、`1002` 未登录、`1003` 无权限、`1004` 资源不存在、`5000` 服务端错误。

主要端点：

```
POST   /api/v1/auth/login | refresh | logout
GET    /api/v1/meta | dicts | version
GET    /api/v1/dashboard
GET    /api/v1/items            筛选 / 搜索 / 排序 / 分页
POST   /api/v1/items            新增
GET    /api/v1/items/:id        详情
PUT    /api/v1/items/:id        更新
DELETE /api/v1/items/:id        删除
POST   /api/v1/items/batch      批量删除 / 改状态
GET    /api/v1/items/export     导出 CSV / Excel
POST   /api/v1/upload           图片上传（自动压缩 + 缩略图）
GET    /api/v1/categories       分类树（含增删改与排序）
GET    /api/v1/stats/overview | by-category | by-brand | by-year | by-season | by-color | cost-per-wear
POST   /api/v1/items/:id/wear-logs   穿着打卡
GET    /api/v1/settings | PUT    运行时设置
GET    /api/v1/backup            下载 JSON 备份
```

## 价格脱敏

`hide_price=true` 时服务端抹除真实金额，响应契约：

```json
{ "purchase_price": null, "purchase_price_display": "¥ ***", "price_masked": true }
```

统计金额返回 `"***"`，导出价格列为 `***` 且响应头带 `X-Price-Masked: true`。前端所有金额只通过 `<Money>` 渲染，切换开关后自动清理 PWA 的 `api-data-*` 缓存。

## License

私有项目，仅供个人与家庭使用。
