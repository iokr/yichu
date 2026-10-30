# 数据库迁移（migrations）

MySQL 8.0+ 的版本化表结构脚本，与 `internal/model` 中的 GORM 模型一一对应。

## 文件

| 文件 | 说明 |
| --- | --- |
| `0001_init.up.sql` | 8 张表完整结构（users / categories / items / item_images / tags / item_tags / wear_logs / settings） |
| `0001_init.down.sql` | 回滚：删除以上全部表（清空数据，仅限开发环境） |
| `0002_seed.up.sql` | 种子数据：26 个二级分类 + 3 个一级分类 + `schema_version`（可选，应用启动时也会幂等补齐） |

命名遵循 `golang-migrate` 约定（`{版本}_{名称}.up.sql` / `.down.sql`），可直接被 `golang-migrate`、`goose` 等工具识别。

## 两种建表方式

**A. 应用自动迁移（默认）**

`config.yaml` 中 `database.auto_migrate: true`，启动时 `AutoMigrate` + `Seed` 自动建表并写入种子数据，无需手工执行本目录脚本。

**B. 纯 SQL 迁移（DBA 手动建库 / 已有流水线）**

```yaml
database:
  auto_migrate: false   # 关闭 AutoMigrate，仅由 SQL 迁移维护结构
```

```bash
# 本地 mysql 客户端
mysql --default-character-set=utf8mb4 -h 127.0.0.1 -u yichu -p yichu < migrations/0001_init.up.sql
mysql --default-character-set=utf8mb4 -h 127.0.0.1 -u yichu -p yichu < migrations/0002_seed.up.sql

# Docker Compose 环境
make migrate          # 等价于 docker compose exec -T mysql mysql -u yichu -p"$MYSQL_PASSWORD" yichu < ...
```

两种方式可以混用：SQL 建表后再以 `auto_migrate: true` 启动，`AutoMigrate` 只补齐缺失的列与索引，不会清空数据。

## 与模型的一致性约定

- 表只声明 `ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`，不写 `COLLATE`，排序规则跟随数据库默认。
- `ID uint64` → `bigint unsigned AUTO_INCREMENT`；其余 `uint64` 外键 → `bigint unsigned`。
- Go `int` → `bigint`（64 位平台）；`uint32` 显式声明为 `int unsigned`。
- `model.Codes` → `json`（颜色 / 季节 / 场合 code 数组，筛选用 `JSON_CONTAINS`）。
- 时间统一 `datetime(3)`；`purchase_date`、`worn_date` 为 `date`；`items.deleted_at` 为软删除列。
- 索引名与模型 tag 完全一致，例如 `idx_category`、`idx_deleted_at`、`uk_username`、`uk_tag_name`；
  连接表 `item_tags` 以 `(item_id, tag_id)` 为联合主键，模型未声明外键约束。
- 金额 `decimal(10,2)`：`hide_price=true` 时由服务端脱敏，数据库中始终保存真实值。

## 新增变更

1. 文件名递增：`0003_add_item_color_sort.up.sql` / `.down.sql`。
2. 同时修改 `internal/model` 中的模型，保证 AutoMigrate 与 SQL 两条路径结果一致。
3. 已发布（合并到主分支）的迁移文件不改写，变更一律新增文件。
