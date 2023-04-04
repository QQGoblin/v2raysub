package contants

import "os"

var (
	DefaultVMessID   = ""
	DefaultAlterID   = "0"
	DefaultVMessPort = ""

	DefaultSockPort = ""
	DefaultSockUser = ""
	DefaultSockPass = ""

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

}
