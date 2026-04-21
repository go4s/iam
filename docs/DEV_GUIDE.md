# 业务开发对接手册

本手册为业务方（网关、微服务）提供集成 IAM 服务的技术规范和接口说明。

## 1. 快速接入流程

IAM 服务提供两种维度的能力：
1.  **认证 (Authentication)**: 验证用户提供的 Token 是否有效。
2.  **鉴权 (Authorization)**: 判断特定用户是否有权访问特定资源。

---

## 2. 身份认证对接 (针对网关/前置服务)

外部网关在收到客户端请求后，应提取 `Authorization` Header 中的 Token 并发送至 IAM 进行验证。

### 2.1 远程验证接口
-   **Endpoint**: `POST /auth/validate`
-   **请求体**:
    ```json
    { "token": "eyJhbGciOiJIUzI1..." }
    ```
-   **成功返回 (200 OK)**:
    ```json
    {
      "valid": true,
      "username": "zhangsan",
      "role": "Editor"
    }
    ```
-   **失败返回 (401 Unauthorized)**:
    ```json
    { "error": "Invalid token" }
    ```

### 2.2 本地验证 (进阶)
网关如果性能要求极高，可以采用本地验证 JWT。
-   **算法**: HS256
-   **密钥 (Secret)**: 由环境变量 `JWT_SECRET` 指定（默认值为 `your-secret-key`）
-   **Payload 结构**:
    -   `sub`: 用户名
    -   `role`: 角色名
    -   `exp`: 过期时间戳

---

## 3. 鉴权对接 (针对业务微服务)

当业务服务需要针对特定操作进行精确权限判断时，可调用此接口。

### 3.1 远程鉴权接口
-   **Endpoint**: `POST /api/v1/enforce`
-   **注意**: 此接口受 JWT 保护，调用方需在 Header 携带有效的管理/服务 Token。
-   **请求体**:
    ```json
    {
      "sub": "editor1",           // 用户名或角色
      "obj": "/api/v1/projects",  // 待访问的资源路径
      "act": "POST"               // 操作类型 (GET, POST, DELETE 等)
    }
    ```
-   **返回结果**:
    ```json
    { "allowed": true }
    ```

---

## 4. API 异常处理规范

| 状态码 | 含义 | 建议操作 |
| :--- | :--- | :--- |
| **401 Unauthorized** | Token 缺失、格式错误或已过期 | 引导用户重新登录 |
| **403 Forbidden** | 用户认证成功，但无权访问该资源 | 提示“无操作权限” |
| **400 Bad Request** | 请求参数错误（如 JSON 格式不对） | 检查代码实现 |
| **500 Internal Error** | IAM 服务内部故障 | 联系系统管理员 |

---

## 5. 开发建议

1.  **缓存机制**：业务服务可以对常用的鉴权结果（如 `sub+obj+act`）进行短时间缓存（如 1 分钟），以减轻 IAM 服务的压力。
2.  **Token 传递**：微服务间调用时，建议在 Header 中透明传递原始 `Authorization` Token。
3.  **环境隔离**：确保开发环境和生产环境使用不同的 JWT Secret。
