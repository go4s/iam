# IAM API 对接指南（v1）

## 快速开始

### 1. 基础信息

- **基础 URL**: `http://localhost:8080/api/v1`
- **内容类型**: `application/json`
- **认证方式**: `Authorization: Bearer {token}`

### 2. 第一个请求

```bash
# 登录获取 Token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

**响应**:
```json
{
  "code": "0000",
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "user:1",
      "username": "admin"
    }
  }
}
```

```bash
# 使用 Token 访问接口
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

---

## 认证流程

### 登录

```http
POST /api/v1/auth/login
```

**请求体**:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**成功响应**:
```json
{
  "code": "0000",
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "user:1",
      "username": "admin"
    }
  }
}
```

### 获取当前用户

```http
GET /api/v1/auth/me
Authorization: Bearer {token}
```

**响应**:
```json
{
  "code": "0000",
  "message": "success",
  "data": {
    "id": "user:1",
    "username": "admin",
    "roles:": ["role:admin"],
    "roles-fold": false
  }
}
```

---

## 通用响应格式

### 成功响应

```json
{
  "code": "0000",
  "message": "success",
  "data": { ... }
}
```

### 错误响应

```json
{
  "code": "1001",
  "message": "实体不存在",
  "data": {
    "template": "user",
    "id": "999"
  }
}
```

### 业务码说明

| 业务码 | HTTP 状态 | 含义 |
|--------|-----------|------|
| `0000` | 200 | 请求成功 |
| `1001` | 404 | 实体不存在 |
| `1002` | 400 | 参数错误 |
| `1003` | 401 | 未认证（Token 无效或过期）|
| `1004` | 403 | 无权限 |
| `1005` | 409 | 业务冲突 |
| `1006` | 422 | 验证失败 |
| `9999` | 500 | 服务器内部错误 |

---

## 实体操作模式

所有实体（user、role、permission）遵循统一的三种操作模式：

### 1. 列出（List）

```http
GET /api/v1/{entity}?page=1&size=20&q=keyword
```

**响应**:
```json
{
  "code": "0000",
  "message": "success",
  "data": {
    "items": [
      {
        "id": "user:1",
        "username": "admin"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 20,
      "total": 100
    }
  }
}
```

**查询参数**:
- `page`: 页码，默认 1
- `size`: 每页数量，默认 20，最大 100
- `q`: 搜索关键词（可选）

### 2. 聚焦（Focus）

```http
GET /api/v1/{entity}/{id}
```

**说明**:
- URL 中的 `{id}` 使用纯数字，如 `/api/v1/user/1`
- 响应中的 `id` 字段使用完整格式，如 `"user:1"`
- 返回完整实体对象，包含 `commands` 数组
- 引用字段附带 `{ref}-fold` 意图字段

**响应示例**:
```json
{
  "code": "0000",
  "message": "success",
  "data": {
    "id": "user:1",
    "username": "admin",
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
      }
    ]
  }
}
```

### 3. 执行命令（Command）

```http
POST /api/v1/{entity}/{id}/commands
```

**请求体**:
```json
{
  "action": "reset_password",
  "params": {
    "new_password": "newpass123"
  }
}
```

**响应**:
```json
{
  "code": "0000",
  "message": "密码重置成功",
  "data": {
    "message": "密码重置成功"
  }
}
```

---

## 模块接口详解

### 用户管理

#### 列出用户

```http
GET /api/v1/user?page=1&size=20&q=admin
Authorization: Bearer {token}
```

#### 获取用户详情

```http
GET /api/v1/user/1
Authorization: Bearer {token}
```

#### 用户命令

| 命令 | 说明 | 参数 |
|------|------|------|
| `reset_password` | 重置密码 | `new_password` (string, 必填) |
| `assign_role` | 分配角色 | `role_id` (number, 必填) |
| `disable` | 禁用账户 | 无 |
| `create_user` | 创建用户（admin 专属）| `username`, `password`, `role_ids` |

**创建用户示例**:
```http
POST /api/v1/user/1/commands
Authorization: Bearer {token}
Content-Type: application/json

{
  "action": "create_user",
  "params": {
    "username": "newuser",
    "password": "password123",
    "role_ids": [1]
  }
}
```

---

### 角色管理

#### 列出角色

```http
GET /api/v1/role?page=1&size=20
Authorization: Bearer {token}
```

#### 获取角色详情

```http
GET /api/v1/role/1
Authorization: Bearer {token}
```

#### 角色命令

| 命令 | 说明 | 参数 |
|------|------|------|
| `add_permission` | 添加权限 | `permission_id` (number, 必填) |
| `remove_permission` | 移除权限 | `permission_id` (number, 必填) |
| `clone` | 复制角色 | `new_name` (string, 必填) |

---

### 权限管理

#### 列出权限

```http
GET /api/v1/permission?page=1&size=20
Authorization: Bearer {token}
```

#### 获取权限详情

```http
GET /api/v1/permission/1
Authorization: Bearer {token}
```

**说明**: Permission 实体不支持命令操作，`commands` 数组始终为空。

---

## 实体 ID 格式说明

### URL 路径 vs 响应体

| 场景 | 格式 | 示例 |
|------|------|------|
| URL 路径参数 | 纯数字 | `/api/v1/user/1` |
| 响应体 `id` 字段 | `{template}:{id}` | `"user:1"` |
| 响应体引用字段 | `{template}:{code}` | `"role:admin"` |

### 引用字段与意图

实体引用字段以 `:` 结尾，附带 `{field}-fold` 布尔值控制展开/折叠：

```json
{
  "roles:": ["role:admin", "role:editor"],
  "roles-fold": false
}
```

- `fold: false` — 默认展开显示
- `fold: true` — 默认折叠显示（可点击标签）

---

## Command 模式详解

### 命令结构

```json
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
}
```

### 参数类型

| 类型 | 说明 | 示例 |
|------|------|------|
| `string` | 字符串 | `"new_password"` |
| `entity` | 实体引用（单选）| `{"template": "role"}` |
| `entities` | 实体引用（多选）| `{"template": "role"}` |

### 执行策略

- `sync: true` — 同步执行，立即返回结果（HTTP 200）
- `sync: false` — 异步执行，返回任务 ID（HTTP 202）**[当前版本暂未实现]**

---

## 前端对接最佳实践

### 1. 动态渲染表单

从 Focus 接口获取 `commands` 数组，动态生成操作按钮和表单：

```javascript
// 示例：根据 commands 生成操作按钮
const commands = userData.commands;
commands.forEach(cmd => {
  const button = document.createElement('button');
  button.textContent = cmd.label;
  button.onclick = () => executeCommand(cmd);
});
```

### 2. 统一错误处理

```javascript
async function apiRequest(url, options = {}) {
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      ...options.headers
    }
  });
  
  const data = await response.json();
  
  if (data.code !== '0000') {
    throw new Error(`[${data.code}] ${data.message}`);
  }
  
  return data.data;
}
```

### 3. Token 管理

```javascript
// 登录时存储 Token
const { token, user } = await login(credentials);
localStorage.setItem('token', token);

// 请求时携带 Token
const token = localStorage.getItem('token');
const userData = await apiRequest('/api/v1/auth/me', {
  headers: { 'Authorization': `Bearer ${token}` }
});

// Token 过期处理（401）
if (response.code === '1003') {
  localStorage.removeItem('token');
  window.location.href = '/login';
}
```

### 4. 分页处理

```javascript
async function fetchUsers(page = 1, size = 20) {
  const data = await apiRequest(`/api/v1/user?page=${page}&size=${size}`);
  return {
    items: data.items,
    pagination: data.pagination
  };
}
```

---

## 系统接口

### 重载格式配置

```http
POST /api/v1/system/reload-formats
Authorization: Bearer {token}
```

**说明**: 当 `entity_format` 表数据变更后，调用此接口刷新内存缓存。**仅限 admin 角色**。

---

## 常见问题

### Q: 为什么 URL 用纯数字，响应体用 `user:1`？

A: URL 遵循 RESTful 惯例，路径前缀 `/user/` 已表明实体类型。响应体使用完整格式是为了全局唯一性，方便跨实体引用。

### Q: `commands` 数组从哪里来？

A: 通过 Focus 接口获取（`GET /api/v1/user/1`）。Commands 根据当前用户角色动态生成，admin 用户会看到更多操作。

### Q: 如何区分同步和异步命令？

A: 查看 command 对象的 `sync` 字段。当前版本所有命令均为同步（`sync: true`），异步模式将在后续版本支持。

### Q: 实体格式配置可以自定义吗？

A: 可以通过修改 `entity_format` 表调整字段显示。修改后调用 `POST /api/v1/system/reload-formats` 生效。

---

## 实体模板列表

| 模板 | 路径 | 说明 |
|------|------|------|
| `user` | `/api/v1/user` | 用户管理 |
| `role` | `/api/v1/role` | 角色管理 |
| `permission` | `/api/v1/permission` | 权限管理 |

---

*文档版本：v1.0*
*更新时间：2026-05-02*
