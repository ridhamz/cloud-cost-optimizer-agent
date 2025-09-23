package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
)

// Resource represents an AWS resource with usage metrics
type Resource struct {
	ID             string
	Type           string
	Usage          float64
	Cost           float64
	Region         string
	Recommendation string
}

// FetchResources fetches EC2, RDS, and Lambda resources with CloudWatch metrics
func FetchResources(region string) ([]Resource, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	cw := cloudwatch.NewFromConfig(cfg)
	var resources []Resource

	///////////////////
	// EC2 Instances
	///////////////////
	ec2Svc := ec2.NewFromConfig(cfg)
	ec2Resp, err := ec2Svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("error describing EC2: %v", err)
	}

	for _, r := range ec2Resp.Reservations {
		for _, i := range r.Instances {
			usage := getAverageCPU("AWS/EC2", "CPUUtilization", "InstanceId", *i.InstanceId, cw)
			resources = append(resources, Resource{
				ID:     *i.InstanceId,
				Type:   "EC2",
				Usage:  usage,
				Cost:   estimateEC2Cost(i),
				Region: region,
			})
		}
	}

	///////////////////
	// RDS Instances
	///////////////////
	rdsSvc := rds.NewFromConfig(cfg)
	rdsResp, err := rdsSvc.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("error describing RDS: %v", err)
	}

	for _, db := range rdsResp.DBInstances {
		usage := getAverageCPU("AWS/RDS", "CPUUtilization", "DBInstanceIdentifier", *db.DBInstanceIdentifier, cw)
		resources = append(resources, Resource{
			ID:     *db.DBInstanceIdentifier,
			Type:   "RDS",
			Usage:  usage,
			Cost:   estimateRDSCost(db),
			Region: region,
		})
	}

	///////////////////
	// Lambda Functions
	///////////////////
	lambdaSvc := lambda.NewFromConfig(cfg)
	lambdaResp, err := lambdaSvc.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{})
	if err != nil {
		return nil, fmt.Errorf("error listing Lambda: %v", err)
	}

	for _, f := range lambdaResp.Functions {
		invocations := getLambdaInvocations(*f.FunctionName, cw)
		resources = append(resources, Resource{
			ID:     *f.FunctionName,
			Type:   "Lambda",
			Usage:  float64(invocations),
			Cost:   estimateLambdaCost(f),
			Region: region,
		})
	}

	return resources, nil
}

/////////////////////////
// CloudWatch helpers
/////////////////////////

func getAverageCPU(namespace, metric, dimName, dimValue string, cw *cloudwatch.Client) float64 {
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	resp, err := cw.GetMetricStatistics(context.TODO(), &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String(namespace),
		MetricName: aws.String(metric),
		Dimensions: []cwTypes.Dimension{
			{
				Name:  aws.String(dimName),
				Value: aws.String(dimValue),
			},
		},
		StartTime:  &start,
		EndTime:    &end,
		Period:     aws.Int32(3600),
		Statistics: []cwTypes.Statistic{cwTypes.StatisticAverage},
	})
	if err != nil || len(resp.Datapoints) == 0 {
		return 0
	}

	// return the most recent datapoint
	return *resp.Datapoints[len(resp.Datapoints)-1].Average
}

func getLambdaInvocations(functionName string, cw *cloudwatch.Client) int {
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	resp, err := cw.GetMetricStatistics(context.TODO(), &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Invocations"),
		Dimensions: []cwTypes.Dimension{
			{
				Name:  aws.String("FunctionName"),
				Value: aws.String(functionName),
			},
		},
		StartTime:  &start,
		EndTime:    &end,
		Period:     aws.Int32(3600),
		Statistics: []cwTypes.Statistic{cwTypes.StatisticSum},
	})
	if err != nil || len(resp.Datapoints) == 0 {
		return 0
	}
	return int(*resp.Datapoints[len(resp.Datapoints)-1].Sum)
}

/////////////////////////
// Cost Estimators (placeholders)
/////////////////////////

func estimateEC2Cost(i ec2Types.Instance) float64                    { return 50.0 }
func estimateRDSCost(i rdsTypes.DBInstance) float64                  { return 200.0 }
func estimateLambdaCost(f lambdaTypes.FunctionConfiguration) float64 { return 10.0 }
