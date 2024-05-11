package common

import (
	"fmt"
	"openfly/conf"
	"openfly/logger"
	"testing"
)

func TestBackupFile(t *testing.T) {
	conf.ParseConfig("../config.toml")
	conf.PathData = "../test/l4"
	conf.GConf.Openfly.PathBak = "../test/bak"
	logger.InitLog()
	output, gerr := BackupFile(conf.PathData)
	if gerr != nil {
		t.Error(gerr)
		return
	}
	fmt.Println(output)
}
