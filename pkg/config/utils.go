package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
)

func LoadConfigOrDie(path string) *SubConfig {

	cf, err := os.Open(path)
	if err != nil {
		klog.Fatalf("load config file from %s failed, err %v", path, err)
	}
	defer cf.Close()

	bs, err := ioutil.ReadAll(cf)
	if err != nil {
		klog.Fatalf("read config file from %s failed, err %v", path, err)
	}

	config := SubConfig{}
	if err := yaml.Unmarshal(bs, &config); err != nil {
		klog.Fatalf("unmarshal config failed, content: %v，", bs)
	}
	return &config
}
