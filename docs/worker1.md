# Worker 1 — Courses & Scheduling

Primary focus
- 负责课程相关后端逻辑：课程数据建模、课程时间/周次二进制逻辑、按学部/教学楼的课程结构化接口、课程查询 SQL 和服务层实现。

主要负责的代码文件/模块（建议归属）
- internal/handlers/course/
  - `building.go` (构建学部/教学楼的 handler/逻辑)
  - `handler.go` (课程相关 HTTP handler)
  - `time.go`, `new_course.go`, `const.go`
- internal/services/
  - `course_service.go` (课程服务：查询、GetAllCourses、课程详情/评价的业务逻辑)
  - `course_structure_service.go` (课程结构相关逻辑)
  - `course/status.go`, `course/tools.go`
- internal/repo/
  - `course_repo.go` (原生 SQL：SearchByAreaAndWeekday 等)
- internal/models/dto/
  - `course_helper.go` (CourseInfo / TimeInfo / CourseReview DTO & Model)
- internal/models/vo/
  - `course_info_vo.go`, `course_info_vo.go` 等 VO 定义
- pkg/generator/
  - `binary.go` (week/lesson 二进制编码/解码与显示逻辑)

已完成 / 已实现（现状说明）
- 实现了周次/节次到二进制编码与反向解码（`WeekLesson2Bin`, `Bin2WeekLesson`, `IsWeekLessonMatch` 等），并提供 `NearestToDisplay` 用于前端显示。
- 编写了按教学楼 + course_num 去重的 SQL（见 `course_repo.go` 中的 QueryStr），用于返回代表行。
- 在 `course_service.go` 中实现了 `GetAllCourses` 等服务方法，用于按学部/教学楼聚合课程并输出 VO。
- 在 DTO 层定义了课程评价相关模型 `CourseReviewModel` 与 `CourseReviewCreateDTO`（最近已把 DTO 的 `rating` 从 int 改为 float32，以兼容前端传 4.5 的场景）。
- 为 `pkg/generator` 编写了几条单元测试（见 `pkg/generator/binary_test.go` 与新增的单元测试文件）。

存在的问题 / 风险点（便于展示时说明）
- SQL 查询可能返回重复（需要按 course_num 去重或确保聚合字段一致）。
- `GetStructuredCoursesHandler` 与 `handlers/course/building` 的逻辑存在重复，实现应合并或复用服务层查询。
- 二进制位序（高19位为周次、低13位为节次）的假设需文档化并补充边界单测。

联系人与约定
- 负责 Pull Request 的 reviewer（课程域）优先由 Worker1 自己审查一次后交给 Worker3 做 CI/安全检查。
- 所有 SQL 变更必须包含复现用例或测试。

