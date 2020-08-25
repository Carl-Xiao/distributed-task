package common

import (
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//Appconfig INI配置文件
type Appconfig struct {
	*Server     `ini:"server"`
	*LogConfig  `ini:"logconfig"`
	*EtcdConfig `ini:"etcd"`
}

//Server INI配置文件
type Server struct {
	PORT         int `ini:"PORT"`
	ReadTimeout  int `ini:"READ_TIMEOUT"`
	WriteTimeout int `ini:"READ_TIMEOUT"`
}

//LogConfig 日志配置文件
type EtcdConfig struct {
	EndPoint []string `ini:"endPoint"`
	TimeOut  string   `ini:"timeOut"`
}

//LogConfig 日志配置文件
type LogConfig struct {
	Level      string `ini:"level"`
	Filename   string `ini:"filename"`
	MaxSize    int    `ini:"maxsize"`
	MaxAge     int    `ini:"max_age"`
	MaxBackups int    `ini:"max_backups"`
}

var (
	App          *Appconfig
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PORT         int
	EndPoint     []string
)

//InitBase 初始化配置文件
func InitBase() (err error) {
	App = new(Appconfig)
	err = ini.MapTo(App, "app.ini")
	if err != nil {
		return
	}
	ReadTimeout = time.Duration(App.ReadTimeout) * time.Second
	WriteTimeout = time.Duration(App.WriteTimeout) * time.Second
	PORT = App.PORT
	EndPoint = App.EndPoint
	//加载logger
	err = InitLogger(App.LogConfig)
	if err != nil {
		return
	}
	return
}
