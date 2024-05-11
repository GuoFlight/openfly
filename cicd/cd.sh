#!/bin/bash
server="172.16.1.128"
passwordSsh="1"
pathDeploy="/opt/openfly"

# 初始化
cd $(dirname "$0")
source ./env.sh

# 创建openfly的目录
sshpass -p ${passwordSsh} ssh root@${server} "mkdir -p ${pathDeploy}" || exit 1

# 停止老openfly
sshpass -p ${passwordSsh} ssh root@${server} "ps aux | grep openfly | grep -v grep | awk '{print \$2}' | xargs -I {} kill -9 {}"

# 部署
sshpass -p ${passwordSsh} scp ${dir_output}/${file_openfly} ${dir_output}/${file_config}  root@${server}:${pathDeploy} || exit 1

# 启动
sshpass -p ${passwordSsh} ssh root@${server} "nohup ${pathDeploy}/${file_openfly} -c ${pathDeploy}/${file_config} > nohup.out &"

exit

