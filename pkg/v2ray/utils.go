package v2ray

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/QQGoblin/v2raysub/pkg/config"
)

type VMessVesion string

const (
	VMessVersion2 VMessVesion = "2"
)
const (
	SocksProxyName       string = "socks"
	ShadowsocksProxyName string = "shadowsocks"
	VMessProxyName       string = "vmess"
)

type vmess struct {
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

func defaultVMess(address string, alterID string) string {

	defaultVMESS := vmess{
		Version: VMessVersion2,
		Name:    VMessProxyName,
		Address: address,
		Port:    config.DefaultVMessPort,
		ID:      config.DefaultVMessID,
		AlterId: alterID,
		SCY:     "auto",
		Network: "ws",
		Type:    "none",
		Host:    "goodx.com",
		Path:    "/goodx",
		TLS:     "tls",
		SNI:     "",
		Alpn:    "",
	}
	t, _ := json.Marshal(defaultVMESS)
	bs64 := base64.StdEncoding.EncodeToString(t)
	return fmt.Sprintf("vmess://%s", bs64)
}

func defaultSocks(address string, alterID string) string {
	secret := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.DefaultSocksUser, config.DefaultSocksPass)))
	return fmt.Sprintf("socks://%s@%s:%s#%s", secret, address, config.DefaultSocksPort, SocksProxyName)
}

func defaultShadowsocks(address string, alterID string) string {

	methodAndPassword := fmt.Sprintf("%s:%s", config.DefaultShadowsocksMethod, config.DefaultShadowsocksPassword)
	bs64 := base64.StdEncoding.EncodeToString([]byte(methodAndPassword))
	return fmt.Sprintf("ss://%s@%s:%s#%s", bs64, address, config.DefaultShadowsocksPort, ShadowsocksProxyName)
}

const (
	ProxyVMess       string = "vmess"
	ProxySocks       string = "socks"
	ProxyShadowsocks string = "shadowsocks"
	ProxyAll         string = "all"
)

var subscribeMapFunc = map[string]func(string, string) string{
	ProxyVMess:       defaultVMess,
	ProxySocks:       defaultSocks,
	ProxyShadowsocks: defaultShadowsocks,
}

func Subscribe(proxyType string, address string, alterId string) string {

	urlFunc := subscribeMapFunc[proxyType]
	subURL := urlFunc(address, alterId)
	return subURL

}
