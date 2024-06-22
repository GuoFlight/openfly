package common

import (
	"fmt"
	"openfly/conf"
	"openfly/logger"
	"testing"
)

func TestBackupFile(t *testing.T) {
	conf.ParseConfig("../config.toml")
	conf.PathData = "../Test/l4"
	conf.GConf.Openfly.PathBak = "../Test/bak"
	logger.InitLog()
	output, gerr := BackupFile(conf.PathData)
	if gerr != nil {
		t.Error(gerr)
		return
	}
	fmt.Println(output)
}
