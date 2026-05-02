# 后端 API 接口文档（v1）

## 基础约定

- **基础路径**：所有 API 统一前缀 `/api/v1`
- **实体 ID 格式**：`{template}:{id}`，例如 `user:1`、`role:admin`
- **实体引用**：对象键以 `:` 结尾表示引用，值为实体 ID 或 ID 数组
- **内容类型**：请求/响应均为 `application/json`
- **认证方式**：请求头携带 `Authorization: Bearer {token}`

---

## 一、认证接口

### 1.1 用户登录

```
POST /api/v1/auth/login
```

**请求体**：
```json
{
  "username": "admin",
  "password": "password123"
}
```

**成功响应（200）**：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "user:1",
    "username": "admin",
    "roles:": ["role:admin"]
  }
}
```

### 1.2 获取当前用户信息

```
GET /api/v1/auth/me
```

**成功响应（200）**：
```json
{
  "id": "user:1",
  "username": "admin",
  "email": "admin@example.com",
  "roles:": ["role:admin"],
  "roles-fold": false
}
```

---

## 二、通用实体接口

所有实体类型共享以下 RESTful 接口模式：

### 2.1 列出实体（List）

```
GET /api/v1/{template}
```

**支持查询参数**：
- `page`：页码（默认 1）
- `size`：每页数量（默认 20）
- `sort`：排序字段
- `q`：全文搜索关键词
- 其他过滤参数根据实体类型而定

**成功响应（200）**：
```json
{
  "data": [
    {
      "id": "user:1",
      "username": "admin",
      "email": "admin@example.com",
      "status": "active"
    },
    {
      "id": "user:2",
      "username": "editor",
      "email": "editor@example.com",
      "status": "active"
    }
  ],
  "pagination": {
    "page": 1,
    "size": 20,
    "total": 100
  }
}
```

**说明**：
- 返回实体摘要数组，不包含嵌套引用详细内容
- 列表响应中**不包含** `commands` 字段

### 2.2 聚焦实体（Focus）

```
GET /api/v1/{template}/{id}
```

**成功响应（200）**：
```json
{
  "id": "user:1",
  "username": "admin",
  "email": "admin@example.com",
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z",
  "roles:": ["role:admin", "role:editor"],
  "roles-fold": false,
  "commands": [
    {
      "action": "reset_password",
      "label": "重置密码",
      "params": [
        {
          "name": "new_password",
          "type": "string",
          "required": true,
          "label": "新密码"
        }
      ],
      "trigger": "/api/v1/user/1/commands",
      "sync": true
    },
    {
      "action": "disable",
      "label": "禁用账户",
      "params": [],
      "trigger": "/api/v1/user/1/commands",
      "sync": true
    }
  ]
}
```

**说明**：
- 返回完整实体对象
- 必须包含 `commands` 数组
- 实体引用字段必须附带 `{ref}-fold` 意图字段

### 2.3 执行命令（Command）

```
POST /api/v1/{template}/{id}/commands
```

**请求体**：
```json
{
  "action": "reset_password",
  "params": {
    "new_password": "newpass123"
  }
}
```

**同步命令响应（200）**：
```json
{
  "success": true,
  "message": "密码重置成功",
  "data": {
    "updated_at": "2024-01-02T12:00:00Z"
  }
}
```

**异步命令响应（202）**：
```json
{
  "success": true,
  "task_id": "task:abc123",
  "message": "处理中，请等待推送通知"
}
```

---

## 三、权限管理模块实体接口

### 3.1 用户（User）

**列出用户**：
```
GET /api/v1/user
```

**聚焦用户**：
```
GET /api/v1/user/{id}
```

**用户数据结构**：
```json
{
  "id": "user:1",
  "username": "admin",
  "email": "admin@example.com",
  "phone": "13800138000",
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-02T00:00:00Z",
  "roles:": ["role:admin"],
  "roles-fold": false,
  "commands": [
    {
      "action": "reset_password",
      "label": "重置密码",
      "params": [
        {
          "name": "new_password",
          "type": "string",
          "required": true,
          "label": "新密码"
        }
      ],
      "trigger": "/api/v1/user/1/commands",
      "sync": true
    },
    {
      "action": "assign_role",
      "label": "分配角色",
      "params": [
        {
          "name": "role_id",
          "type": "entity",
          "template": "role",
          "required": true,
          "label": "角色"
        }
      ],
      "trigger": "/api/v1/user/1/commands",
      "sync": true
    },
    {
      "action": "disable",
      "label": "禁用账户",
      "params": [],
      "trigger": "/api/v1/user/1/commands",
      "sync": true
    }
  ]
}
```

### 3.2 角色（Role）

**列出角色**：
```
GET /api/v1/role
```

**聚焦角色**：
```
GET /api/v1/role/{id}
```

**角色数据结构**：
```json
{
  "id": "role:admin",
  "name": "管理员",
  "code": "admin",
  "description": "系统管理员，拥有所有权限",
  "created_at": "2024-01-01T00:00:00Z",
  "permissions:": [
    "permission:user:create",
    "permission:user:read",
    "permission:user:update",
    "permission:user:delete",
    "permission:role:manage"
  ],
  "permissions-fold": true,
  "users:": ["user:1", "user:3"],
  "users-fold": true,
  "commands": [
    {
      "action": "add_permission",
      "label": "添加权限",
      "params": [
        {
          "name": "permission_id",
          "type": "entity",
          "template": "permission",
          "required": true,
          "label": "权限"
        }
      ],
      "trigger": "/api/v1/role/admin/commands",
      "sync": true
    },
    {
      "action": "remove_permission",
      "label": "移除权限",
      "params": [
        {
          "name": "permission_id",
          "type": "entity",
          "template": "permission",
          "required": true,
          "label": "权限"
        }
      ],
      "trigger": "/api/v1/role/admin/commands",
      "sync": true
    },
    {
      "action": "clone",
      "label": "复制角色",
      "params": [
        {
          "name": "new_name",
          "type": "string",
          "required": true,
          "label": "新角色名称"
        }
      ],
      "trigger": "/api/v1/role/admin/commands",
      "sync": true
    }
  ]
}
```

### 3.3 权限（Permission）

**列出权限**：
```
GET /api/v1/permission
```

**聚焦权限**：
```
GET /api/v1/permission/{id}
```

**权限数据结构**：
```json
{
  "id": "permission:user:create",
  "name": "创建用户",
  "code": "user:create",
  "resource": "user",
  "action": "create",
  "description": "允许创建新用户",
  "roles:": ["role:admin", "role:manager"],
  "roles-fold": true,
  "commands": []
}
```

---



---

## 七、实体模板列表

| 模板 | 路径 | 说明 |
|------|------|------|
| `user` | `/api/v1/user` | 用户 |
| `role` | `/api/v1/role` | 角色 |
| `permission` | `/api/v1/permission` | 权限 |

---

## 八、意图字段规范

### 8.1 解引用意图（fold）

控制实体引用默认展开/折叠：

```json
{
  "author:": "user:1",
  "author-fold": false
}
```

- `false`：默认展开（显示嵌套实体详情）
- `true` 或缺失：默认折叠（显示为可点击标签）

### 8.2 执行策略意图（sync）

控制命令执行方式：

```json
{
  "action": "generate_report",
  "sync": false
}
```

- `true`（默认）：同步执行，等待响应
- `false`：异步执行，返回任务 ID，通过 SSE 推送结果

---

## 九、响应状态码

| 状态码 | 含义 |
|--------|------|
| 200 | 请求成功（同步命令） |
| 202 | 已接受（异步命令） |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 实体不存在 |
| 409 | 业务冲突 |
| 422 | 验证失败 |
| 500 | 服务器内部错误 |

---

## 十、错误响应格式

```json
{
  "error": {
    "code": "ENTITY_NOT_FOUND",
    "message": "指定的实体不存在",
    "details": {
      "template": "user",
      "id": "999"
    }
  }
}
```

---

## 十一、版本化设计

- **当前版本**：`v1`
- **版本前缀**：`/api/v1`
- **向前兼容**：v2 开发时保留 v1 端点
- **版本协商**：通过 URL 路径明确版本，不支持通过 Header 协商

---

*文档版本：v1.0*
*归档时间：2026-05-02*
