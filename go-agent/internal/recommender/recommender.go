package recommender

import "github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/analyzer"

// Recommendation is the final actionable advice derived from analysis
type Recommendation struct {
	ResourceID string  // AWS Resource identifier
	Type       string  // AWS service type (EC2, RDS, etc.)
	Action     string  // Suggested action (e.g., "Stop instance", "Downsize DB")
	Savings    float64 // Estimated monthly savings
	Reason     string  // Why this action is recommended
}

// GenerateRecommendations transforms analysis results into actionable recommendations
func GenerateRecommendations(analysis []analyzer.Analysis) []Recommendation {
	var recs []Recommendation
	for _, a := range analysis {
		if a.UnderUtilized {
			recs = append(recs, Recommendation{
				ResourceID: a.ResourceID,
				Type:       a.Type,
				Action:     a.Recommendation,
				Savings:    a.Savings,
				Reason:     "Resource is underutilized based on threshold analysis",
			})
		}
	}
	return recs
}
