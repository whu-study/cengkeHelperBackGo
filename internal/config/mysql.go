package config

import (
	"cengkeHelperBackGo/internal/models/dto"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var Client *gorm.DB

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
		Conf.Mysql.User, Conf.Mysql.Password,
		Conf.Mysql.Host, Conf.Mysql.Port,
		Conf.Mysql.Database)
	// 连接到SQLite数据库
	if Client, err = gorm.Open(mysql.Open(dsn), &cfg); err != nil {
		panic(err)
	}

	TableAutoMigrate()
}

func TableAutoMigrate() {
	//if !config.EnvCfg.AutoMigrate {
	//	logger.Info("未启用迁移数据库")
	//	return
	//}
	if err := Client.AutoMigrate(&dto.User{}); err != nil {
		panic(err)
		return
	}

}
