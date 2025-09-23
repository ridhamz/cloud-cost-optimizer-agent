package main

import (
	"fmt"
	"log"

	config "github.com/ridhamz/AI-cloud-cost-optimizer-agent/configs"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/analyzer"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/aws"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/executor"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/recommender"
)

func main() {

	// Load config once at startup
	config.Load("configs/config.yaml")

	fmt.Println("AWS Region:", config.AppConfig.AWS.Region)
	fmt.Println("Claude API Key:", config.AppConfig.AI.ClaudeAPIKey)
	fmt.Println("EC2 CPU Threshold:", config.AppConfig.Thresholds.EC2CPUUtilization)
	fmt.Println("Protected resources:", config.AppConfig.ProtectedResources)
	fmt.Println("🚀 Starting AI-Powered Cloud Cost Optimizer Agent...")

	// Step 1: Collect AWS resource data
	resources, err := aws.FetchResources(config.AppConfig.AWS.Region)
	if err != nil {
		log.Fatalf("Error fetching AWS resources: %v", err)
	}

	// Step 2: Analyze resources
	analysis := analyzer.AnalyzeResources(resources)

	for _, a := range analysis {
		fmt.Println("Resource:", a.ResourceID)
		fmt.Println("Underutilized:", a.UnderUtilized)
		fmt.Println("AI Recommendation:", a.Recommendation)
		fmt.Println("Estimated Savings:", a.Savings)
		fmt.Println("---------------------------------------------------")
	}

	// Step 3: Generate recommendations
	recs := recommender.GenerateRecommendations(analysis)

	// Step 4: Execute safe actions
	executor.ExecuteActions(recs)

	fmt.Println("✅ Optimization run complete!")
}
