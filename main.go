package main

import (
	"bytes"
	"github.com/QQGoblin/v2raysub/pkg/aliyun"
	"github.com/QQGoblin/v2raysub/pkg/contants"
	"github.com/QQGoblin/v2raysub/pkg/v2ray"
	"io"
	"net/http"
)

func subscribe(w http.ResponseWriter, r *http.Request, proxyType v2ray.ProxyType) {

	ecss, err := aliyun.ListInstances(contants.Region, contants.Key, contants.Secret, contants.Endpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, ecs := range ecss {
		if ecs.Status != "Running" {
			continue
		}
		urls := v2ray.Subscribe(ecs.PublicIpAddress[0], proxyType)
		if _, err := io.Copy(w, bytes.NewReader([]byte(urls))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "no running ecs", http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/sub", func(writer http.ResponseWriter, request *http.Request) {
		subscribe(writer, request, v2ray.ProxyTypeVMess)
	})
	http.ListenAndServe("0.0.0.0:80", nil)
}
