package seetings

import (
	"fmt"
	"github.com/spf13/viper"
)

type Setting struct {
	Mysql MySQLSetting `yaml:"mysql"`
	RPC   RPCSetting   `yaml:"rpc"`
}

type MySQLSetting struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	Timeout  int    `yaml:"timeout"`
}

type RPCSetting struct {
	DbServPort string `yaml:"dbservport"`
	Timeout    int    `yaml:"timeout"`
}

var setting *Setting
var err error

func init() {
	setting = &Setting{}
	vp := viper.New()
	vp.AddConfigPath("configs")
	vp.SetConfigName("conf")
	vp.SetConfigType("yaml")
	err = vp.ReadInConfig()
	if err != nil {
		fmt.Println("read config err! ", err)
		return
	}
	err = vp.Unmarshal(&setting)
	if err != nil {
		fmt.Println("unmarshal config err! ", err)
		return
	}
}

func GetSetting() (*Setting, error) {
	return setting, err
}
