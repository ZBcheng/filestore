package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RdConf    redisConfig `toml:"redis"`
	MysqlConf mysqlConfig `toml:"mysql"`
	PgConf    pgConfig    `toml:"postgres"`
}

type redisConfig struct {
	Host      string `toml:"host"`
	Port      string `toml:"port"`
	MaxIdle   int    `toml:"maxIdle"`
	MaxActive int    `toml:"maxActive"`
}

type mysqlConfig struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"username"`
	DBName   string `toml:"dbname"`
	Password string `toml:"password"`
}

type pgConfig struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"username"`
	DBName   string `toml:"dbname"`
	Password string `toml:"password"`
}

var conf *Config

func init() {
	confPath := "/Users/zhangbicheng/PycharmProjects/filestore/conf/conf.toml"
	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		fmt.Println(err)
	}
}

func GetConfig() *Config {
	return conf
}
