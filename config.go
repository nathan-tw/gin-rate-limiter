package main

import "github.com/spf13/viper"

type Redis struct {
	Host string `yaml:"Host"`
	Port string    `yaml:"Port"`
}

type App struct {
	Timeout int `yaml:"Timeout"`
	MaxIP   int `yaml:"MaxIP"`
}

var (
	redisSetting Redis
	appSetting App
)


func initConfig() (*viper.Viper, error) {
	v := viper.New()	
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./")
	
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	} 

	v.WatchConfig()
	return v, nil
}
