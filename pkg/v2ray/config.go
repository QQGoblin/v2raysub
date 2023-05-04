package v2ray

type Config struct {
	Path      string             `yaml:"path"`
	Tmpl      string             `yaml:"tmpl"`
	Subscribe []*SubscribeConfig `yaml:"subscribe"`
}

type SubscribeConfig struct {
	Name      string            `yaml:"name"`      // 代理名称
	Proxy     string            `yaml:"proxy"`     // 代理协议类型
	Address   string            `yaml:"address"`   // 代理地址
	Port      string            `yaml:"port"`      // 代理地址
	Dynamic   bool              `yaml:"dynamic"`   // 动态变更 ip 地址
	Variables map[string]string `yaml:"variables"` // 其他代理参数
}
