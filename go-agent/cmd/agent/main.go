package main

import (
	"fmt"
	"log"

	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/analyzer"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/aws"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/executor"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/recommender"
)

func main() {
	fmt.Println("🚀 Starting AI-Powered Cloud Cost Optimizer Agent...")

	// Step 1: Collect AWS resource data
	resources, err := aws.FetchEC2Instances()
	if err != nil {
		log.Fatalf("Error fetching AWS resources: %v", err)
	}

	// Step 2: Analyze resources
	analysis := analyzer.AnalyzeResources(resources)

	// Step 3: Generate recommendations
	recs := recommender.GenerateRecommendations(analysis)

	// Step 4: Execute safe actions
	executor.ExecuteActions(recs)

	fmt.Println("✅ Optimization run complete!")
}
