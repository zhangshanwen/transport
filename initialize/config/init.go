package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	v *viper.Viper
)

func InitConf(project string, c interface{}) {
	v = viper.New()
	v.SetConfigFile(fmt.Sprintf("./env/%s.yaml", project))
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		logrus.Fatal("Fatal error config file: %s \n", err)
	}
	if err := v.Unmarshal(&c); err != nil {
		logrus.Fatal(err)
	}
	return
}

func WriteInConfig(key string, value interface{}) error {
	v.Set(key, value)
	return v.WriteConfig()
}
