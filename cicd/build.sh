#!/bin/bash
# 初始化
cd $(dirname "$0")
source ./env.sh
ouput_dir="build/${version}"
mkdir -p ${ouput_dir}
cp ../config.toml ${ouput_dir}/
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${ouput_dir}/${app}_mac_amd64 ../main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ${ouput_dir}/${app}_mac_arm64 ../main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${ouput_dir}/${app}_linux_amd64 ../main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${ouput_dir}/${app}_linux_arm64 ../main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${ouput_dir}/${app}_win_amd64.exe ../main.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ${ouput_dir}/${app}_win_arm64.exe ../main.go
chmod +x ${ouput_dir}/${app}*

