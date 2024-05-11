# 简介

作者：京城郭少

# 简介

> 基于nginx的4层代理项目

支持的功能：
* 被动健康检查
* 白名单
* include导入文件
* 哈希
* backup冗余互备
* weight权重
* ......

# 启动

```shell
go run main.go
go run main.go -c ./config.toml
```

# 部署openfly

部署nginx：

* 目标：部署一个支持stream模块的nginx。

```shell
systemctl stop firewalld
systemctl disable firewalld
setenforce 0
sed -i '/^SELINUX=/cSELINUX=disabled' /etc/selinux/config
yum install -y gc gcc gcc-c++ pcre-devel zlib-devel openssl-devel libxslt-devel GeoIP-devel perl-ExtUtils-Embed make
wget http://nginx.org/download/nginx-1.24.0.tar.gz
tar -xvf nginx-1.24.0.tar.gz
cd nginx-1.24.0/
mkdir -p /usr/local/nginx
# 关键在于--with-stream=dynamic --with-stream_ssl_module
./configure --prefix=/usr/local/nginx --with-file-aio --with-http_auth_request_module --with-http_ssl_module --with-http_v2_module --with-http_realip_module --with-http_addition_module --with-http_xslt_module=dynamic --with-http_geoip_module=dynamic --with-http_sub_module --with-http_dav_module --with-http_flv_module --with-http_mp4_module --with-http_gunzip_module --with-http_gzip_static_module --with-http_random_index_module --with-http_secure_link_module --with-http_degradation_module --with-http_slice_module --with-http_stub_status_module --with-http_perl_module=dynamic --with-pcre --with-pcre-jit --with-stream=dynamic --with-stream_ssl_module
make && make install
cp -r contrib/vim/* /usr/share/vim/vimfiles/
ln -s /usr/local/nginx/conf/ /etc/nginx
ln -s /usr/local/nginx/sbin/nginx /usr/sbin/nginx
vim /usr/lib/systemd/system/nginx.service
    [Unit]
    Description=The nginx HTTP and reverse proxy server
    After=network.target remote-fs.target nss-lookup.target

    [Service]
    Type=forking
    PIDFile=/usr/local/nginx/logs/nginx.pid
    ExecStartPre=/usr/bin/rm -f /usr/local/nginx/logs/nginx.pid
    ExecStartPre=/usr/sbin/nginx -t
    ExecStart=/usr/sbin/nginx -c /etc/nginx/nginx.conf
    ExecReload=/bin/kill -s HUP $MAINPID
    KillSignal=SIGQUIT
    TimeoutStopSec=5
    KillMode=process
    PrivateTmp=true

    [Install]
    WantedBy=multi-user.target
systemctl daemon-reload
systemctl restart nginx
systemctl enable nginx
```

配置nginx：

* 目标：配置一个4层代理的目录（此目录托管给openfly，生成相应的4层代理配置）

```shell
mkdir -p /etc/nginx/stream.d
vim /etc/nginx/nginx.conf
    load_module /usr/local/nginx/modules/ngx_stream_module.so;      # 此配置放在文件的首行
    ......
    stream {
        include /etc/nginx/stream.d/*.conf;
    }
nginx -t
nginx -s reload
```

部署etcd：

```shell
yum install -y etcd
vim /etc/etcd/etcd.conf
    ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:2379"
systemctl enable --now etcd
vim /etc/bashrc
    export ETCDCTL_API=3
source /etc/bashrc
```

# 启动openfly

```shell
vim config-vx.x.x.toml                  # 编辑配置文件
./openfly-vx.x.x -c config-vx.x.x.toml  # 启动openfly
```

<br>

（未经许可请不要应用在任何商业用途上）