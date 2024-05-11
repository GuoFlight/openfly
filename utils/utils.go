package utils

import (
	"errors"
	"net"
	"openfly/logger"
	"os"
	"path/filepath"
	"strconv"
)

func IsPortListening(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return true
	}
	defer ln.Close()
	return false
}
func DelFile(path string) error {
	logger.GLogger.Info("即将删除文件：", path)
	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if pathAbs == "/" {
		return errors.New("无法删除/")
	}
	return os.RemoveAll(pathAbs)
}
