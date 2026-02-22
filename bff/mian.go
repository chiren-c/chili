package main

import (
	"github.com/chiren-c/chili/bff/ioc"
	"github.com/chiren-c/chili/user/repository/dao"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	initViperWatch()
	_ = InitTables(ioc.InitDB())
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // 默认监听 0.0.0.0:8080
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
