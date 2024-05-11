#!/bin/bash
# 初始化
cd $(dirname "$0")
source ./env.sh
mkdir -p ${dir_output}

# 编译
if [ "$version" != "vx.x.x" ] && ([ -e "${dir_output}/$file_openfly" ] || [ -e "${dir_output}/$file_config" ]); then
  echo "${dir_output}中目标文件已存在"
  exit 1
fi
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o "${dir_output}/${file_openfly}" ../main.go  || exit 1
cp ../config.toml "${dir_output}/$file_config"