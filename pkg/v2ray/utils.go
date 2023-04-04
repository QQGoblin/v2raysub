package v2ray

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/QQGoblin/v2raysub/pkg/contants"
	"strings"
)

type VMessVesion string

const (
	VMessVersion2 VMessVesion = "2"
)
const (
	DefaultVMessName = "VMess"
	DefaultSockName  = "Socks"
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

func defaultVMess(address string) string {

	defaultVMESS := vmess{
		Version: VMessVersion2,
		Name:    DefaultVMessName,
		Address: address,
		Port:    contants.DefaultVMessPort,
		ID:      contants.DefaultVMessID,
		AlterId: contants.DefaultAlterID,
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

func defaultSock(address string) string {
	secret := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", contants.DefaultSockUser, contants.DefaultSockPass)))
	return fmt.Sprintf("socks://%s@%s:%s#%s", secret, address, contants.DefaultSockPort, DefaultSockName)
}

type ProxyType string

const (
	ProxyTypeVMess ProxyType = "VMess"
	ProxyTypeSocks ProxyType = "Socks"
	ProxyTypeSS    ProxyType = "Shadowsocks"
	ProxyTypeAll   ProxyType = "All"
)

var subscribeMapFunc = map[ProxyType]func(string) string{
	ProxyTypeVMess: defaultVMess,
	ProxyTypeSocks: defaultSock,
}

func Subscribe(address string, proxyType ProxyType) string {
	all := make([]string, 0)

	for k, v := range subscribeMapFunc {
		if proxyType == ProxyTypeAll || k == proxyType {
			all = append(all, v(address))
		}
	}

	allURLs := strings.Join(all, "\n")
	return base64.StdEncoding.EncodeToString([]byte(allURLs))

}
