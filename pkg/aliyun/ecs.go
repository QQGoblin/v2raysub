package aliyun

import (
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Config struct {
	Endpoint string `yaml:"endpoint"`
	Region   string `yaml:"region"`
	Key      string `yaml:"key"`
	Secret   string `yaml:"secret"`
}

type ECS struct {
	Name            string   `json:"name"`
	PublicIpAddress []string `json:"public_ip_address"`
	Status          string   `json:"status"`
}

func ListInstances(c *Config) ([]*ECS, error) {

	config := &openapi.Config{
		AccessKeyId:     &c.Key,
		AccessKeySecret: &c.Secret,
		Endpoint:        tea.String(c.Endpoint),
	}

	client, err := ecs.NewClient(config)
	if err != nil {
		return nil, err
	}
	r := &ecs.DescribeInstancesRequest{
		RegionId: tea.String(c.Region),
	}

	resp, err := client.DescribeInstancesWithOptions(r, &util.RuntimeOptions{})
	if err != nil {
		return nil, err
	}
	if *resp.StatusCode != 200 {
		return nil, errors.New("error status code")
	}

	ecss := make([]*ECS, 0)
	for _, ins := range resp.Body.Instances.Instance {

		pubIPs := make([]string, 0)
		for _, ip := range ins.PublicIpAddress.IpAddress {
			pubIPs = append(pubIPs, *ip)
		}
		ecss = append(ecss, &ECS{
			Name:            *ins.InstanceName,
			PublicIpAddress: pubIPs,
			Status:          *ins.Status,
		})
	}
	return ecss, nil
}

func PublicAddress(c *Config) (string, error) {
	ecss, err := ListInstances(c)
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
