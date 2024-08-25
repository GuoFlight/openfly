package conf

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
)

var (
	GConf            ConfigFileStruct
	PathData         = "./data"
	DefaultPathBak   = "./bak"
	DefaultPathL4Log = "/var/log/nginx/l4"
)

// ConfigFileStruct 配置文件结构体
type ConfigFileStruct struct {
	Openfly struct {
		PathBak string `toml:"pathBak"`
	} `toml:"openfly"`
	API struct {
		AdminUser     string `toml:"admin_user"`
		AdminPassword string `toml:"admin_password"`
		Expire        int    `toml:"expire"`
	} `toml:"api"`
	Log struct {
		Level         string `toml:"level"`
		Path          string `toml:"path"`
		RotationCount int    `toml:"rotationCount"`
	} `toml:"log"`
	Http struct {
		Port int `toml:"port"`
	} `toml:"http"`
	Etcd struct {
		Server  string `toml:"server"`
		Timeout int    `toml:"timeout"`
		Prefix  string `toml:"prefix"`
	} `toml:"etcd"`
	Nginx struct {
		WeightDefault      int `toml:"weightDefault"`
		MaxFailsDefault    int `toml:"maxFailsDefault"`
		FailTimeoutDefault int `toml:"failTimeoutDefault"`
	} `toml:"nginx"`
	L4 struct {
		Log struct {
			Path              string `toml:"path"`
			FormatNameDefault string `toml:"formatNameDefault"`
		} `toml:"log"`
	} `toml:"l4"`
}

// ParseConfig 解析配置文件
func ParseConfig(pathConfFile string) {
	if _, err := toml.DecodeFile(pathConfFile, &GConf); err != nil {
		log.Fatal(err)
	}
	CheckAndInit()
}
func CheckAndInit() {
	// 设置工作目录
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	wd := filepath.Dir(path)
	log.Println("[INFO] work directory:", wd)
	err = os.Chdir(wd)
	if err != nil {
		log.Fatal(err)
	}
	// 创建数据目录
	err = os.MkdirAll(PathData, 0770)
	if err != nil {
		log.Fatal(err)
	}
	// 创建备份目录
	if GConf.Openfly.PathBak == "" {
		GConf.Openfly.PathBak = DefaultPathBak
	}
	err = os.MkdirAll(GConf.Openfly.PathBak, 0770)
	if err != nil {
		log.Fatal(err)
	}
	// 创建L4日志目录
	if GConf.L4.Log.Path == "" {
		GConf.L4.Log.Path = DefaultPathL4Log
	}
	err = os.MkdirAll(GConf.L4.Log.Path, 0770)
	if err != nil {
		log.Fatal(err)
	}
}
