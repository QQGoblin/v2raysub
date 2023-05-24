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
	InstanceId      string   `json:"instance_id"`
}

func NewClient(c *Config) (*ecs.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     &c.Key,
		AccessKeySecret: &c.Secret,
		Endpoint:        tea.String(c.Endpoint),
	}

	return ecs.NewClient(config)
}

func ListInstances(client *ecs.Client, region string) ([]*ECS, error) {

	r := &ecs.DescribeInstancesRequest{
		RegionId: tea.String(region),
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
			InstanceId:      *ins.InstanceId,
			Name:            *ins.InstanceName,
			PublicIpAddress: pubIPs,
			Status:          *ins.Status,
		})
	}
	return ecss, nil
}

func StopAllInstances(client *ecs.Client, region string) error {

	ecss, err := ListInstances(client, region)
	if err != nil {
		return err
	}

	forceStop := false
	stoppedMode := "StopCharging"

	for _, ins := range ecss {
		if ins.Status != "Running" {
			continue
		}
		r := &ecs.StopInstanceRequest{
			ForceStop:   &forceStop,
			InstanceId:  &ins.InstanceId,
			StoppedMode: &stoppedMode,
		}
		resp, err := client.StopInstance(r)
		if err != nil {
			return err
		}
		if *resp.StatusCode != 200 {
			return errors.New("error status code")
		}

	}

	return nil
}
