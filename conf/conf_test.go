package conf

import "testing"

func TestGetConf(t *testing.T) {
	t.Log(Load().DstPath.Path)
}
