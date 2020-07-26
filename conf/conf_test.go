package conf

import (
	"os"
	"testing"
)

func TestGetConf(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Error("Failed to get pwd")
		t.Error(err)
		return
	}
	confPath := pwd + "/conf.toml"
	t.Log(confPath)
}

func TestLoadConf(t *testing.T) {
	conf := Load().MysqlConf.DBName
	t.Log(conf)
}
