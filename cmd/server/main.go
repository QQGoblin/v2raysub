package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/QQGoblin/go-sdk/pkg/tmpl"
	"github.com/QQGoblin/v2raysub/pkg/aliyun"
	"github.com/QQGoblin/v2raysub/pkg/config"
	"github.com/QQGoblin/v2raysub/pkg/v2ray"
	"io"
	"k8s.io/klog/v2"
	"net/http"
	"os/exec"
)

var (
	subConfig   = &config.SubConfig{}
	cancelV2Ray context.CancelFunc
)

func getProxyServerAddress(c *aliyun.Config) (string, error) {

	client, err := aliyun.NewClient(c)
	if err != nil {
		return "", err
	}

	ecss, err := aliyun.ListInstances(client, c.Region)
	if err != nil {
		return "", err
	}

	for _, ecs := range ecss {
		if ecs.Status != "Running" {
			continue
		}
		return ecs.PublicIpAddress[0], nil
	}
	return "", errors.New("not running ecs")
}

func subscribe(w http.ResponseWriter, r *http.Request) {

	dynamicAddress, err := getProxyServerAddress(subConfig.Aliyun)
	if err != nil {
		klog.Errorf("get publicAddress failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s, err := v2ray.Subscribe(subConfig.V2ray.Subscribe, dynamicAddress)
	if err != nil {
		klog.Errorf("create subscribe string failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, bytes.NewReader([]byte(s))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func refresh(w http.ResponseWriter, r *http.Request) {

	klog.Info("refresh local v2ray config")
	dynamicAddress, err := getProxyServerAddress(subConfig.Aliyun)
	if err != nil {
		klog.Errorf("get publicAddress failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = v2ray.WriteConfig(subConfig.V2ray, tmpl.Data{
		"Address": dynamicAddress,
	}); err != nil {
		klog.Errorf("write local config failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		ctx context.Context
	)

	if cancelV2Ray != nil {
		cancelV2Ray()
	}
	ctx, cancelV2Ray = context.WithCancel(context.Background())
	if err = run(ctx, "/usr/bin/v2ray", subConfig.V2ray.Path); err != nil {
		klog.Errorf("start local v2ray failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = io.Copy(w, bytes.NewReader([]byte("refresh success"))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func run(ctx context.Context, bin, config string) error {
	v2rayCmd := exec.CommandContext(ctx, bin, "run", "-c", config)
	return v2rayCmd.Start()
}

func main() {

	subConfig = config.LoadConfigOrDie("/v2ray/v2raysub.yaml")

	var (
		ctx context.Context
	)
	ctx, cancelV2Ray = context.WithCancel(context.Background())
	if err := run(ctx, "/usr/bin/v2ray", subConfig.V2ray.Path); err != nil {
		klog.Errorf("start local v2ray failed: %v", err)
	}

	http.HandleFunc("/v2ray/sub", subscribe)
	http.HandleFunc("/v2ray/refresh", refresh)
	http.ListenAndServe("0.0.0.0:18088", nil)
}
