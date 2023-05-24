package config

import (
	"github.com/QQGoblin/v2raysub/pkg/aliyun"
	"github.com/QQGoblin/v2raysub/pkg/v2ray"
)

type SubConfig struct {
	V2ray  *v2ray.Config  `yaml:"v2ray"`
	Aliyun *aliyun.Config `yaml:"aliyun"`
}
