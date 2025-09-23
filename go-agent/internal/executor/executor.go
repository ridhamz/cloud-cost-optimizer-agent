package executor

import (
	"context"
	"fmt"
	"strings"

	sdkAws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	localAws "github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/aws"
)

// SendReportSNS sends an AWS resource report via SNS
func SendReportSNS(cfg sdkAws.Config, resources []localAws.Resource, topicARN string) error {
	snsClient := sns.NewFromConfig(cfg)

	// Generate report body
	body := generateReportBody(resources)

	input := &sns.PublishInput{
		TopicArn: sdkAws.String(topicARN),
		Message:  sdkAws.String(body),
		Subject:  sdkAws.String("AWS Resource Usage Report with Recommendations"),
	}

	_, err := snsClient.Publish(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to publish SNS message: %v", err)
	}

	fmt.Println("SNS report sent successfully!")
	return nil
}

// generateReportBody creates a simple text report including AI recommendations
func generateReportBody(resources []localAws.Resource) string {
	var sb strings.Builder
	sb.WriteString("AWS Resource Usage Report with Recommendations\n\n")
	sb.WriteString("ID\tType\tUsage\tCost\tRegion\tRecommendation\n")
	for _, r := range resources {
		sb.WriteString(fmt.Sprintf("%s\t%s\t%.2f\t%.2f\t%s\t%s\n",
			r.ID, r.Type, r.Usage, r.Cost, r.Region, r.Recommendation))
	}
	return sb.String()
}
