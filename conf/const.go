package conf

var (
	Version         = "v0.3.0"
	NgFileExtension = ".conf"
	EtcdSubPathL4   = "/l4"
)

type OpWhiteList string

const (
	Allow OpWhiteList = "allow"
	Deny  OpWhiteList = "deny"
)
