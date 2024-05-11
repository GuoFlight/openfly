#!/bin/bash
app="openfly"
dir_output="output"

version=$(go run ../main.go -v)   # 获取版本号

file_openfly="${app}-${version}"
file_config="config-${version}.toml"