# Docker 部署

本项目镜像会在构建阶段完成两件事：

1. 使用 Node.js 构建 `web/` 下的 Next.js 静态前端。
2. 将前端产物复制到 `internal/static/dist`，再编译为内嵌静态资源的 Go 服务。

## 构建镜像

```bash
docker build -t cliplink:latest .
```

## 使用 SQLite 运行

SQLite 是默认模式，数据保存在容器用户目录 `/home/cliplink/.cliplink/cliplink.db`。

```bash
docker run -d \
  --name cliplink \
  --restart unless-stopped \
  -p 8080:8080 \
  -v cliplink-data:/home/cliplink/.cliplink \
  cliplink:latest
```

访问：

```text
http://服务器IP:8080
```

剪贴板权限在现代浏览器中通常需要 HTTPS，生产环境建议用 Nginx、Caddy 或云厂商负载均衡在镜像前面终止 TLS。

## 使用 docker compose

```bash
docker compose up -d --build
```

## 使用配置文件

如需 MySQL 或自定义端口，把 `config.example.yml` 复制为 `config.yml`，编辑后挂载到容器：

```bash
docker run -d \
  --name cliplink \
  --restart unless-stopped \
  -p 8080:8080 \
  -v cliplink-data:/home/cliplink/.cliplink \
  -v ./config.yml:/app/config.yml:ro \
  cliplink:latest
```

如果在 `config.yml` 中修改了服务端口，也要同步修改 `-p 宿主机端口:容器端口`。

## 安全相关环境变量

```bash
-e CLIPLINK_ALLOWED_ORIGINS=https://cliplink.example.com
-e CLIPLINK_MAX_BODY_BYTES=2097152
```

`CLIPLINK_ALLOWED_ORIGINS` 只影响跨域浏览器请求。同域访问不需要配置 CORS。生产环境建议只填写你的 HTTPS 域名，不要使用 `*`。

## 临时覆盖端口

```bash
docker run -d \
  --name cliplink \
  --restart unless-stopped \
  -p 3000:3000 \
  -v cliplink-data:/home/cliplink/.cliplink \
  cliplink:latest -port 3000
```
