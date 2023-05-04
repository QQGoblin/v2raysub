package main

import (
	"bytes"
	"context"
	"github.com/QQGoblin/go-sdk/pkg/tmpl"
	"github.com/QQGoblin/v2raysub/pkg/aliyun"
	"github.com/QQGoblin/v2raysub/pkg/v2ray"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"os/exec"
)

type SubConfig struct {
	V2ray  *v2ray.Config  `yaml:"v2ray"`
	Aliyun *aliyun.Config `yaml:"aliyun"`
}

var (
	subConfig   = &SubConfig{}
	cancelV2Ray context.CancelFunc
)

func loadConfig(path string) {

	cf, err := os.Open(path)
	if err != nil {
		klog.Fatalf("load config file from %s failed, err %v", path, err)
	}
	defer cf.Close()

	bs, err := ioutil.ReadAll(cf)
	if err != nil {
		klog.Fatalf("read config file from %s failed, err %v", path, err)
	}

	if err := yaml.Unmarshal(bs, subConfig); err != nil {
		klog.Fatalf("unmarshal config failed, content: %v，", bs)
	}
}

func subscribe(w http.ResponseWriter, r *http.Request) {

	dynamicAddress, err := aliyun.PublicAddress(subConfig.Aliyun)
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
	dynamicAddress, err := aliyun.PublicAddress(subConfig.Aliyun)
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
	loadConfig("/v2ray/v2raysub.yaml")

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
