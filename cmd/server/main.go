package main

import (
	"bytes"
	"github.com/QQGoblin/go-sdk/pkg/tmpl"
	"github.com/QQGoblin/v2raysub/pkg/aliyun"
	"github.com/QQGoblin/v2raysub/pkg/v2ray"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"os"
)

type SubConfig struct {
	V2ray  *v2ray.Config  `yaml:"v2ray"`
	Aliyun *aliyun.Config `yaml:"aliyun"`
}

var subConfig = &SubConfig{}

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s, err := v2ray.Subscribe(subConfig.V2ray.Subscribe, dynamicAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, bytes.NewReader([]byte(s))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, "no running ecs", http.StatusInternalServerError)
}

func refresh(w http.ResponseWriter, r *http.Request) {

	dynamicAddress, err := aliyun.PublicAddress(subConfig.Aliyun)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = v2ray.WriteConfig(subConfig.V2ray, tmpl.Data{
		"Address": dynamicAddress,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {
	loadConfig("/v2ray/v2raysub.yaml")
	http.HandleFunc("/v2ray/sub", subscribe)
	http.HandleFunc("/v2ray/refresh", refresh)
	http.ListenAndServe("0.0.0.0:80", nil)
}
