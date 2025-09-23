package analyzer

import (
	"fmt"

	config "github.com/ridhamz/AI-cloud-cost-optimizer-agent/configs"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/aws"
)

type Analysis struct {
	ResourceID     string
	UnderUtilized  bool
	Savings        float64
	Recommendation string
}

func AnalyzeResources(resources []aws.Resource) []Analysis {
	var results []Analysis

	for _, r := range resources {
		// Simple underutilization check from config thresholds
		underUtilized := false
		threshold := config.AppConfig.Thresholds.EC2CPUUtilization

		if r.Usage < threshold {
			underUtilized = true
		}

		// Prompt Claude for recommendation
		prompt := fmt.Sprintf(
			"Analyze this AWS resource and give cost optimization advice:\nResource ID: %s\nType: %s\nUsage: %.2f\nCost: %.2f",
			r.ID, r.Type, r.Usage, r.Cost,
		)

		recommendation, err := CallClaude(prompt)
		if err != nil {
			recommendation = "Error generating AI recommendation"
		}

		results = append(results, Analysis{
			ResourceID:     r.ID,
			UnderUtilized:  underUtilized,
			Savings:        r.Cost * 0.5,
			Recommendation: recommendation,
		})
	}

	return results
}
