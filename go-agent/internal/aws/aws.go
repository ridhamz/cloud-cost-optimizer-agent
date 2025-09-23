package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Resource struct {
	ID     string
	Type   string
	Region string
	Cost   float64
	Usage  float64
}

func FetchEC2Instances() ([]Resource, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println("Error loading AWS config:", err)
		return nil, err
	}

	svc := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(context.TODO(), input)
	if err != nil {
		log.Println("Error describing instances:", err)
		return nil, err
	}

	var resources []Resource
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			resources = append(resources, Resource{
				ID:     *instance.InstanceId,
				Type:   "EC2",
				Region: cfg.Region,
			})
		}
	}

	return resources, nil
}
