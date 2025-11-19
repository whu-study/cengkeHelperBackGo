# Worker 2 — Forum & Social Features

Primary focus
- 负责社区（论坛）模块：帖子（Post）与评论（Comment）的 handler、service、repo（如果存在）、以及与用户互动相关的逻辑（点赞、收藏、计数）。

主要负责的代码文件/模块（建议归属）
- internal/handlers/
  - `post_handler.go` (帖子相关 HTTP handler)
  - `comment_handler.go` (评论相关 handler)
  - `users.go` (用户展示/用户相关接口)
- internal/services/
  - `post_sevice.go` (注意目录名拼写，服务实现)
  - `comment_service.go`
  - `user.go` (用户相关服务)
- internal/models/dto/
  - `post.go`, `comment_dto.go`, `like_favourite_dto.go`
- internal/models/vo/
  - `post_vo.go`, `comment_vo.go`, `user_vo.go`

已完成 / 已实现（现状说明）
- 实现了帖子发布、查询、分页、按用户/话题过滤等基本接口。
- 实现了对帖子的评论、回复、点赞/收藏等功能的数据模型与 handler/服务调用路径。
- 在用户 VO/DTO 中保留了 `postsCount`、`commentsCount`、`likesCount` 等字段（部分为即时计算）。

存在的问题 / 风险点（便于展示时说明）
- 活跃用户（近 3 天发帖数 Top N）接口尚未实现（可作为后续展示的扩展）。
- 社区统计接口（总帖子、注册用户、今日新帖）和用户收到的点赞数统计仍可优化（缓存或计数字段）。

协作建议
- 与 Worker3 协调 Redis 缓存键命名与过期策略（例如 `community:stats`, `user:active:3d`）。
- 与 Worker1 对“今日课程统计”与“当前时段课程数”接口进行对接（Worker1 提供 `GetSingleNumOfCourses`、`GetOneDayNumOfCourses` 方法）。

验收标准（演示时可说明）
- 已实现的帖子、评论、点赞等功能完整可用，并包含基本的接口文档与示例。
