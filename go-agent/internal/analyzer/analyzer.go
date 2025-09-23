package analyzer

import (
	"fmt"

	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/aws"
)

type Analysis struct {
	ResourceID     string
	UnderUtilized  bool
	Savings        float64
	Recommendation string
}

// AnalyzeResources uses Claude AI to generate recommendations
func AnalyzeResources(resources []aws.Resource) []Analysis {
	var results []Analysis

	for _, r := range resources {
		// Create a prompt for Claude
		prompt := fmt.Sprintf(
			"Analyze the following AWS resource and give cost optimization advice:\nResource ID: %s\nType: %s\nUsage: %.2f\nCost: %.2f",
			r.ID, r.Type, r.Usage, r.Cost,
		)

		recommendation, err := CallClaude(prompt)
		if err != nil {
			recommendation = "Error generating AI recommendation"
		}

		results = append(results, Analysis{
			ResourceID:     r.ID,
			UnderUtilized:  r.Usage < 30, // simple threshold
			Savings:        r.Cost * 0.5, // mock savings
			Recommendation: recommendation,
		})
	}

	return results
}
