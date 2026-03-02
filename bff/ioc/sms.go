package ioc

import (
	"github.com/chiren-c/chili/sms/repository"
	"github.com/chiren-c/chili/sms/repository/dao"
	"github.com/chiren-c/chili/sms/service"
	"github.com/chiren-c/chili/sms/service/async"
	"github.com/chiren-c/chili/sms/service/tencent"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func InitSmsTencentService() service.Service {
	type Config struct {
		SecretID  string `yaml:"secretId"`
		SecretKey string `yaml:"secretKey"`
	}
	var cfg Config
	err := viper.UnmarshalKey("tencentSms", &cfg)
	c, err := tencentSMS.NewClient(common.NewCredential(cfg.SecretID, cfg.SecretKey),
		"ap-nanjing",
		profile.NewClientProfile())
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "", "")
}

var SmsProvider = wire.NewSet(
	InitSmsTencentService,
	async.NewService,
	repository.NewAsyncSMSRepository,
	dao.NewGORMAsyncSmsDAO,
)
