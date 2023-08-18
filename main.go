package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"os/exec"
)

var (
	cancelV2Ray context.CancelFunc
)

const SubscribeFile = "subscribe"

func subscribe(w http.ResponseWriter, r *http.Request) {

	s, err := ioutil.ReadFile(SubscribeFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = io.Copy(w, bytes.NewReader(s)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func run(ctx context.Context, bin, config string) error {
	v2rayCmd := exec.CommandContext(ctx, bin, "run", "-c", config)
	pipReader, pipWriter := io.Pipe()
	v2rayCmd.Stdout = pipWriter

	go func() {
		_, err := io.Copy(os.Stdout, pipReader)
		if err != nil {
			klog.Errorf("start local v2ray failed: %v", err)
		}
	}()

	return v2rayCmd.Start()
}

func main() {

	var (
		ctx context.Context
	)

	ctx, cancelV2Ray = context.WithCancel(context.Background())
	if err := run(ctx, "v2ray", "config.json"); err != nil {
		klog.Errorf("start local v2ray failed: %v", err)
	}

	http.HandleFunc("/v2ray/sub", subscribe)
	http.ListenAndServe("0.0.0.0:18088", nil)
}
