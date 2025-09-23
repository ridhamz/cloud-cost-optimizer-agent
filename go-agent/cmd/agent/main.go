package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/analyzer"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/aws"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/executor"
	"github.com/ridhamz/AI-cloud-cost-optimizer-agent/internal/recommender"
	"gopkg.in/yaml.v3"
)

// Structs for YAML
type Config struct {
	AWS struct {
		Region   string   `yaml:"region"`
		Services []string `yaml:"services"`
	} `yaml:"aws"`

	Thresholds struct {
		EC2CPUUtilization float64 `yaml:"ec2_cpu_utilization"`
		RDSCPUUtilization float64 `yaml:"rds_cpu_utilization"`
	} `yaml:"thresholds"`

	ProtectedResources []string `yaml:"protected_resources"`
}

func main() {
	fmt.Println("🚀 Starting AI-Powered Cloud Cost Optimizer Agent...")

	// Read YAML file
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config.yaml: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Failed to parse config.yaml: %v", err)
	}

	fmt.Println("AWS Region:", cfg.AWS.Region)
	fmt.Println("EC2 CPU threshold:", cfg.Thresholds.EC2CPUUtilization)
	fmt.Println("Protected resources:", cfg.ProtectedResources)

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
