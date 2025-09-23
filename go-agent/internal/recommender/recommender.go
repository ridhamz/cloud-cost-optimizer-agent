package recommender

import "github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/analyzer"

type Recommendation struct {
	ResourceID string
	Action     string
}

func GenerateRecommendations(analysis []analyzer.Analysis) []Recommendation {
	var recs []Recommendation
	for _, a := range analysis {
		if a.UnderUtilized {
			recs = append(recs, Recommendation{
				ResourceID: a.ResourceID,
				Action:     a.Recommendation,
			})
		}
	}
	return recs
}
