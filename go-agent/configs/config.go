package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Structs for config.yaml
type AWSConfig struct {
	Region      string   `yaml:"region"`
	SNSTopicARN string   `yaml:"sns_topic_arn"`
	Services    []string `yaml:"services"`
}

type Thresholds struct {
	EC2CPUUtilization float64 `yaml:"ec2_cpu_utilization"`
	RDSCPUUtilization float64 `yaml:"rds_cpu_utilization"`
}

type AIConfig struct {
	ClaudeAPIKey        string  `yaml:"claude_api_key"`
	ConfidenceThreshold float64 `yaml:"confidence_threshold"`
}

type Notifications struct {
	SlackWebhookURL string `yaml:"slack_webhook_url"`
	Email           struct {
		From string   `yaml:"from"`
		To   []string `yaml:"to"`
	} `yaml:"email"`
}

type Config struct {
	AWS                AWSConfig     `yaml:"aws"`
	Thresholds         Thresholds    `yaml:"thresholds"`
	ProtectedResources []string      `yaml:"protected_resources"`
	AI                 AIConfig      `yaml:"ai"`
	Notifications      Notifications `yaml:"notifications"`
}

// Global variable to use anywhere
var AppConfig Config

// Load reads the YAML file and populates AppConfig
func Load(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config.yaml: %v", err)
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("Failed to parse config.yaml: %v", err)
	}

	log.Println("✅ Config loaded successfully")
}
