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

func IsValidPort(port int) bool {
	if port <= 0 || port >= 65536 {
		return false
	}
	return true
}

func IsValidIp(ip string) bool {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false
	} else {
		return true
	}
}

// IsValidNetSeg 判断是否为合理的网段
func IsValidNetSeg(netSeg string) bool {
	_, _, err := net.ParseCIDR(netSeg)
	if err != nil {
		return false
	}
	return true
}
