package main

import (
	"github.com/chiren-c/chili/bff/ioc"
	"github.com/chiren-c/chili/user/repository/dao"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	initViperWatch()
	_ = InitTables(ioc.InitDB())
	app := InitApp()
	err := app.WebServer.Start()
	panic(err)
}

func initViperWatch() {
	cfile := pflag.String("config",
		"bff/config/config.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&dao.User{},
	)
}
