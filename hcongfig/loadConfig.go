package hconfig

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func ReadConfigFromFile(fileName string, v any, needListen bool) error {
	vp := viper.New()
	if err := loadConfig(vp, fileName, v); err != nil {
		return err
	}
	if needListen {
		vp.WatchConfig()
		vp.OnConfigChange(func(in fsnotify.Event) {
			fmt.Println("读取配置文件")
			if err := loadConfig(vp, fileName, v); err != nil {
				fmt.Printf("读取配置文件失败：%s\n", err)
			} else {
				fmt.Printf("读取配置文件成功!\n")
			}
		})
	}
	return nil
}

func loadConfig(vp *viper.Viper, fileName string, v any) error {
	//设置配置文件路径
	vp.SetConfigFile(fileName)
	//读取配置文件消息
	if err := vp.ReadInConfig(); err != nil {
		return err
	}
	return vp.Unmarshal(v)

}
