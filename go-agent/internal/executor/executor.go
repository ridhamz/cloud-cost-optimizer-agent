package executor

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	appConfig "github.com/ridhamz/AI-cloud-cost-optimizer-agent/configs"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/recommender"
)

// ExecuteActions generates a summary of recommendations and sends it via SNS
func ExecuteActions(recs []recommender.Recommendation) error {
	if len(recs) == 0 {
		fmt.Println("ℹ️ No recommendations to execute.")
		return nil
	}

	// Build a human-readable report
	report := buildReport(recs)

	// Send report to SNS
	if err := sendReportToSNS(report); err != nil {
		return fmt.Errorf("failed to send SNS report: %w", err)
	}

	fmt.Println("📧 Report successfully sent via SNS.")
	return nil
}

// buildReport formats recommendations into a clear text report
func buildReport(recs []recommender.Recommendation) string {
	var sb strings.Builder

	sb.WriteString("📊 AI-Powered Cloud Cost Optimizer Report\n")
	sb.WriteString("=========================================\n\n")

	for i, r := range recs {
		sb.WriteString(fmt.Sprintf("🔹 Recommendation #%d\n", i+1))
		sb.WriteString(fmt.Sprintf("Resource ID: %s\n", r.ResourceID))
		sb.WriteString(fmt.Sprintf("Service: %s\n", r.Type))
		sb.WriteString(fmt.Sprintf("Action: %s\n", r.Action))
		sb.WriteString(fmt.Sprintf("Estimated Monthly Savings: $%.2f\n", r.Savings))
		sb.WriteString(fmt.Sprintf("Reason: %s\n", r.Reason))
		sb.WriteString("-----------------------------------------\n")
	}

	sb.WriteString("\n✅ End of report.\n")
	return sb.String()
}

// sendReportToSNS publishes the report to an SNS topic
func sendReportToSNS(report string) error {
	// Load AWS SDK config (uses ~/.aws/credentials or IAM role)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(appConfig.AppConfig.AWS.Region))
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	snsClient := sns.NewFromConfig(cfg)

	_, err = snsClient.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(appConfig.AppConfig.AWS.SNSTopicARN),
		Subject:  aws.String("AI Cloud Cost Optimizer Report"),
		Message:  aws.String(report),
	})
	if err != nil {
		return fmt.Errorf("SNS publish failed: %w", err)
	}

	return nil
}
