package conf

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/arstd/log"
)

type Config struct {
	RdConf    redisConfig `toml:"redis"`
	MysqlConf mysqlConfig `toml:"mysql"`
	PgConf    pgConfig    `toml:"postgres"`
	Secret    secret      `toml:"secret"`
	DstPath   dstPath     `toml:"dst_path"`
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

type secret struct {
	TokenSecret string `toml:"token_secret"`
	PwdSecret   string `toml:"pwd_secret"`
}

type dstPath struct {
	Path string `toml:"path"`
}

func Load() *Config {
	var conf = &Config{}
	pwd, err := os.Getwd()
	if err != nil {
		log.Error("Failed to get pwd")
		return nil
	}
	confPath := pwd + "/conf/conf.toml"
	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		log.Error(err)
		return nil
	}
	return conf
}
