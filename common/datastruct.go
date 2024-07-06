package common

import "openfly/conf"

type NginxConfL4 struct {
	Disable      bool            `json:"disable,omitempty"`
	Listen       int             `json:"listen"`
	Upstream     Upstream        `json:"upstream"`
	IncludeFiles []string        `json:"includeFiles,omitempty"`
	WhiteList    []WhiteListItem `json:"whiteList,omitempty"`
	Comments     []string        `json:"comments,omitempty"`
	// 限速
	ProxyUploadRate   string `json:"proxyUploadRate,omitempty"`
	ProxyDownloadRate string `json:"proxyDownloadRate,omitempty"`
	// 超时时间
	ProxyConnectTimeout string `json:"proxyConnectTimeout,omitempty"`
	ProxyTimeout        string `json:"proxyTimeout,omitempty"`
}
type WhiteListItem struct {
	Type   conf.OpWhiteList `json:"type"`
	Target string           `json:"target"`
}

type Upstream struct {
	// 后端机器
	Hosts []UpstreamHost `json:"hosts"`
	// 哈希
	IsHash    bool   `json:"isHash,omitempty"`
	HashField string `json:"hashField,omitempty"`
	// 主动健康检查
	Interval int `json:"interval,omitempty"`
	Rise     int `json:"rise,omitempty"`
	Fall     int `json:"fall,omitempty"`
	Timeout  int `json:"timeout,omitempty"`
}
type UpstreamHost struct {
	Ip                string `json:"ip"`
	Port              int    `json:"port"`
	Weight            int    `json:"weight,omitempty"`
	MaxFails          int    `json:"maxFails,omitempty"`
	FailTimeoutSecond int    `json:"failTimeoutSecond,omitempty"`
	IsBackup          bool   `json:"isBackup,omitempty"`
}
