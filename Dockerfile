# syntax=docker/dockerfile:1

# ---------- 前端构建 ----------
FROM node:20-alpine AS fe
WORKDIR /src
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm ci
COPY frontend/ ./frontend/
COPY web/ ./web/
RUN cd frontend && npm run build

# ---------- 后端构建 ----------
FROM golang:1.25-alpine AS be
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=fe /src/web/dist ./web/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/yichu ./cmd/yichu

# ---------- 运行镜像 ----------
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
