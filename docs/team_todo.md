# Team TODO (汇总)

下面列出原来每个 worker 文件中的短期 TODO，集中管理，便于展示后跟进实现。

## Worker1 (课程域)
1. 修复/确认 SQL 去重问题：复现重复 case，修改 `QueryStr` 或在 repo 层做去重，保证按 `course_num` 唯一。
2. 把 `GetStructuredCoursesHandler` 的逻辑改为复用 `SearchByAreaAndWeekday` 或把查询合并为一条 SQL（降低重复与 bug 面）。
3. 给 `pkg/generator` 的边界条件写更全面的测试（包含非法值、极端周/节值）。
4. 为课程相关的服务编写单元/集成测试（建议使用 sqlite 内存 DB 做集成测试）。
5. 文档化二进制位分布（周/节位序），方便新同学理解。

## Worker2 (社区域)
1. 实现“活跃用户”接口：SQL 聚合（WHERE created_at >= now()-3days GROUP BY user_id ORDER BY COUNT DESC LIMIT N）并支持缓存（Redis）。
2. 实现“社区数据”接口（总帖子、注册用户、今日新帖），并把结果缓存 60s 或更长，避免频繁扫描大表。
3. 为用户统计（posts/comments/likes/likesReceived）添加 API：推荐做两步：a) 尝试先做事实计算 SQL；b) 如果频繁，增加计数字段并在写操作（发帖/评论/点赞）维护计数（事务或异步任务队列）。
4. 补充单元测试与集成测试（post/comment 服务），覆盖关键业务路径（创建、删除、统计）。

## Worker3 (平台/Infra)
1. 将配置初始化由隐式 `init()` 改为显式调用（可在 `cmd/main.go` 或测试 setup 中调用），避免测试污染与不可预期的 side-effect。
2. 为 JWT key、DB 密码建立更安全的存放/注入（环境变量、CI secret），并移除日志中的明文打印。
3. 提供一个简单的 `sqlite` 内存配置用于单元测试（让 services 的测试无需依赖真实 MySQL）。
4. 为 Redis 添加 wrapper 工具（Get/Set JSON with type-safe helpers），并提供典型缓存 primitives（GetOrSet、Invalidate）。
5. 在 CI 中添加基础的静态检查（go vet / golangci-lint），并在 Pull Request 流程中强制运行 `go test ./...`。

---

建议流程：
- 在下次演示前（Sprint Demo），把 `team_todo.md` 中排前 1-2 项优先实现为小 PR（每项配套测试与演示脚本）。
- 在合并业务逻辑前，Worker3 负责做一次集成/安全/配置审核。
