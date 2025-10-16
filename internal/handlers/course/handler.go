package course

import (
	"cengkeHelperBackGo/internal/config"
	// database "cengkeHelperBackGo/internal/db" // Handler 不再直接使用 database.Client
	"cengkeHelperBackGo/internal/models/dto"
	"cengkeHelperBackGo/internal/models/vo"
	"cengkeHelperBackGo/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CourseHandler 处理课程相关的HTTP请求
type CourseHandler struct {
	courseService          *services.CourseService
	courseStructureService *services.CourseStructureService
}

// NewCourseHandler 创建一个新的 CourseHandler
func NewCourseHandler() *CourseHandler {
	return &CourseHandler{
		courseService:          services.NewCourseService(),
		courseStructureService: services.NewCourseStructureService(),
	}
}

// 从 context 中获取用户 ID 的辅助函数 (与 comment_handler 类似)
// 请确保此实现与您的认证中间件如何存储用户ID的方式一致
func getCourseHandlerUserIDFromContext(c *gin.Context) (*uint32, bool) {
	userIDVal, exists := c.Get("userId") // 假设认证中间件将 userID 存入 context，键名为 "userID"
	if !exists {
		return nil, false
	}
	parseUint, err := strconv.ParseUint(userIDVal.(string), 10, 32)
	if err != nil {
		return nil, false
	}
	userId := uint32(parseUint)
	return &userId, true
}

// GetCoursesHandler godoc
// @Summary 获取所有课程列表 (按学部和教学楼分组)
// @Description 获取所有课程的列表，按学部和教学楼进行分组。
// @Tags Courses
// @Accept json
// @Produce json
// @Success 200 {object} vo.RespData{data=[][]vo.BuildingInfoVO} "成功"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /courses [get]
func (h *CourseHandler) GetCoursesHandler(c *gin.Context) {
	// 不再需要原始 handler 中的 ToVO 和 convertCoursesToVO 辅助函数，
	// 因为数据转换的逻辑移到了 service 层。

	infos := GetTeachInfos()
	// 假设 vo.RespondSuccess 存在
	vo.RespondSuccess(c, "课程数据获取成功", infos)
}

// GetStructuredCoursesHandler godoc
// @Summary 获取结构化的课程数据（学部 → 教学楼 → 楼层 → 课程）
// @Description 获取按照四级结构组织的课程数据。默认返回当前时间的课程。参数说明：-1表示不限（查询所有），0或不传表示使用当前时间
// @Tags Courses
// @Accept json
// @Produce json
// @Param weekNum query int false "周次（-1=不限, 0或不传=当前周次）"
// @Param weekday query int false "星期几（-1=不限, 0或不传=当前星期）"
// @Param lessonNum query int false "节次（-1=不限, 0或不传=当前节次）"
// @Param divisionId query int false "学部ID（1-4，不传表示所有学部）"
// @Param useCache query bool false "是否使用缓存（默认true）"
// @Success 200 {object} vo.RespData{data=[]vo.DivisionVO} "成功"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /courses/structured [get]
func (h *CourseHandler) GetStructuredCoursesHandler(c *gin.Context) {
	// 解析查询参数，默认值为0（表示使用当前时间）
	params := &services.CourseQueryParams{
		WeekNum:   0, // 0 表示使用当前周次
		Weekday:   0, // 0 表示使用当前星期
		LessonNum: 0, // 0 表示使用当前节次
		UseCache:  true,
	}

	// 解析 weekNum
	if weekNumStr := c.Query("weekNum"); weekNumStr != "" {
		if weekNum, err := strconv.Atoi(weekNumStr); err == nil {
			params.WeekNum = weekNum
		}
	}

	// 解析 weekday
	if weekdayStr := c.Query("weekday"); weekdayStr != "" {
		if weekday, err := strconv.Atoi(weekdayStr); err == nil {
			params.Weekday = weekday
		}
	}

	// 解析 lessonNum
	if lessonNumStr := c.Query("lessonNum"); lessonNumStr != "" {
		if lessonNum, err := strconv.Atoi(lessonNumStr); err == nil {
			params.LessonNum = lessonNum
		}
	}

	// 解析 divisionId
	if divisionIDStr := c.Query("divisionId"); divisionIDStr != "" {
		if divisionID, err := strconv.Atoi(divisionIDStr); err == nil {
			params.DivisionID = &divisionID
		}
	}

	// 解析 useCache
	if useCacheStr := c.Query("useCache"); useCacheStr != "" {
		if useCache, err := strconv.ParseBool(useCacheStr); err == nil {
			params.UseCache = useCache
		}
	}

	divisions, err := h.courseStructureService.GetStructuredCourses(params)
	if err != nil {
		vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "获取课程数据失败", err)
		return
	}
	vo.RespondSuccess(c, "课程数据获取成功", divisions)
}

// GetCourseDetailHandler godoc
// @Summary 获取课程详情
// @Description 根据课程ID获取特定课程的详细信息。
// @Tags Courses
// @Accept json
// @Produce json
// @Param courseId path uint true "课程ID"
// @Success 200 {object} vo.RespData{data=vo.CourseDetailVO} "成功"
// @Failure 400 {object} vo.RespData "请求参数错误 (无效的课程ID)"
// @Failure 404 {object} vo.RespData "课程未找到"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /courses/{courseId} [get]
func (h *CourseHandler) GetCourseDetailHandler(c *gin.Context) {
	courseIDStr := c.Param("courseId")
	courseIDUint64, err := strconv.ParseUint(courseIDStr, 10, 32) // 解析为 uint64 再转 uint
	if err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeInvalidParams, "无效的课程ID格式", err)
		return
	}
	courseID := uint(courseIDUint64)

	courseDetail, serviceErr := h.courseService.GetCourseDetailByID(courseID)
	if serviceErr != nil {
		if serviceErr.Error() == config.MsgCourseNotFound { // 使用统一定义的错误消息
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, serviceErr.Error(), nil)
		} else {
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "获取课程详情处理失败", serviceErr)
		}
		return
	}
	vo.RespondSuccess(c, "课程详情获取成功", courseDetail)
}

// GetCourseReviewsHandler godoc
// @Summary 获取课程的评价列表
// @Description 根据课程ID获取特定课程的所有评价。
// @Tags Courses
// @Accept json
// @Produce json
// @Param courseId path uint true "课程ID"
// @Success 200 {object} vo.RespData{data=[]vo.CourseReviewInfoVO} "成功"
// @Failure 400 {object} vo.RespData "请求参数错误 (无效的课程ID)"
// @Failure 404 {object} vo.RespData "课程未找到"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /courses/{courseId}/reviews [get]
func (h *CourseHandler) GetCourseReviewsHandler(c *gin.Context) {
	courseIDStr := c.Param("courseId")
	courseIDUint64, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeInvalidParams, "无效的课程ID格式", err)
		return
	}
	courseID := uint(courseIDUint64)

	reviews, serviceErr := h.courseService.GetCourseReviewsByCourseID(courseID)
	if serviceErr != nil {
		if serviceErr.Error() == config.MsgCourseNotFound {
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, serviceErr.Error(), nil)
		} else {
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "获取课程评价列表失败", serviceErr)
		}
		return
	}

	if reviews == nil { // 即使 service 返回 nil，也确保前端收到空数组
		reviews = []vo.CourseReviewInfoVO{}
	}
	vo.RespondSuccess(c, "课程评价列表获取成功", reviews)
}

// SubmitCourseReviewHandler godoc
// @Summary 提交课程评价
// @Description 为指定课程提交一条新的评价。需要用户认证。
// @Tags Courses
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param reviewData body dto.CourseReviewCreateDTO true "评价数据"
// @Success 201 {object} vo.RespData{data=vo.CourseReviewInfoVO} "评价提交成功"
// @Failure 400 {object} vo.RespData "请求参数错误"
// @Failure 401 {object} vo.RespData "用户未授权"
// @Failure 404 {object} vo.RespData "相关资源未找到 (课程/用户)"
// @Failure 500 {object} vo.RespData "服务器内部错误"
// @Router /reviews [post]
func (h *CourseHandler) SubmitCourseReviewHandler(c *gin.Context) {
	userIDStr := c.GetString("userId")
	userIDVal, err := strconv.ParseUint(userIDStr, 10, 64)
	if userIDStr == "" || err != nil {
		vo.RespondError(c, http.StatusUnauthorized, config.CodeUnauthorized, "用户未授权或无法获取用户ID", nil)
		return
	}
	userID := uint32(userIDVal)
	var payload dto.CourseReviewCreateDTO // 使用 dto.CourseReviewCreateDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		vo.RespondError(c, http.StatusBadRequest, config.CodeInvalidParams, "请求参数无效", err)
		return
	}
	// courseStore.ts 中已将 courseId 转为 Number，后端 payload.CourseID 是 uint

	createdReviewVO, serviceErr := h.courseService.SubmitCourseReview(userID, payload)
	if serviceErr != nil {
		errMsg := serviceErr.Error()
		switch {
		case errMsg == config.MsgCourseForReviewNotFound || errMsg == config.MsgUserForReviewNotFound:
			vo.RespondError(c, http.StatusNotFound, config.CodeNotFound, errMsg, nil)
		default: // 其他所有来自 service 层的错误都视为内部错误
			vo.RespondError(c, http.StatusInternalServerError, config.CodeServerError, "提交评价处理失败", serviceErr)
		}
		return
	}

	// HTTP 201 Created 表示资源成功创建，并返回创建的资源
	c.JSON(http.StatusCreated, vo.NewSuccessResp("评价提交成功", createdReviewVO))
}
