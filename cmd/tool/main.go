package main

import (
	"flag"
	"github.com/QQGoblin/v2raysub/pkg/aliyun"
	"github.com/QQGoblin/v2raysub/pkg/config"
	"k8s.io/klog/v2"
)

var (
	v2rayconfig string
)

func init() {
	flag.StringVar(&v2rayconfig, "config", "/root/v2ray/v2raysub.yaml", "")
}

func main() {

	flag.Parse()

	subConfig := config.LoadConfigOrDie(v2rayconfig)

	client, err := aliyun.NewClient(subConfig.Aliyun)
	if err != nil {
		klog.Fatalf("create aliyun client failed: %v", err)
	}

	if err != aliyun.StopAllInstances(client, subConfig.Aliyun.Region) {
		klog.Fatalf("stop instances failed: %v", err)
	}

	klog.Info("stop all instances success")

}
