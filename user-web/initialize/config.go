package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"mxshop-api/user-web/global"
)

func GetEnvInfo(env string)int{
	viper.AutomaticEnv()
	return viper.GetInt(env)
}



func InitConfig(){
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFileName := "user-web/config-pro.yaml"
	if debug == 1{
		configFileName = "user-web/config-debug.yaml"
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig();err != nil{
		panic(err)
	}

	if err := v.Unmarshal(global.ServerConfig);err != nil{
		panic(err)
	}
	fmt.Println(global.ServerConfig)


	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event){
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		fmt.Println(global.ServerConfig)
	})
}
