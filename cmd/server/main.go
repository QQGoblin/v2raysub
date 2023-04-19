package main

import (
	"bytes"
	"github.com/QQGoblin/v2raysub/pkg/aliyun"
	"github.com/QQGoblin/v2raysub/pkg/config"
	"github.com/QQGoblin/v2raysub/pkg/v2ray"
	"io"
	"net/http"
)

func subscribe(w http.ResponseWriter, r *http.Request) {

	ecss, err := aliyun.ListInstances(config.Region, config.Key, config.Secret, config.Endpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	proxyType := r.URL.Query().Get("type")
	if proxyType == "" {
		proxyType = v2ray.ProxyVMess
	}
	alterid := r.URL.Query().Get("alterid")
	for _, ecs := range ecss {
		if ecs.Status != "Running" {
			continue
		}
		urls := v2ray.Subscribe(proxyType, ecs.PublicIpAddress[0], alterid)
		if _, err := io.Copy(w, bytes.NewReader([]byte(urls))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "no running ecs", http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/sub", subscribe)
	http.ListenAndServe("0.0.0.0:80", nil)
}
