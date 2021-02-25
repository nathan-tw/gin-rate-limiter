package global

import (
	"github.com/nathan-tw/gin-rate-limiter/pkg/logger"
	"github.com/nathan-tw/gin-rate-limiter/pkg/setting"
)

var (
	RedisSetting  *setting.RedisSettingS
	ServerSetting *setting.ServerSettingS
	AppSetting    *setting.AppSettingS
	Logger        *logger.Logger
)
