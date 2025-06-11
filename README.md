# 临时邮箱系统 (Temporary Mailbox Service)

一个专业的临时邮箱服务系统，允许用户创建临时邮箱用于接收邮件。系统提供DNS管理、邮件接收、临时邮箱生命周期管理和完整的API访问等功能。

## ✨ 特性

- 🚀 **快速创建**：一键创建临时邮箱，支持自定义有效期
- 📧 **实时接收**：基于SMTP服务实时接收邮件
- 🌐 **域名管理**：支持Cloudflare DNS自动配置MX记录
- 🔒 **安全认证**：完整的用户认证系统和JWT令牌管理
- 📱 **现代界面**：基于Vue3和Ant Design的现代化前端界面
- 🐳 **容器化部署**：支持Docker一键部署和开发环境
- 🗄️ **多数据库支持**：支持PostgreSQL、MySQL、SQLite多种数据库
- ⚡ **高性能**：基于Go语言的高性能后端服务

## 🛠️ 技术栈

### 后端技术栈
- **语言**: Go
- **Web框架**: Gin
- **配置管理**: Viper
- **数据库**: PostgreSQL / MySQL / SQLite
- **ORM**: GORM
- **SMTP服务**: go-guerrilla
- **DNS服务**: Cloudflare SDK
- **认证**: JWT
- **部署**: Docker

### 前端技术栈
- **框架**: Vue 3
- **构建工具**: Vite
- **UI组件库**: Ant Design Vue
- **路由管理**: Vue Router 4
- **状态管理**: Pinia
- **HTTP客户端**: Axios
- **开发语言**: TypeScript

## 🚀 快速开始

### 环境要求

- Go 1.19+
- Node.js 16+
- Docker & Docker Compose
- Git

### 克隆项目

```bash
git clone https://github.com/your-username/temp-mailbox-service.git
cd temp-mailbox-service
```

### 开发环境启动

1. **复制环境配置文件**
```bash
cp .env.example .env.dev
```

2. **启动Docker开发环境**
```bash
docker-compose -f docker-compose.dev.yml up -d
```

3. **启动后端服务**
```bash
cd backend
go mod download
go run cmd/server/main.go
```

4. **启动前端服务**
```bash
cd frontend
npm install
npm run dev
```

5. **访问应用**
- 前端界面: http://localhost:3000
- 后端API: http://localhost:8080
- API文档: http://localhost:8080/swagger

## 📖 项目结构

```
temp-mailbox-service/
├── backend/           # Go后端项目
│   ├── cmd/          # 应用程序入口
│   ├── internal/     # 内部业务逻辑
│   ├── pkg/          # 公共工具包
│   └── configs/      # 配置文件
├── frontend/         # Vue前端项目
│   ├── src/          # 源代码
│   ├── public/       # 静态资源
│   └── dist/         # 构建输出
├── docs/             # 项目文档
├── scripts/          # 构建脚本
└── deploy/           # 部署配置
```

详细的目录结构说明请参考：[项目结构文档](docs/临时邮箱系统开发计划.md#项目目录结构)

## 🔧 配置说明

### 环境变量配置

创建 `.env.dev` 文件用于开发环境：

```env
# 应用配置
APP_NAME=temp-mailbox-service
APP_VERSION=1.0.0
APP_ENV=development

# 服务端口
HTTP_PORT=8080
SMTP_PORT=2525

# 数据库配置
DB_TYPE=sqlite                    # 开发环境使用SQLite
DB_PATH=./dev.db                  # SQLite数据库文件路径

# 生产环境使用PostgreSQL
# DB_TYPE=postgres
# DB_HOST=localhost
# DB_PORT=5432
# DB_NAME=temp_mailbox
# DB_USER=postgres
# DB_PASSWORD=password

# JWT配置
JWT_SECRET=your-jwt-secret-key
JWT_EXPIRES_IN=24h

# Cloudflare DNS配置
CLOUDFLARE_API_TOKEN=your-cloudflare-api-token
CLOUDFLARE_ZONE_ID=your-zone-id

# SMTP配置
SMTP_HOST=0.0.0.0
SMTP_PORT=2525
SMTP_ALLOWED_HOSTS=your-domain.com

# Redis配置 (可选)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

### 数据库配置

系统支持多种数据库，通过 `DB_TYPE` 环境变量切换：

- **SQLite** (开发推荐): 无需额外安装，适合开发和测试
- **PostgreSQL** (生产推荐): 高性能，支持高并发
- **MySQL**: 兼容性好，广泛使用

## 📚 API 文档

### 核心API端点

| 接口类型 | 路径 | 功能描述 |
|---------|------|----------|
| POST | `/api/auth/register` | 用户注册 |
| POST | `/api/auth/login` | 用户登录 |
| POST | `/api/domains` | 添加域名 |
| GET | `/api/domains` | 获取域名列表 |
| POST | `/api/emails` | 创建临时邮箱 |
| GET | `/api/emails` | 获取临时邮箱列表 |
| GET | `/api/emails/{id}/messages` | 获取邮箱邮件 |
| GET | `/api/messages/{id}` | 获取邮件详情 |

完整的API文档请参考：[API文档](docs/api.md)

## 🚢 部署说明

### Docker部署 (推荐)

1. **生产环境配置**
```bash
cp .env.example .env.production
# 编辑 .env.production 配置生产环境参数
```

2. **构建和启动**
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### 手动部署

1. **后端部署**
```bash
cd backend
go build -o bin/server cmd/server/main.go
go build -o bin/worker cmd/worker/main.go
./bin/server
```

2. **前端部署**
```bash
cd frontend
npm run build
# 将 dist/ 目录部署到Web服务器
```

详细部署指南请参考：[部署文档](docs/deployment.md)

## 🔒 安全注意事项

- 🔑 **JWT密钥**: 生产环境请使用强随机密钥
- 🌐 **HTTPS**: 生产环境必须启用HTTPS
- 🔒 **防火墙**: 正确配置防火墙规则
- 📧 **SMTP安全**: 配置适当的SMTP访问控制
- 🗄️ **数据库安全**: 使用强密码并限制访问

## 🧪 开发指南

### 代码规范

- Go代码遵循 `gofmt` 和 `golint` 标准
- TypeScript代码使用 ESLint + Prettier
- 提交信息遵循 Conventional Commits 规范

### 测试

```bash
# 后端测试
cd backend
go test ./...

# 前端测试
cd frontend
npm run test
```

### 热重载开发

项目支持热重载开发，代码修改后自动重启：

- 后端：使用 Air 工具自动重载
- 前端：Vite 原生支持热重载

详细开发指南请参考：[开发计划](docs/临时邮箱系统开发计划.md)

## 🤝 贡献指南

我们欢迎所有形式的贡献！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

### 提交规范

请使用以下格式提交代码：

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型说明：
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建工具、依赖管理等

## 📝 版本历史

### v1.0.0 (计划中)
- [ ] 基础用户认证系统
- [ ] DNS域名管理功能
- [ ] 临时邮箱创建和管理
- [ ] SMTP邮件接收服务
- [ ] 现代化前端界面
- [ ] Docker容器化部署

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源协议。

## 📞 支持

如果您在使用过程中遇到问题，可以通过以下方式获取支持：

- 🐛 [提交Issue](https://github.com/your-username/temp-mailbox-service/issues)
- 📧 邮件联系：your-email@domain.com
- 💬 讨论区：[GitHub Discussions](https://github.com/your-username/temp-mailbox-service/discussions)

## 🙏 致谢

感谢以下开源项目的支持：

- [Gin](https://github.com/gin-gonic/gin) - 高性能Go Web框架
- [Vue.js](https://vuejs.org/) - 渐进式JavaScript框架
- [Ant Design Vue](https://antdv.com/) - 企业级UI组件库
- [GORM](https://gorm.io/) - Go ORM库
- [go-guerrilla](https://github.com/flashmob/go-guerrilla) - Go SMTP服务器

---

⭐ 如果这个项目对您有帮助，请给我们一个Star！ 