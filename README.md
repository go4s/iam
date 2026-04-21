# IAM (Identity and Access Management) Service

本项目是一个基于 Go 语言开发的身份管理与访问控制基础服务，集成了用户认证 (Authentication) 和基于角色的授权 (Authorization - RBAC) 功能。

## ⚙️ 配置说明

系统支持通过环境变量进行配置：
- `JWT_SECRET`: JWT 签名密钥（建议生产环境必须设置）。

## 🚀 快速启动

1.  **编译运行**
    ```bash
    go run ./cmd/server/main.go
    ```
    服务默认监听 `:8080` 端口。

2.  **API 测试**
    项目根目录提供 `client.http` 文件，可使用 IntelliJ IDEA 或 VS Code 的 REST Client 插件直接运行测试。

## 📚 文档指南

为了方便不同角色的人员快速上手，我们准备了以下详细手册：

### 1. [管理员使用手册 (ADMIN_GUIDE.md)](docs/ADMIN_GUIDE.md)
-   如何配置 RBAC 权限策略。
-   如何管理角色继承关系。
-   如何注册和管理系统资源。

### 2. [业务开发对接手册 (DEV_GUIDE.md)](docs/DEV_GUIDE.md)
-   网关如何通过远程或本地方式校验 Token。
-   微服务如何调用远程鉴权接口 (`/enforce`)。
-   系统错误码及处理建议。

## 🛠 技术栈
-   **框架**: Gin
-   **鉴权引擎**: Casbin
-   **存储**: XORM + SQLite
-   **安全**: JWT + Bcrypt

## 📂 项目结构
-   `cmd/server`: 入口程序。
-   `configs`: 权限模型定义 (`model.conf`)。
-   `internal/handler`: 业务逻辑接口。
-   `internal/middleware`: JWT 与鉴权中间件。
-   `docs`: 详细使用文档。
