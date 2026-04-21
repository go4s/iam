# IAM 系统管理员使用手册

本手册旨在指导系统管理员如何使用 IAM (Identity and Access Management) 服务进行权限配置、角色管理和资源管理。

## 1. 核心概念

本系统基于 **Casbin RBAC (基于角色的访问控制)** 模型实现。

-   **Subject (sub)**: 主体。可以是用户名（如 `user1`）或角色名（如 `Admin`, `Editor`）。
-   **Object (obj)**: 资源对象。通常是 API 路径（如 `/api/v1/articles`）。
-   **Action (act)**: 操作。如 `GET`, `POST`, `DELETE`。
-   **Role (g)**: 角色继承关系。定义谁属于哪个角色。

### 特殊规则
-   **超级管理员**: 拥有 `Admin` 角色的用户自动拥有系统内**所有资源**的操作权限，无需显式配置策略。

---

## 2. 账号与认证

### 2.1 注册管理员/用户
通过 `/auth/register` 接口创建账号。注册时需指定初始角色。
-   **接口**: `POST /auth/register`
-   **Payload**:
    ```json
    {
      "username": "admin_user",
      "password": "secure_password",
      "role": "Admin"
    }
    ```

### 2.2 登录获取 Token
后续所有管理接口均需在 Header 中携带 JWT Token。
-   **接口**: `POST /auth/login`
-   **Header**: `Authorization: Bearer <token>`

---

## 3. 权限策略管理 (Policy)

策略定义了“谁可以对什么资源做什么”。

### 3.1 添加策略
-   **接口**: `POST /api/v1/policies`
-   **示例**: 让 `Editor` 角色可以发布文章。
    ```json
    {
      "sub": "Editor",
      "obj": "/api/v1/articles",
      "act": "POST"
    }
    ```

### 3.2 删除策略
-   **接口**: `DELETE /api/v1/policies`
-   **Payload**: 需提供完整的 sub, obj, act 匹配项。

### 3.3 查询策略
-   **接口**: `GET /api/v1/policies`
-   **返回**: 当前系统中所有的权限规则列表。

---

## 4. 角色与继承管理 (Grouping)

### 4.1 角色继承
支持角色间的层级关系。例如，让 `Admin` 继承 `Editor` 的所有权限。
-   **接口**: `POST /api/v1/policies/grouping`
-   **Payload**:
    ```json
    {
      "sub": "Admin",
      "role": "Editor"
    }
    ```

### 4.2 查看继承关系
-   **接口**: `GET /api/v1/policies/grouping`

---

## 5. 资源管理实战

### 场景：将一个新资源（如 `projects`）纳入管理
1.  **明确路径**：确定新资源的 API 路径，如 `/api/v1/projects`。
2.  **分配基础权限**：调用 `POST /api/v1/policies`。
    -   `sub`: `Editor`, `obj`: `/api/v1/projects`, `act`: `POST`
3.  **用户关联**：如果一个新用户需要管理项目，在注册时赋予其 `Editor` 角色，或通过 `grouping` 接口将该用户关联到 `Editor` 角色。

---

## 6. 管理建议
-   **最小权限原则**：优先给角色分配权限，再将用户关联到角色，避免直接给用户（sub 为用户名）分配权限。
-   **路径规范**：`obj` 建议使用完整的 API 路径，以确保权限校验的准确性。
