package initialize

import (
	"log"
	"tiktok/cmd/api/global"

	"github.com/spf13/viper"
)

func Viper() {
	// 设置配置文件类型和路径
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config/config.yaml")
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("读取配置文件错误")
	}
	// 将配置反序列化到全局配置中
	err = viper.Unmarshal(&global.Config)
	if err != nil {
		log.Panic("配置信息反序列化错误")
	}
}
