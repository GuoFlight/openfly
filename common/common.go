package common

import (
	"github.com/GuoFlight/gerror"
	"github.com/sirupsen/logrus"
	"math/rand"
	"openfly/conf"
	"openfly/logger"
	"os/exec"
	"path"
	"strconv"
	"time"
)

func BackupFile(srcFile string) (string, *gerror.Gerr) {
	srcBaseName := path.Base(srcFile)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	destBaseName := srcBaseName + "_" + time.Now().Format("20060102_150405_") + strconv.Itoa(rand.Intn(1000))
	destFile := path.Join(conf.GConf.Openfly.PathBak, destBaseName)
	logger.GLogger.WithFields(logrus.Fields{"源文件": srcFile, "目标文件": destFile}).Info("即将备份nginx配置文件")
	cmd := exec.Command("cp", "-r", srcFile, destFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), gerror.NewErr("备份文件失败: " + err.Error() + " " + string(output))
	}
	logger.GLogger.WithFields(logrus.Fields{"源文件": srcFile, "目标文件": destFile}).Info("文件备份完成")
	return string(output), nil
}
