package contants

import "os"

var (
	DefaultVMessID   = ""
	DefaultAlterID   = "0"
	DefaultVMessPort = ""

	DefaultSocksPort = ""
	DefaultSocksUser = ""
	DefaultSocksPass = ""

	DefaultShadowsocksPort     = ""
	DefaultShadowsocksPassword = ""
	DefaultShadowsocksMethod   = ""

	Region   = ""
	Key      = ""
	Secret   = ""
	Endpoint = "ecs.cn-hangzhou.aliyuncs.com"
)

func init() {

	Region = os.Getenv("ALI_ECS_REGION")
	Key = os.Getenv("ALI_ECS_KEY")
	Secret = os.Getenv("ALI_ECS_SECRET")

	DefaultVMessID = os.Getenv("VMESS_CLIENT_ID")
	DefaultVMessPort = os.Getenv("VMESS_PORT")

	DefaultShadowsocksPassword = os.Getenv("SHADOWSOCKS_PASSWORD")
	DefaultShadowsocksPort = os.Getenv("SHADOWSOCKS_PORT")
	DefaultShadowsocksMethod = os.Getenv("SHADOWSOCKS_METHOD")

}
