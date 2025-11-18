package services

import (
	"cengkeHelperBackGo/internal/config"      // 假设config包路径正确
	database "cengkeHelperBackGo/internal/db" // 确保这是您项目中统一的数据库客户端包
	"cengkeHelperBackGo/internal/models/dto"  // 使用您提供的dto包
	"cengkeHelperBackGo/internal/models/vo"   // 使用您提供的vo包
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"math"
	"slices"
	"strconv"
)

// CourseService 结构体用于组织课程相关的服务方法
type CourseService struct {
	// 此处不需要 db *gorm.DB 字段, 因为我们将直接使用 database.Client
}

// NewCourseService 创建 CourseService 实例
func NewCourseService() *CourseService {
	return &CourseService{}
}

// convertCourseModelToInfoVO 将 dto.CourseInfo (模型) 转换为 vo.CourseInfoVO
// 注意：原始 GetCoursesHandler 中的 convertCoursesToVO 参数是 dto.CourseInfo，返回 vo.CourseInfoVO
// 我们在这里保持这个转换逻辑，但确保字段映射正确。
// dto.CourseInfo 包含 ID，而原始 convertCoursesToVO 的输出 vo.CourseInfoVO 没有ID，但前端 course.ts 定义了 id。
// 我将假设 vo.CourseInfoVO 也需要 ID。
//
//	func convertCourseModelToInfoVO(course dto.CourseInfo) vo.CourseInfoVO {
//		//return vo.CourseInfoVO{
//		//	ID:           course.ID, // 确保 vo.CourseInfoVO 有 ID 字段
//		//	Room:         course.Room,
//		//	Faculty:      course.Faculty,
//		//	CourseName:   course.CourseName,
//		//	TeacherName:  course.TeacherName,
//		//	TeacherTitle: course.TeacherTitle,
//		//	CourseTime:   course.CourseTime,
//		//	CourseType:   course.CourseType,
//		//}
//	}
//
// courseQueryRow 用于接收 Raw SQL 查询的聚合结果
// 这基于 const.go 中的 queryStr 的 SELECT 字段
type courseQueryRow struct {
	ID           uint32 `gorm:"column:id"`
	CourseNum    string `gorm:"column:course_num"`
	Classroom    string `gorm:"column:classroom"`
	Building     string `gorm:"column:building"`
	CourseType   string `gorm:"column:course_type"`
	Faculty      string `gorm:"column:faculty"`
	CourseName   string `gorm:"column:course_name"`
	Teacher      string `gorm:"column:teacher"`
	TeacherTitle string `gorm:"column:teacher_title"`
	WeekAndTime  uint32 `gorm:"column:week_and_time"`
	DayOfWeek    uint8  `gorm:"column:day_of_week"` // 数据库中是 uint8
}

// queryStrAllByArea 是一个新的查询字符串，它只按 'area' 过滤
// 它基于 const.go 中的 queryStr，但移除了 ti.day_of_week = ?
const queryStrAllByArea = `
        SELECT 
            MAX(ci.id) AS id,
            ci.course_num,
            ti.classroom,
            ti.building,
            MAX(ci.course_type) AS course_type,
            MAX(ci.faculty) AS faculty,
            MAX(ci.course_name) AS course_name,
            MAX(ci.teacher) AS teacher,
            MAX(ci.teacher_title) AS teacher_title,
            MAX(ti.week_and_time) AS week_and_time,
            MAX(ti.day_of_week) AS day_of_week
        FROM time_infos ti 
        JOIN course_infos ci ON ci.id = ti.course_info_id
        WHERE ti.area = ? 
        GROUP BY 
            ti.building, 
            ti.classroom,
            ci.course_num
    `

// GetAllCourses 获取数据库中所有的课程信息，并按学部和教学楼分组
// (返回与 course.GetTeachInfos 相同的结构)
func (s *CourseService) GetAllCourses() ([][]vo.BuildingInfoVO, error) {

	// 模仿 building.go，我们假设有4个学部(area 1-4)
	allFacultiesData := make([][]vo.BuildingInfoVO, 4)

	// 循环遍历4个学部 (Area 1 到 4)
	for areaNum := 1; areaNum <= 4; areaNum++ {

		var results []courseQueryRow

		// 1. 执行 Raw SQL 查询，获取该学部的所有课程
		if err := database.Client.Raw(queryStrAllByArea, areaNum).Scan(&results).Error; err != nil {
			log.Printf("Service: GetAllCourses (Area %d) 查询失败: %v", areaNum, err)
			return nil, fmt.Errorf("获取所有课程 (Area %d) 的数据库操作失败: %w", areaNum, err)
		}

		// 2. 将结果按教学楼分组 (模仿 building.go)
		buildingMap := make(map[string][]vo.CourseInfoVO)

		for _, row := range results {
			// 关键区别：我们不再调用 generator.IsWeekLessonMatch
			// 我们接受所有查询到的课程

			// 转换课程时间
			// TODO: 你应该使用你的 generator 包中的函数来将 row.WeekAndTime 转换为可读字符串
			// 暂时我们使用一个占位符
			dayStr := strconv.Itoa(int(row.DayOfWeek))
			courseTimeStr := fmt.Sprintf("周%s (Raw: %d)", dayStr, row.WeekAndTime)

			// 确保 vo.CourseInfoVO 的字段被正确填充
			res := vo.CourseInfoVO{
				ID:           row.ID,
				Room:         row.Classroom,
				Faculty:      row.Faculty,
				CourseName:   row.CourseName,
				TeacherName:  row.Teacher,
				TeacherTitle: row.TeacherTitle,
				CourseTime:   courseTimeStr,
				CourseType:   row.CourseType,
			}

			buildingMap[row.Building] = append(buildingMap[row.Building], res)
		}

		// 3. 将 map 转换为 []vo.BuildingInfoVO
		buildingInfos := make([]vo.BuildingInfoVO, 0, len(buildingMap))
		for key, infos := range buildingMap {
			// 填充 vo.BuildingInfoVO
			buildingInfos = append(buildingInfos, vo.BuildingInfoVO{
				Building: key,
				Label:    key, // Label 和 Value 也填充一下
				Value:    0,   //
				Infos:    infos,
			})
		}

		// 4. 按课程数量对教学楼进行排序 (模仿 building.go)
		slices.SortFunc(buildingInfos, func(a, b vo.BuildingInfoVO) int {
			return len(b.Infos) - len(a.Infos)
		})

		// 5. 将这个学部的教学楼列表存入最终结果
		allFacultiesData[areaNum-1] = buildingInfos
	}

	return allFacultiesData, nil
}

// GetCourseDetailByID 根据课程 ID 获取课程详细信息
func (s *CourseService) GetCourseDetailByID(courseID uint) (*vo.CourseDetailVO, error) {
	//var courseModel dto.CourseInfo // 使用 dto.CourseInfo 作为课程模型
	//
	//// 使用 database.Client
	//if err := database.Client.First(&courseModel, courseID).Error; err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return nil, errors.New(config.MsgCourseNotFound) // 使用统一定义的错误消息
	//	}
	//	log.Printf("Service: 获取课程详情 (ID %d) 失败: %v", courseID, err)
	//	return nil, fmt.Errorf("获取课程详情数据库操作失败: %w", err)
	//}
	//
	//// 将模型转换为 VO
	//courseDetailVO := &vo.CourseDetailVO{
	//	ID:            courseModel.ID,
	//	Room:          courseModel.Room,
	//	Faculty:       courseModel.Faculty,
	//	CourseName:    courseModel.CourseName,
	//	TeacherName:   courseModel.TeacherName,
	//	TeacherTitle:  courseModel.TeacherTitle,
	//	CourseTime:    courseModel.CourseTime,
	//	CourseType:    courseModel.CourseType,
	//	Description:   courseModel.Description,
	//	Credits:       courseModel.Credits,
	//	AverageRating: courseModel.AverageRating,
	//	ReviewCount:   courseModel.ReviewCount,
	//}

	return nil, nil
}

// GetCourseReviewsByCourseID 根据课程 ID 获取课程评价列表
func (s *CourseService) GetCourseReviewsByCourseID(courseID uint) ([]vo.CourseReviewInfoVO, error) {
	// 1. 检查课程是否存在
	var courseExists dto.CourseInfo // 使用 dto.CourseInfo 作为课程模型
	// 使用 database.Client
	if err := database.Client.Select("id").First(&courseExists, courseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(config.MsgCourseNotFound)
		}
		log.Printf("Service: 检查课程 (ID %d) 是否存在失败: %v", courseID, err)
		return nil, fmt.Errorf("查询课程是否存在时数据库操作失败: %w", err)
	}

	var reviews []dto.CourseReviewModel // 使用 dto.CourseReviewModel
	// 预加载 User 以获取 ReviewerName。
	// dto.CourseReviewModel 有 User   dto.User `gorm:"foreignKey:UserID"`
	// 假设 dto.User 有 Username 字段
	// 使用 database.Client
	if err := database.Client.Preload("User").Where("course_id = ?", courseID).Order("created_at desc").Find(&reviews).Error; err != nil {
		log.Printf("Service: 获取课程 (ID %d) 的评价列表失败: %v", courseID, err)
		return nil, fmt.Errorf("获取课程评价列表数据库操作失败: %w", err)
	}

	reviewVOs := make([]vo.CourseReviewInfoVO, 0, len(reviews))
	for _, review := range reviews {
		reviewerName := "匿名用户"
		// 确保 review.User 结构体和其 Username 字段被正确加载和访问
		if review.User.Id != 0 && review.User.Username != "" { // dto.User 假设有 ID 和 Username
			reviewerName = review.User.Username
		}
		reviewVOs = append(reviewVOs, vo.CourseReviewInfoVO{
			ID:           review.ID, // 评价本身的 ID
			CourseID:     review.CourseID,
			Rating:       review.Rating,
			Comment:      review.Comment,
			ReviewerName: reviewerName,
			// UserID:    review.UserID, // 前端可能不直接需要，VO 中可省略或用 `json:"-"` 标记
			CreatedAt: review.CreatedAt,
		})
	}

	return reviewVOs, nil
}

// SubmitCourseReview 提交课程评价
// userID 从认证中间件中获取
// 返回创建的评价VO和错误
func (s *CourseService) SubmitCourseReview(userID uint32, payload dto.CourseReviewCreateDTO) (*vo.CourseReviewInfoVO, error) {
	// 1. 检查课程是否存在
	var course dto.CourseInfo // 使用 dto.CourseInfo
	// 使用 database.Client
	if err := database.Client.First(&course, payload.CourseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(config.MsgCourseForReviewNotFound)
		}
		log.Printf("Service: 为评价查找课程 (ID %d) 失败: %v", payload.CourseID, err)
		return nil, fmt.Errorf("为评价查询课程数据库操作失败: %w", err)
	}

	// 2. 检查用户是否存在
	var user dto.User // 假设您的 dto 包中有 User 模型，且与 CourseReviewModel 中的 User 一致
	// 使用 database.Client
	if err := database.Client.Select("id", "username").First(&user, userID).Error; err != nil { // 只选择需要的字段
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(config.MsgUserForReviewNotFound)
		}
		log.Printf("Service: 为评价查找用户 (ID %d) 失败: %v", userID, err)
		return nil, fmt.Errorf("为评价查询用户数据库操作失败: %w", err)
	}

	var createdReviewModel dto.CourseReviewModel

	// 3. 使用事务创建评价记录并更新课程统计信息
	// 使用 database.Client
	err := database.Client.Transaction(func(tx *gorm.DB) error {
		// 将前端可能传来的小数评分（如 4.5）四舍五入为最近的整数以兼容后端存储（int）
		roundedRating := int(math.Round(float64(payload.Rating)))
		reviewToCreate := dto.CourseReviewModel{ // 使用 dto.CourseReviewModel
			CourseID: payload.CourseID,
			UserID:   userID,
			Rating:   roundedRating,
			Comment:  payload.Comment,
			// User 和 Course 关联字段通常由 GORM 自动处理或在查询时 Preload，创建时不需要手动赋值模型实例
		}
		if err := tx.Create(&reviewToCreate).Error; err != nil {
			log.Printf("Service: 创建课程评价记录失败: %v", err)
			return fmt.Errorf("提交评价数据库操作失败: %w", err)
		}
		createdReviewModel = reviewToCreate // 保存刚创建的记录

		// 更新课程的平均评分和评论数 (更健壮的方式)
		var newAvgRating float32
		var newReviewCount uint

		// 使用原生 SQL 或 GORM 的聚合函数来计算，避免竞态条件
		// 这里使用 GORM 的聚合查询，更安全
		type RatingStats struct {
			Average float32
			Count   uint
		}
		var stats RatingStats
		if err := tx.Model(&dto.CourseReviewModel{}).
			Where("course_id = ?", payload.CourseID).
			Select("AVG(rating) as average, COUNT(*) as count").
			Scan(&stats).Error; err != nil {
			log.Printf("Service: 计算课程 (ID %d) 新的平均分和评价数失败: %v", payload.CourseID, err)
			return fmt.Errorf("更新课程评分信息时计算失败: %w", err)
		}
		newAvgRating = stats.Average
		newReviewCount = stats.Count

		if err := tx.Model(&dto.CourseInfo{}).Where("id = ?", payload.CourseID).Updates(map[string]interface{}{
			"average_rating": newAvgRating,
			"review_count":   newReviewCount,
		}).Error; err != nil {
			log.Printf("Service: 更新课程 (ID %d) 的评分和评价数失败: %v", payload.CourseID, err)
			return fmt.Errorf("保存课程评分信息失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err // 错误已在事务中被格式化和记录
	}

	// 构建成功响应的 VO
	// createdReviewModel 此时有 ID, CourseID, UserID, Rating, Comment, CreatedAt
	// User 关联需要从之前查询到的 user 变量获取 username
	reviewVO := &vo.CourseReviewInfoVO{
		ID:           createdReviewModel.ID,
		CourseID:     createdReviewModel.CourseID,
		Rating:       createdReviewModel.Rating,
		Comment:      createdReviewModel.Comment,
		ReviewerName: user.Username, // 从之前查询到的 user 中获取
		CreatedAt:    createdReviewModel.CreatedAt,
	}

	return reviewVO, nil
}
