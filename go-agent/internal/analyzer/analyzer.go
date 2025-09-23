package analyzer

import "github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/aws"

type Analysis struct {
	ResourceID    string
	UnderUtilized bool
	Savings       float64
	Confidence    float64 // AI confidence score
}

// Fake AI scoring for demonstration
func AnalyzeResources(resources []aws.Resource) []Analysis {
	var result []Analysis
	for _, r := range resources {
		underUtilized := false
		confidence := 0.0
		// Simulate AI detecting underutilization
		if r.Type == "EC2" {
			underUtilized = true // pretend AI detected low usage
			confidence = 0.85
		}

		result = append(result, Analysis{
			ResourceID:    r.ID,
			UnderUtilized: underUtilized,
			Savings:       50.0, // mock value
			Confidence:    confidence,
		})
	}
	return result
}
