package database

import (
	"cengkeHelperBackGo/internal/config"
	"cengkeHelperBackGo/internal/models/dto"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var Client *gorm.DB

func GetDB() *gorm.DB {
	return Client
}
func init() {
	var err error
	var cfg gorm.Config
	cfg = gorm.Config{
		PrepareStmt: true,
		Logger:      gormLogger.Default.LogMode(gormLogger.Info),
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix: "test",
		//},
		ConnPool: nil,
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Conf.Mysql.User, config.Conf.Mysql.Password,
		config.Conf.Mysql.Host, config.Conf.Mysql.Port,
		config.Conf.Mysql.Database)
	// 连接到SQLite数据库
	if Client, err = gorm.Open(mysql.Open(dsn), &cfg); err != nil {
		panic(err)
	}

	TableAutoMigrate()
}

func TableAutoMigrate() {
	if !config.Conf.Mysql.AutoMigrate {
		fmt.Println("未启用迁移数据库")
		return
	}

	// 要迁移的所有模型
	modelsToMigrate := []interface{}{
		&dto.User{},
		&dto.Division{},
		&dto.BuildingInfo{},
		&dto.CourseInfo{},
		&dto.Post{},
		&dto.Comment{},
		&dto.UserPostCollect{},
		&dto.UserPostLike{},
		&dto.UserCommentLike{},
	}

	// 批量执行自动迁移
	if err := Client.AutoMigrate(modelsToMigrate...); err != nil {
		panic(fmt.Errorf("数据库迁移失败: %v", err))
	}

}
