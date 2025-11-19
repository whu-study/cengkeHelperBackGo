# Worker 3 — Infra, Auth, Router & Platform

Primary focus
- 负责平台基础设施：数据库与 Redis 连接、项目配置加载、认证中间件（JWT）、路由与全局中间件、邮件服务和 CI/部署脚本。

主要负责的代码文件/模块（建议归属）
- internal/db/
  - `mysql.go`, `mysql_test.go`, `redis.go` (DB/Redis 连接与客户端包装)
- internal/config/
  - `config.go`, `constant.go` (配置加载与热重载)
- internal/filter/
  - `auth_filter.go` (认证中间件)
- internal/handlers/auth/
  - `handler.go`, `service.go` (登录/注册/鉴权 handler)
- internal/router/
  - `router.go` (注册所有路由及中间件)
- pkg/utils/
  - `jwt.go` (JWT 生成与解析)
- cmd/
  - `main.go` (程序入口，服务启动)
- internal/services/
  - `email_service.go` (邮件发送)

已完成 / 已实现（现状说明）
- 配置加载（`internal/config/config.go`）实现了从 `config.yaml` 加载与热重载（fsnotify）。
- 初始化并导出了 `database.Client`（MySQL）及 `redis` 客户端（`internal/db/*`）。
- 实现了 JWT 生成/解析逻辑（`pkg/utils/jwt.go`），并在认证中间件中消费，并注入 userID 到请求上下文。
- 提供了路由注册点（`internal/router/router.go`），并在 `cmd/main.go` 中启动 HTTP 服务。

存在的问题 / 风险点（便于展示时说明）
- 配置热重载当前做法在测试/CI 环境下可能会打印敏感信息到 stdout（注意 jwt key、数据库密码），需要在日志中遮掩或限定环境。
- `internal/config/config.go` 中的 `init()` 会在测试时自动触发配置加载，导致测试环境输出噪音或依赖本地文件。建议改为在 main 或测试初始化显式调用 `LoadConfig`。
- 需要为 DB/Redis 的错误与重连策略加更多健壮性（超时、断线重试、连接池参数）。
- 缓存方案（Redis）需要和 Worker2/Worker1 对齐键与过期策略。

协作与交付
- Worker3 负责 CI/CD 管线和基础库代码合并，其他两位完成业务功能后提交 PR，由 Worker3 负责合并前的安全/配置审查。
- 对关键变更（DB schema、缓存策略）先在 docs 中写明迁移步骤与回滚策略。
