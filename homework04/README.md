# Blog API - Gin + GORM 博客系统

基于 Gin + GORM 的 Go Web 应用，实现用户认证、文章管理的 RESTful API。

## 运行环境

- **Go**: 1.25+
- **MySQL**:  8.0
- **操作系统**: Linux / macOS / Windows

## 依赖安装

### 1. 安装 Go 环境

确保已安装 Go 1.25 或更高版本：

```bash
go version
```

### 2. 安装 MySQL

确保 MySQL 服务已启动并创建数据库：

```sql
CREATE DATABASE gorm_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 安装项目依赖

进入项目根目录（`awesomeProject/`，即 `go.mod` 所在位置）：

```bash
cd /path/to/awesomeProject
go mod download
```

主要依赖：

| 依赖 | 说明 |
|------|------|
| github.com/gin-gonic/gin | Web 框架 |
| gorm.io/gorm | ORM 框架 |
| gorm.io/driver/mysql | MySQL 驱动 |
| github.com/spf13/viper | 配置管理 |
| go.uber.org/zap | 结构化日志 |
| github.com/dgrijalva/jwt-go | JWT 认证 |
| golang.org/x/crypto | 密码加密 |

## 配置说明

配置文件位于 `homework04/config.yaml`：

```yaml
server:
  port: 8081              # 服务端口

db:
  type: "MySQL"           # 数据库类型
  host: "127.0.0.1"       # 数据库地址
  port: 3306              # 数据库端口
  user: "root"            # 数据库用户名
  password: "root"        # 数据库密码
  database: "gorm_test"   # 数据库名称

jwt:
  secret: "Luke"          # JWT 签名密钥

api:
  version: "v1"           # API 版本
  prefix: "api"           # API 前缀
```

请根据实际环境修改数据库连接信息和 JWT 密钥。

## 启动方式

**必须从父目录 `awesomeProject/` 运行**（`go.mod` 所在位置）：

```bash
cd /path/to/awesomeProject

# 运行
go run homework04/main.go

# 或构建后运行
go build -o homework04/app homework04/main.go
./homework04/app
```

启动成功后，服务运行在 `http://localhost:8081`。

## API 接口

基础路径：`/api/v1`

### 公开接口（无需认证）

#### 用户注册

```
POST /api/v1/user/register
```

请求体：
```json
{
  "username": "testuser",
  "password": "123456",
  "email": "test@example.com"
}
```

#### 用户登录

```
POST /api/v1/user/login
```

请求体：
```json
{
  "username": "testuser",
  "password": "123456"
}
```

响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    }
  }
}
```

### 鉴权接口（需要 JWT Token）

请求头需携带：
```
Authorization: Bearer <token>
```

#### 创建文章

```
POST /api/v1/post/create
```

请求体：
```json
{
  "title": "文章标题",
  "content": "文章内容"
}
```

#### 获取文章列表

```
GET /api/v1/post/list/:userId?page=1&pageSize=10
```

#### 获取文章详情

```
GET /api/v1/post/detail/:id
```

#### 更新文章

```
PUT /api/v1/post/update/:id
```

请求体：
```json
{
  "title": "新标题",
  "content": "新内容"
}
```

#### 删除文章

```
DELETE /api/v1/post/delete/:id
```

## 项目结构

```
homework04/
├── main.go              # 应用入口
├── config.yaml          # 配置文件
├── api/                 # 控制器层
│   ├── enter.go         # Service 实例入口
│   ├── userApi.go       # 用户接口
│   └── postApi.go       # 文章接口
├── service/             # 业务逻辑层
│   ├── enter.go         # Service 入口
│   ├── interface.go     # 接口定义
│   ├── userService.go   # 用户服务
│   └── postService.go   # 文章服务
├── model/               # 数据模型
│   ├── baseModel.go     # 基础模型
│   ├── user.go          # 用户模型
│   ├── post.go          # 文章模型
│   └── comment.go       # 评论模型
├── router/              # 路由定义
│   ├── enter.go         # 路由初始化
│   ├── routerGroup.go   # 路由组聚合
│   ├── userRouter.go    # 用户路由
│   └── PostRouter.go    # 文章路由
├── middleware/          # 中间件
├── common/              # 公共模块
│   ├── request/         # 请求 DTO
│   ├── response/        # 响应 VO
│   └── config.go        # 配置结构
├── global/              # 全局状态
└── initialize/          # 初始化模块
```

## 启动流程

1. `LoadSystemConfig()` - 加载 config.yaml 配置
2. `InitLogger()` - 初始化 Zap 日志
3. `OpenDatabase()` - 连接 MySQL 数据库
4. `InitDatabase()` - AutoMigrate 自动建表
5. `InitServer()` - 启动 Gin HTTP 服务
