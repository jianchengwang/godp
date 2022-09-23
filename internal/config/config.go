package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type appConfig struct {
	CrossOrigin string `yaml:"cross_origin"mapstructure:"cross_origin"` // underscore need add mapstructure tag or cant unmarshal
	WorkPath    string `yaml:"work_path"mapstructure:"work_path"`
}

type mysqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type config struct {
	App   appConfig   `yaml:"app"`
	Mysql mysqlConfig `yaml:"mysql"`
}

var (
	Config config
	Env    = ""
)

func init() {
	env := os.Getenv("APP_ENV")
	if len(env) > 0 {
		Env = env
		viper.SetConfigName("config_" + Env)
	} else {
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/godp/") // path to look for the config file in
	viper.AddConfigPath(".")          // lock for current path
	err := viper.ReadInConfig()       // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		panic(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
	fmt.Println(Config)
}
