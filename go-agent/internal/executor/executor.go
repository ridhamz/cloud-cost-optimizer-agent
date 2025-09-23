package executor

import "github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/recommender"

var protectedResources = map[string]bool{
	"i-prod-123": true,
}

func ExecuteActions(recs []recommender.Recommendation) {
	for _, r := range recs {
		if protectedResources[r.ResourceID] {
			println("⚠️ Skipping protected resource:", r.ResourceID)
			continue
		}
		// TODO: Call AWS SDK or Terraform to apply changes
		println("🔹 Action for resource", r.ResourceID, ":", r.Action)
	}
}
