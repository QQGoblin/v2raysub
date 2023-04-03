package main

import (
	"bytes"
	"github.com/QQGoblin/v2raysub/pkg/aliyun"
	"github.com/QQGoblin/v2raysub/pkg/contants"
	"github.com/QQGoblin/v2raysub/pkg/v2ray"
	"io"
	"net/http"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {

	ecss, err := aliyun.ListInstances(contants.Region, contants.Key, contants.Secret, contants.Endpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, ecs := range ecss {
		if ecs.Status != "Running" {
			continue
		}
		vmessURL := v2ray.Subscribe(ecs.PublicIpAddress[0])
		if _, err := io.Copy(w, bytes.NewReader([]byte(vmessURL))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "on running ecs", http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/sub", Subscribe)
	http.ListenAndServe("0.0.0.0:80", nil)
}
