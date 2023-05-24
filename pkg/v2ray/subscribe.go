package v2ray

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type VMessVesion string

const (
	VMessVersion2 VMessVesion = "2"
)
const (
	SocksProxy       string = "socks"
	ShadowsocksProxy string = "shadowsocks"
	VMessProxy       string = "vmess"
)

type vmessConfig struct {
	Version VMessVesion `json:"v"`
	Name    string      `json:"ps"`
	Address string      `json:"add"`
	Port    string      `json:"port"`
	ID      string      `json:"id"`
	AlterId string      `json:"aid"`
	SCY     string      `json:"scy"`
	Network string      `json:"net"`
	Type    string      `json:"type"`
	Host    string      `json:"host"`
	Path    string      `json:"path"`
	TLS     string      `json:"tls"`
	SNI     string      `json:"sni"`
	Alpn    string      `json:"alpn"`
}

type SubscribeFunc func(string, string, string, map[string]string) string

func vmess(name, address, port string, variables map[string]string) string {

	// TODO: 参数太多建议优化一下
	defaultVMESS := vmessConfig{
		Version: VMessVersion2,
		Name:    name,
		Address: address,
		Port:    port,
		ID:      variables["clientid"],
		AlterId: variables["alterid"],
		SCY:     variables["scy"],
		Network: variables["network"],
		Type:    variables["type"],
		Host:    variables["host"],
		Path:    variables["path"],
		TLS:     variables["tls"],
		SNI:     variables["sni"],
		Alpn:    variables["alpn"],
	}
	t, _ := json.Marshal(defaultVMESS)
	bs64 := base64.StdEncoding.EncodeToString(t)
	return fmt.Sprintf("vmess://%s", bs64)
}

func socks(name, address, port string, variables map[string]string) string {
	secret := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", variables["user"], variables["password"])))
	return fmt.Sprintf("socks://%s@%s:%s#%s", secret, address, port, name)
}

func shadowSocks(name, address, port string, variables map[string]string) string {

	methodAndPassword := fmt.Sprintf("%s:%s", variables["method"], variables["password"])
	bs64 := base64.StdEncoding.EncodeToString([]byte(methodAndPassword))
	return fmt.Sprintf("ss://%s@%s:%s#%s", bs64, address, port, name)
}

func Subscribe(subConfigs []*SubscribeConfig, dynamicAddress string) (string, error) {

	all := make([]string, 0)

	funcMap := map[string]SubscribeFunc{
		SocksProxy:       socks,
		VMessProxy:       vmess,
		ShadowsocksProxy: shadowSocks,
	}

	for _, s := range subConfigs {

		f, isOK := funcMap[s.Proxy]
		if !isOK {
			return "", fmt.Errorf("error subscribe config: %v", s)
		}
		if s.Dynamic {
			all = append(all, f(s.Name, dynamicAddress, s.Port, s.Variables))
		} else {
			all = append(all, f(s.Name, s.Address, s.Port, s.Variables))
		}
	}

	allStr := strings.Join(all, "\n")
	return base64.StdEncoding.EncodeToString([]byte(allStr)), nil
}
