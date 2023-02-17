package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Port     int    `mapstructure:"port"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init() (err error) {
	//下面两行并不是简单的拼接，上面才是指定文件名字，下面那个用于viper远程访问时才会使用
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	viper.SetConfigFile("config.yaml") // 指定配置文件
	viper.AddConfigPath(".")           // 指定查找配置文件的路径
	err = viper.ReadInConfig()         // 读取配置信息
	if err != nil {                    // 读取配置信息失败
		fmt.Printf("viper 读取配置错误 err : %#v\n", err)
		return
	}
	//把读取到的信息反序列化到结构体Conf中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper unmarshall failed err : %#v\n", err)
	}
	// 监控配置文件变化
	viper.WatchConfig()
	// 配置发生变化的回调函数
	var cnt int = 0
	viper.OnConfigChange(func(in fsnotify.Event) {
		//把读取到的信息反序列化到结构体Conf中
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper unmarshall failed err : %#v\n", err)
		}
		cnt++
		fmt.Println(cnt)
		fmt.Println("Config file changed", in.Name)
	})
	return
}
