package main

import (
	"errors"
	"flag"
	"github.com/QQGoblin/go-sdk/pkg/tmpl"
	"github.com/QQGoblin/v2raysub/pkg/aliyun"
	"github.com/QQGoblin/v2raysub/pkg/config"
	"io/ioutil"
	"os"
	"path"
)

var (
	output    string
	aliRegion string
	aliSecret string
	aliKey    string
)

func init() {
	flag.StringVar(&output, "output", "/root/v2ray/config.json", "")
	flag.StringVar(&aliRegion, "region", "", "")
	flag.StringVar(&aliSecret, "secret", "", "")
	flag.StringVar(&aliKey, "key", "", "")
}

func generateConfig(filename string, data tmpl.Data) error {

	content, err := tmpl.Render(config.GoblinV2RayConfig, data)

	if err != nil {
		return err
	}
	dir, _ := path.Split(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, []byte(content), 0755)
}

func address() (string, error) {
	ecss, err := aliyun.ListInstances(aliRegion, aliKey, aliSecret, config.Endpoint)
	if err != nil {
		return "", err
	}
	var currentAddress = ""
	for _, ecs := range ecss {
		if ecs.Status != "Running" {
			continue
		}
		currentAddress = ecs.PublicIpAddress[0]
	}

	if currentAddress == "" {
		return "", errors.New("no running ecs on aliyun")
	}
	return currentAddress, nil
}

func main() {
	flag.Parse()
	addr, err := address()
	if err != nil {
		panic("no running ecs on aliyun")
	}
	if err := generateConfig(output, tmpl.Data{
		"Address": addr,
	}); err != nil {
		panic("write config failed")
	}
}
