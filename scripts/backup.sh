#!/usr/bin/env bash
# 每日 MySQL 备份，保留 14 天。建议 crontab：0 3 * * * /opt/yichu/scripts/backup.sh
set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-/data/backup}"
KEEP_DAYS="${KEEP_DAYS:-14}"
COMPOSE_FILE="${COMPOSE_FILE:-/opt/yichu/docker-compose.yml}"
DB_USER="${DB_USER:-yichu}"
DB_NAME="${DB_NAME:-yichu}"

mkdir -p "$BACKUP_DIR"
STAMP="$(date +%Y%m%d_%H%M)"
TARGET="$BACKUP_DIR/${DB_NAME}_${STAMP}.sql.gz"

if docker compose -f "$COMPOSE_FILE" ps --status running --services | grep -qx mysql; then
  docker compose -f "$COMPOSE_FILE" exec -T mysql \
    mysqldump --single-transaction --quick --default-character-set=utf8mb4 \
    -u"$DB_USER" -p"${MYSQL_PASSWORD:?MYSQL_PASSWORD 未设置}" "$DB_NAME" | gzip >"$TARGET"
else
  echo "❌ mysql 容器未运行，跳过本次备份" >&2
  exit 1
fi

echo "✅ 备份完成: $TARGET ($(du -h "$TARGET" | cut -f1))"

find "$BACKUP_DIR" -name "${DB_NAME}_*.sql.gz" -type f -mtime "+$KEEP_DAYS" -delete
echo "🧹 已清理 ${KEEP_DAYS} 天前的备份"
