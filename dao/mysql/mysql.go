package mysql

import (
	"fmt"
	"github.com/yongtenglei/newThing/model"
	"github.com/yongtenglei/newThing/settings"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		settings.UserServiceConf.MySQLConf.User,
		settings.UserServiceConf.MySQLConf.Password,
		settings.UserServiceConf.MySQLConf.Host,
		settings.UserServiceConf.MySQLConf.Port,
		settings.UserServiceConf.MySQLConf.DbName)
	println(dsn)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.S().Errorw("Open MySQL failed", "err", err.Error())
		panic(err)
	}

	err = DB.AutoMigrate(&model.User{}, &model.TokenSession{})
	if err != nil {
		zap.S().Errorw("AutoMigrate model failed", "err", err.Error())
		panic(err)
	}
}
