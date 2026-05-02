# AGENTS.md - IAM 项目协作指南

## 项目概述

基于 Go + Gin 的 IAM（身份与访问管理）服务，提供 RESTful API 用于用户、角色、权限管理。

## 技术栈

- **语言**: Go 1.25+
- **Web 框架**: Gin
- **ORM**: XORM
- **数据库**: SQLite (modernc.org/sqlite)
- **认证**: JWT (golang-jwt/jwt/v5)
- **密码加密**: bcrypt

## 项目结构

```
iam/
├── cmd/server/          # 入口
│   └── main.go
├── internal/
│   ├── auth/            # (已移除 Casbin，保留 JWT)
│   ├── db/              # 数据库连接
│   ├── handler/         # HTTP Handler
│   ├── middleware/      # JWT 中间件
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   ├── service/         # 业务逻辑层
│   └── test/            # 集成测试
├── migrations/          # SQL 初始化脚本
├── configs/             # 配置文件
└── client.http          # API 测试脚本
```

## API 规范

### 基础约定

- **基础路径**: `/api/v1`
- **认证方式**: `Authorization: Bearer {token}`
- **响应格式**:
  ```json
  {
    "code": "0000",
    "message": "success",
    "data": { ... }
  }
  ```

### 业务码

| Code | 含义 |
|------|------|
| 0000 | 成功 |
| 1001 | 实体不存在 |
| 1002 | 参数错误 |
| 1003 | 未认证 |
| 1004 | 无权限 |
| 1005 | 业务冲突 |
| 1006 | 验证失败 |
| 9999 | 服务器内部错误 |

### 实体 ID 格式

- **URL 路径**: 纯数字，如 `/api/v1/user/1`
- **响应体**: 带模板前缀，如 `"user:1"`, `"role:admin"`

## 数据模型

### 核心实体

- **User**: 用户（id, username, password_hash, role）
- **Role**: 角色（id, name, code, description）
- **Permission**: 权限（id, name, code, resource, action, description）
- **UserRole**: 用户-角色关联
- **RolePermission**: 角色-权限关联
- **EntityFormat**: 实体响应格式配置

### 格式系统

实体响应字段通过 `entity_format` 表动态配置，支持：
- **summary**: 列表模式（摘要信息）
- **detail**: 详情模式（完整信息 + commands）

## 开发指南

### 运行项目

```bash
# 初始化数据库（首次运行）
sqlite3 iam.db < migrations/init.sql

# 运行服务
go run ./cmd/server
```

### 运行测试

```bash
# 全部测试
go test ./internal/test/... -v

# 单个测试
go test ./internal/test -run TestAuth -v
```

### 默认账号

- **用户名**: admin
- **密码**: admin123

## 关键决策

### 1. 移除 Casbin
- **原因**: 项目已通过 `user_role` + `role_permission` 实现 RBAC，Casbin 冗余
- **影响**: 权限检查在 Service 层实现，不再通过中间件

### 2. 动态格式配置
- **实现**: `EntityFormat` 表 + `FormatRegistry` 内存缓存
- **重载**: `POST /api/v1/system/reload-formats`

### 3. Command 模式
- **说明**: 实体操作通过 Command 接口统一处理
- **示例**: `POST /api/v1/user/1/commands` {action: "reset_password"}

## 注意事项

1. **ID 格式转换**: URL 使用纯数字，响应体使用 `{template}:{id}`
2. **密码哈希**: 使用 bcrypt，cost 为默认值
3. **JWT Secret**: 开发环境使用硬编码密钥，生产环境通过环境变量 `JWT_SECRET` 设置
4. **数据库**: 开发使用 SQLite，生产建议迁移至 PostgreSQL/MySQL

## 扩展建议

- [ ] 添加操作日志（审计）
- [ ] 实现用户状态字段（active/disabled）
- [ ] 支持 OAuth2 / SSO 登录
- [ ] 添加 Redis 缓存格式配置
- [ ] API 限流与防暴力破解
