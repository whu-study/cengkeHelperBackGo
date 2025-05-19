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
func convertCourseModelToInfoVO(course dto.CourseInfo) vo.CourseInfoVO {
	return vo.CourseInfoVO{
		ID:           course.ID, // 确保 vo.CourseInfoVO 有 ID 字段
		Room:         course.Room,
		Faculty:      course.Faculty,
		CourseName:   course.CourseName,
		TeacherName:  course.TeacherName,
		TeacherTitle: course.TeacherTitle,
		CourseTime:   course.CourseTime,
		CourseType:   course.CourseType,
	}
}

// convertBuildingModelToVO 将 dto.BuildingInfo (模型) 转换为 vo.BuildingInfoVO
func convertBuildingModelToVO(building dto.BuildingInfo) vo.BuildingInfoVO {
	var courseInfosVO []vo.CourseInfoVO
	for _, courseModel := range building.Courses {
		courseInfosVO = append(courseInfosVO, convertCourseModelToInfoVO(courseModel))
	}
	return vo.BuildingInfoVO{
		Building: building.Name, // dto.BuildingInfo 中的 Name 对应前端的 building
		Label:    building.Label,
		Value:    building.Value,
		Infos:    courseInfosVO,
	}
}

// GetCourseList 获取按学部和教学楼分组的所有课程
func (s *CourseService) GetCourseList() ([][]vo.BuildingInfoVO, error) {
	var divisions []dto.Division // 使用您定义的 dto.Division

	// 使用 database.Client 进行数据库操作
	// 确保 Preload 路径与您的 GORM 模型定义匹配
	// dto.Division 包含 Buildings []BuildingInfo `gorm:"foreignKey:DivisionID"`
	// dto.BuildingInfo 包含 Courses []CourseInfo `gorm:"foreignKey:BuildingID"`
	// dto.CourseInfo 没有直接关联 Teacher，如果需要 TeacherName/Title，它们已是 CourseInfo 的字段
	if err := database.Client.
		Preload("Buildings.Courses"). // 根据您的 dto.Division 和 dto.BuildingInfo 结构预加载
		Find(&divisions).Error; err != nil {
		log.Printf("Service: 获取学部及课程列表失败: %v", err)
		return nil, fmt.Errorf("数据库查询失败: %w", err)
	}

	// 转换为前端需要的格式
	var result [][]vo.BuildingInfoVO
	for _, div := range divisions {
		var buildingsVO []vo.BuildingInfoVO
		for _, buildingModel := range div.Buildings {
			buildingsVO = append(buildingsVO, convertBuildingModelToVO(buildingModel))
		}
		result = append(result, buildingsVO)
	}

	return result, nil
}

// GetCourseDetailByID 根据课程 ID 获取课程详细信息
func (s *CourseService) GetCourseDetailByID(courseID uint) (*vo.CourseDetailVO, error) {
	var courseModel dto.CourseInfo // 使用 dto.CourseInfo 作为课程模型

	// 使用 database.Client
	if err := database.Client.First(&courseModel, courseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(config.MsgCourseNotFound) // 使用统一定义的错误消息
		}
		log.Printf("Service: 获取课程详情 (ID %d) 失败: %v", courseID, err)
		return nil, fmt.Errorf("获取课程详情数据库操作失败: %w", err)
	}

	// 将模型转换为 VO
	courseDetailVO := &vo.CourseDetailVO{
		ID:            courseModel.ID,
		Room:          courseModel.Room,
		Faculty:       courseModel.Faculty,
		CourseName:    courseModel.CourseName,
		TeacherName:   courseModel.TeacherName,
		TeacherTitle:  courseModel.TeacherTitle,
		CourseTime:    courseModel.CourseTime,
		CourseType:    courseModel.CourseType,
		Description:   courseModel.Description,
		Credits:       courseModel.Credits,
		AverageRating: courseModel.AverageRating,
		ReviewCount:   courseModel.ReviewCount,
	}

	return courseDetailVO, nil
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
		reviewToCreate := dto.CourseReviewModel{ // 使用 dto.CourseReviewModel
			CourseID: payload.CourseID,
			UserID:   userID,
			Rating:   payload.Rating,
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
