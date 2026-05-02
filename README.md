# ClipLink - 跨平台剪贴板共享工具

[![Go](https://img.shields.io/badge/Go-1.25%2B-blue)](https://golang.org/)
[![Node.js](https://img.shields.io/badge/Node.js-22%2B-green)](https://nodejs.org/)
[![React](https://img.shields.io/badge/React-19.2%2B-blue)](https://reactjs.org/)
[![Next.js](https://img.shields.io/badge/Next.js-16.2%2B-lightgrey)](https://nextjs.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-4.2%2B-38B2AC)](https://tailwindcss.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

ClipLink 是一个功能强大的跨平台文字剪贴板同步工具，允许您在不同设备（如电脑和手机）之间通过网页界面共享文本、链接、代码、密码和富文本内容。该项目采用前后端分离架构，后端使用 Go 语言构建，数据通过 SQLite 存储并通过网络同步，前端使用 Next.js 和 React 构建。通过内置的编译脚本，可以将前端静态资源嵌入到 Go 二进制文件中，实现前后端一体化部署。

**演示网站：** 👉 [https://cliplink.mmmss.com/](https://cliplink.mmmss.com/) - 立即体验ClipLink的强大功能！

**使用说明：** 在大多数现代浏览器和设备上，Web 页面具有读取剪贴板的权限，可自动获取剪贴板内容。但在 iOS 等部分平台，由于系统安全策略限制，剪贴板内容无法自动读取，需要手动粘贴。使用流程简单：打开网页，内容在获得权限后自动同步，或在需要时手动粘贴，即可在多设备间实现共享。

**⚠️ 重要提示：** 为保证剪贴板权限正常获取，请确保通过 **HTTPS** 协议访问网页。HTTP 协议下，现代浏览器将限制剪贴板访问权限。

![应用预览](docs/home.png)

![应用预览](docs/home-1.png)

![应用预览](docs/home-2.png)


## 💡 项目初衷

在当今多设备环境下，我们经常需要在不同设备间共享剪贴板文字内容 —— 可能是一段文字、一个链接、一段命令或一个临时密码。传统的解决方案往往繁琐且低效：

- 需要登录微信、QQ或其他通讯工具，发送给自己或特定联系人
- 依赖第三方云服务，隐私安全无法保障
- 需要安装专用软件，增加系统负担
- 操作复杂，打断工作流程

ClipLink 项目就是为了解决这些痛点而创建的。它提供了一个轻量、安全、高效的方式，让您可以在任何支持浏览器的设备上快速共享和访问剪贴板内容，无需复杂设置，真正实现随剪随用。

## 🌟 主要功能

- **跨设备同步**: 通过网页界面，在手机、平板和电脑之间共享剪贴板内容
- **自动读取**: 在支持的平台上，自动读取剪贴板内容（Windows/macOS/Android等）
- **富文本同步**: 自动识别并保留剪贴板中的富文本格式（粗体、斜体、链接等），在支持的浏览器间实现格式保真同步，不支持时自动降级为纯文本
- **内容类型识别**: 支持文本、链接、代码片段和密码等常用剪贴板内容类型
- **历史记录**: 自动保存内容历史，随时查看和恢复以前的内容
- **收藏与搜索**: 支持收藏常用内容，并按关键词或内容类型筛选
- **网络同步**: 通过网络实时同步内容，在任何设备上都能获取最新数据
- **安全可靠**: 数据存储在您控制的服务器上，不经过第三方云服务
- **一键部署**: 单一二进制文件，无需复杂配置，一键启动
- **通道管理**: 支持创建通道、退出通道和删除通道（含数据清理）
- **重复内容清理**: 自动检测并清理重复的剪贴板内容，保持历史记录整洁
- **Docker 部署**: 支持 Docker 和 docker-compose 一键容器化部署

## 🚀 未来计划

目前，ClipLink 专注文本类剪贴板同步，已实现文本、链接、代码片段、密码、富文本、历史记录、收藏、搜索、通道管理和重复内容清理等核心网页端能力。我们计划在未来版本中逐步增强文字剪贴板体验：

| 功能 | 状态 | 说明 |
|------|------|------|
| 格式化文本支持 | ✅ 已完成 | 保留富文本格式（粗体、斜体、链接等），支持 HTML 剪贴板读写与预览，不支持时自动降级纯文本 |
| 代码片段优化 | 📅 计划中 | 针对代码片段提供语法高亮和格式化功能 |
| 内容分类与标签 | 📅 计划中 | 支持对剪贴板内容进行分类整理，添加标签 |
| 敏感内容保护 | 🔍 调研中 | 增强密码、token、密钥等文本内容的隐藏展示、确认和加密保护 |
| 端到端加密 | 🔍 调研中 | 增强数据传输和存储安全性，保护敏感文本信息 |
| 桌面客户端 | 🔍 调研中 | 提供Windows/macOS/Linux桌面客户端，后台运行，自动同步 |
| Android App | 🔍 调研中 | 开发原生Android应用，提供后台自动监听剪贴板功能 |
| iOS App | 🔍 调研中 | 开发原生iOS应用，提供后台自动监听剪贴板功能 |

我们欢迎社区贡献者参与这些功能的开发，共同打造更强大的跨平台剪贴板工具！

## 📋 目录结构

```
cliplink/
├── cmd/                    # Go主程序目录
│   └── main.go             # 应用程序入口
├── internal/               # Go内部包
│   ├── app/                # HTTP 路由、控制器、用例和静态站点服务
│   ├── common/             # 通用响应与校验
│   ├── config/             # 配置加载
│   ├── domain/             # 模型、仓库接口和服务接口
│   ├── infra/              # 数据库与持久化实现
│   └── static/             # 嵌入式前端静态资源
├── web/                    # 前端目录
│   ├── src/                # React源码
│   ├── public/             # 静态资源
│   └── package.json        # 依赖配置
├── docs/                   # 文档目录
│   ├── build_script.md     # 构建脚本使用文档
│   ├── database-config.md  # 数据库配置文档
│   ├── docker.md           # Docker 部署文档
│   ├── run_script.md       # 运行脚本使用文档
│   ├── text-clipboard-sync-design.md # 文本剪贴板同步设计
│   ├── home.png            # 应用预览图片
│   ├── home-1.png          # 应用预览图片
│   └── home-2.png          # 应用预览图片
├── .dockerignore           # Docker 构建排除文件
├── Dockerfile              # Docker 镜像构建文件
├── docker-compose.yml      # Docker Compose 配置
├── Makefile                # Make 构建配置
├── auto_deploy.sh          # 自动部署脚本
├── build.sh                # 构建脚本
├── run.sh                  # 运行脚本
├── config.example.yml      # 配置文件示例
├── go.mod                  # Go 模块定义
├── go.sum                  # Go 模块校验
└── LICENSE                 # 许可证文件
```

## 🚀 快速开始

### 在线体验

不想自行部署？直接访问我们的演示站点体验ClipLink功能：

🔗 [https://cliplink.mmmss.com/](https://cliplink.mmmss.com/)

在多个设备上打开此链接，即可立即开始共享剪贴板内容。

### 手动部署

1. 从项目的 GitHub Releases 页面下载适合您系统的压缩包
2. 解压到您选择的目录
3. **（可选）配置数据库**：
   - 默认使用 SQLite，无需任何配置
   - 如需使用 MySQL，请将 `config.example.yml` 复制为 `config.yml` 并修改数据库配置
4. 运行应用程序：

```bash
# 确保run.sh有执行权限
chmod +x run.sh

# 启动应用（默认端口8080）
./run.sh start

# 或使用自定义端口启动
./run.sh start --port 3000
```

5. **重要：** 为确保剪贴板功能正常，请配置反向代理（如Nginx）提供HTTPS访问，或使用SSL证书
6. 在浏览器中访问 `https://<服务器域名>:<端口>` 开始使用
7. 在所有需要共享剪贴板的设备上访问同一地址

### Docker 部署

ClipLink 支持 Docker 容器化部署，适合快速搭建和运维管理。

#### 使用 docker-compose（推荐）

```bash
docker compose up -d --build
```

#### 使用 Docker 命令

```bash
# 构建镜像
docker build -t cliplink:latest .

# 使用 SQLite 运行
docker run -d \
  --name cliplink \
  --restart unless-stopped \
  -p 8080:8080 \
  -v cliplink-data:/home/cliplink/.cliplink \
  cliplink:latest
```

#### 自定义配置

如需 MySQL 或自定义端口，将 `config.example.yml` 复制为 `config.yml` 编辑后挂载到容器：

```bash
docker run -d \
  --name cliplink \
  --restart unless-stopped \
  -p 8080:8080 \
  -v cliplink-data:/home/cliplink/.cliplink \
  -v ./config.yml:/app/config.yml:ro \
  cliplink:latest
```

#### 安全相关环境变量

```bash
docker run -d \
  --name cliplink \
  --restart unless-stopped \
  -p 8080:8080 \
  -v cliplink-data:/home/cliplink/.cliplink \
  -e CLIPLINK_ALLOWED_ORIGINS=https://cliplink.example.com \
  -e CLIPLINK_MAX_BODY_BYTES=2097152 \
  -e CLIPLINK_INSTANCE_TOKEN="your-token" \
  cliplink:latest
```

各环境变量说明：
- `CLIPLINK_ALLOWED_ORIGINS`：跨域来源，生产环境建议只填写 HTTPS 域名
- `CLIPLINK_MAX_BODY_BYTES`：请求体大小限制，默认 2 MiB
- `CLIPLINK_INSTANCE_TOKEN`：实例 Token，配置后创建通道需要提供

> 更多 Docker 部署细节请参考 [Docker 部署文档](docs/docker.md)。

### 自动部署

为简化部署流程，我们提供了自动部署脚本，只需一行命令即可完成安装和启动：

```bash
# 使用curl下载并执行（默认在当前目录安装）
curl -fsSL https://sh.mmmss.com/shell/cliplink/auto-deploy.sh | bash

# 或使用wget
wget -O- https://sh.mmmss.com/shell/cliplink/auto-deploy.sh | bash
```

> **注意**: 脚本会直接在当前执行命令的目录下安装应用，无需指定额外目录。如果您希望应用安装在特定位置，请先进入该目录后执行上述命令，或使用下面的`--dir`参数。

#### 自定义部署选项

自动部署脚本支持多种自定义参数：

```bash
# 指定安装目录
curl -fsSL https://sh.mmmss.com/shell/cliplink/auto-deploy.sh | bash -s -- --dir /opt/cliplink

# 指定应用端口
curl -fsSL https://sh.mmmss.com/shell/cliplink/auto-deploy.sh | bash -s -- --port 3000

# 指定应用版本
curl -fsSL https://sh.mmmss.com/shell/cliplink/auto-deploy.sh | bash -s -- --version v1.1

# 组合使用多个参数
curl -fsSL https://sh.mmmss.com/shell/cliplink/auto-deploy.sh | bash -s -- --dir /opt/cliplink --port 3000 --version v1.1
```

自动部署脚本会：
1. 自动检测系统类型和架构
2. 下载对应版本的应用程序包
3. 默认安装在**当前目录**（无需额外创建子目录）
4. 解压并设置正确的执行权限
5. 根据指定参数配置应用
6. 自动启动服务
7. 显示服务器IP、访问地址和常用管理命令

这是在新服务器上最快速部署ClipLink的方式，特别适合快速测试和生产环境部署。

#### 权限说明

如果您在执行过程中遇到权限问题，可能需要：

```bash
# 下载脚本后手动执行
curl -O https://sh.mmmss.com/shell/cliplink/auto-deploy.sh
chmod +x auto-deploy.sh
./auto-deploy.sh

# 或使用sudo（如果目标目录需要特权）
curl -fsSL https://sh.mmmss.com/shell/cliplink/auto-deploy.sh | sudo bash
```

### 使用方法

1. 打开网页界面（必须通过HTTPS访问）
2. 在大多数平台上（Windows/macOS/Android等），剪贴板内容会在授予权限后自动读取并同步
3. 在 iOS 等限制自动读取剪贴板的平台上，需要手动粘贴内容到输入框
4. 所有内容会自动保存并通过网络同步到服务器
5. 在其他设备上打开同一网址，即可查看和使用已保存的内容
6. 历史内容会在页面上列出，可随时查看和重新使用

#### 通道管理

- **创建通道**: 在通道详情弹窗中创建新通道，如服务端启用了实例 Token，需在创建时填写
- **退出通道**: 退出当前通道，仅断开连接，通道数据仍保留
- **删除通道**: 彻底删除通道及其所有数据（剪贴板内容、同步历史、设备关联），需提供实例 Token 并输入完整通道 ID 确认，操作不可恢复

#### 重复内容清理

在设置中开启「自动清理重复内容」后：
- 保存新内容时，会自动移除历史中内容相同的旧条目
- 列表显示时自动去重，保持历史记录整洁
- 开启设置时会立即执行一次全量清理

### 安全配置

如果实例会暴露到公网，建议配置实例 Token 来限制陌生人创建新通道。配置后，只有创建通道接口需要提供 `X-Instance-Token`；读写已有通道仍使用 `X-Channel-ID`。

```yaml
security:
  instance_token: "replace-with-a-long-random-token"
  # 请求体大小限制，默认 2097152（2 MiB）
  max_body_bytes: 2097152
```

也可以用环境变量覆盖：

```bash
export CLIPLINK_INSTANCE_TOKEN="replace-with-a-long-random-token"
export CLIPLINK_MAX_BODY_BYTES=2097152
```

生成随机 Token 示例：

```powershell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }))
```

删除通道属于管理操作，需要同时提供当前通道的 `X-Channel-ID` 和实例 Token。删除后会清理该通道的剪贴板内容、同步历史和设备关联，并自动删除超过 30 天没有任何通道关联的孤儿设备。

### CORS 跨域配置

生产环境推荐通过同域名反向代理访问，一般不需要额外跨域来源。如果前端和后端分属不同源，请只填写可信来源：

```yaml
cors:
  allowed_origins:
    - "https://cliplink.example.com"
    - "http://localhost:3000"
```

环境变量覆盖：

```bash
export CLIPLINK_ALLOWED_ORIGINS="https://cliplink.example.com,http://localhost:3000"
```

> 默认允许 `http://localhost:3000` 和 `http://127.0.0.1:3000` 跨域访问，方便本地开发。生产环境建议只填写你的 HTTPS 域名。

### 日志配置

可控制 SQL 日志级别，生产环境建议 `warn` 或 `silent`：

```yaml
log:
  sql: warn
```

环境变量覆盖：

```bash
export CLIPLINK_SQL_LOG=warn
```

### 数据库配置

ClipLink 支持两种数据库类型：

#### SQLite（默认，零配置）
- **无需任何配置**，应用启动后自动使用 SQLite 数据库
- 数据存储在 `~/.cliplink/cliplink.db`
- 适合个人使用和小规模部署

#### MySQL（生产推荐）
如果需要使用 MySQL 数据库（推荐生产环境使用），请按以下步骤配置：

1. **准备 MySQL 数据库**：
   ```sql
   CREATE DATABASE cliplink CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   CREATE USER 'cliplink'@'localhost' IDENTIFIED BY 'your_password';
   GRANT ALL PRIVILEGES ON cliplink.* TO 'cliplink'@'localhost';
   FLUSH PRIVILEGES;
   ```

2. **创建配置文件**：
   将 `config.example.yml` 复制为 `config.yml` 并放在与二进制文件相同的目录下：
   ```bash
   cp config.example.yml config.yml
   ```

3. **修改配置文件**：
   编辑 `config.yml`，配置 MySQL 连接信息：
   ```yaml
   mysql:
     host: "localhost"
     port: 3306
     username: "cliplink"
     password: "your_password"
     database: "cliplink"
     charset: "utf8mb4"
   ```

4. **启动应用**：
   ```bash
   ./cliplink
   ```

**注意事项**：
- 配置文件必须与 cliplink 二进制文件放在同一目录下
- 如果 MySQL 连接失败，应用会自动降级使用 SQLite，确保服务正常运行
- 支持热切换：可以随时通过修改或删除配置文件来切换数据库类型

## 💻 开发指南

### 前提条件

- Go 1.25+
- Node.js 22+ 和 npm
- 支持 Linux/macOS/Windows 的开发环境

### 前端开发

前端使用 Next.js 16、React 19、Tailwind CSS 4 和 npm。Tailwind v4 通过 `@tailwindcss/postcss` 接入 PostCSS，样式入口位于 `web/src/app/globals.css`。

```bash
# 进入前端目录
cd web

# 安装依赖
npm install

# 开发模式运行
npm run dev

# 类型检查与生产构建
npm exec tsc -- --noEmit
npm run build
```

### 后端开发

> **注意**：`go build` / `go run` 前需先构建前端，否则会因缺少 dist 目录而编译失败。

```bash
# 1. 先构建前端（生成 internal/static/dist/）
./build.sh --frontend

# 2. 再启动后端
go run cmd/main.go

# 或使用自定义端口
go run cmd/main.go --port 3000
```

## 📝 详细文档

本项目提供了详细的构建和部署文档，可以在以下文件中找到：

- **[构建脚本使用指南 (build.sh)](docs/build_script.md)** - 详细介绍构建工具的使用方法、参数选项和功能
- **[运行脚本使用指南 (run.sh)](docs/run_script.md)** - 全面说明运行脚本的所有命令、选项和流程
- **[Docker 部署指南](docs/docker.md)** - Docker 和 docker-compose 部署方法、配置选项和安全设置
- **[数据库配置指南](docs/database-config.md)** - SQLite 和 MySQL 数据库配置详解
- **[文本剪贴板同步设计](docs/text-clipboard-sync-design.md)** - 当前产品边界、类型约束和后续文字专项方向

这些文档包含了更多技术细节和进阶使用方法，建议开发者和部署人员完整阅读。

## 🛠️ 构建脚本 (build.sh)

`build.sh` 是一个功能完整的构建工具，可以帮助您构建前端、后端或完整应用。

### 基本用法

```bash
# 确保脚本有执行权限
chmod +x build.sh

# 显示交互式菜单
./build.sh

# 使用默认设置快速构建前后端
./build.sh all

# 只构建前端
./build.sh --frontend

# 只构建后端（指定目标系统为macOS）
./build.sh --backend --os darwin
```

更多详细信息，请参考 [构建脚本使用指南](docs/build_script.md)。

## 🔄 运行脚本 (run.sh)

`run.sh` 脚本提供了应用程序的全生命周期管理，包括启动、停止、重启、状态查看和日志查看。

### 基本用法

```bash
# 启动应用
./run.sh start

# 停止应用
./run.sh stop

# 查看应用状态
./run.sh status

# 查看应用日志
./run.sh logs
```

更多详细信息，请参考 [运行脚本使用指南](docs/run_script.md)。

## 🌐 跨平台支持

ClipLink 应用支持以下操作系统和架构组合：

| 操作系统 | 支持的架构 |
|----------|------------|
| Linux    | amd64, arm64, 386 |
| macOS    | amd64, arm64 |
| Windows  | amd64, 386 |

## 👥 贡献指南

欢迎贡献代码、报告问题或提供改进建议！请遵循以下步骤：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。
