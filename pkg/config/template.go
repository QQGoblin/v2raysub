package config

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var GoblinV2RayConfig = template.Must(template.New("config.json").Parse(dedent.Dedent(`{
  "log": {
    "loglevel": "warning"
  },
  "inbounds": [
    {
      "protocol": "socks",
      "port": 10777,
      "listen": "0.0.0.0",
      "settings": {
        "auth": "password",
        "accounts": [
          {
            "user": "proxy",
            "pass": ""
          }
        ],
        "udp": false,
        "ip": "127.0.0.1"
      }
    },
    {
      "tag": "http",
      "port": 10778,
      "listen": "127.0.0.1",
      "protocol": "http"
    }
  ],
  "outbounds": [
    {
      "tag": "proxy",
      "protocol": "shadowsocks",
      "settings": {
        "servers": [
          {
            "address": "{{ .Address }}",
            "method": "chacha20-ietf-poly1305",
            "ota": false,
            "password": "",
            "port": 28398,
            "level": 1
          }
        ]
      },
      "streamSettings": {
        "network": "tcp"
      },
      "mux": {
        "enabled": false,
        "concurrency": -1
      }
    },
    {
      "tag": "direct",
      "protocol": "freedom",
      "settings": {}
    }
  ],
  "routing": {
    "domainStrategy": "AsIs",
    "rules": [
      {
        "id": "1",
        "type": "field",
        "outboundTag": "direct",
        "domain": [
          "geosite:cn",
          "geoip:cn"
        ],
        "enabled": true
      }
    ]
  }
}`)))
