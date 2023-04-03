package aliyun

import (
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type ECS struct {
	Name            string   `json:"name"`
	PublicIpAddress []string `json:"public_ip_address"`
	Status          string   `json:"status"`
}

func ListInstances(regionID, key, secret, endpoint string) ([]*ECS, error) {

	config := &openapi.Config{
		AccessKeyId:     &key,
		AccessKeySecret: &secret,
		Endpoint:        tea.String(endpoint),
	}

	client, err := ecs.NewClient(config)
	if err != nil {
		return nil, err
	}
	r := &ecs.DescribeInstancesRequest{
		RegionId: tea.String(regionID),
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
