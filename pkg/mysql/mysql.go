package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormDefaultLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
	"timeLineGin/internal/model"
	"timeLineGin/pkg/config"
	"timeLineGin/pkg/logger"
)

var instance *gorm.DB

func GetInstance() *gorm.DB {
	return instance
}

func Initialize(conf *config.DB) {
	DSN := fmt.Sprintf(
		"%s:%d@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Name,
		conf.Pass,
		conf.Host,
		conf.Port,
		conf.DB)

	mysqlCfg := mysql.Config{
		DSN: DSN,
	}
	// 配置日志
	logMode := gormDefaultLogger.Error
	if conf.Debug {
		logMode = gormDefaultLogger.Info
	}
	l := gormDefaultLogger.Default.LogMode(logMode)

	// 连接数据库
	gormDB, err := gorm.Open(mysql.New(mysqlCfg), &gorm.Config{
		Logger: l,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tl_",
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 连接池配置
	db, _ := gormDB.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour * 2)
	Migrate(gormDB)

	// 初始化
	instance = gormDB
}

// Migrate 数据库迁移
func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		new(model.UserInputCreate),
	); err != nil {
		logger.GetInstance().Errorf("数据库迁移失败: %v", err)
	}
	logger.GetInstance().Info("数据库迁移成功")
}
