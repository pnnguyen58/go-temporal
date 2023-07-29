package config

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

type TempoConfig struct {
	HostPort string `json:"hostPort"`
	Namespace *Namespace `json:"namespace"`
	Workflow *Workflow `json:"workflow"`
}

type Namespace struct {
	Name string `json:"namespace"`
	WorkflowExecutionRetentionPeriod time.Duration `json:"workflowExecutionRetentionPeriod"` // seconds
}

type Workflow struct {
	TaskQueueName string `json:"taskQueueName"`
	SearchAttributes map[string]interface{} `json:"searchAttributes"`
	ExecutionTimeout time.Duration `json:"executionTimeout"` // seconds
	RunTimeout time.Duration `json:"runTimeout"` // seconds
	TaskTimeout time.Duration `json:"taskTimeout"` // seconds
	ScheduleToCloseTimeout time.Duration `json:"scheduleToCloseTimeout"` // seconds
	StartToCloseTimeout time.Duration `json:"startToCloseTimeout"` // seconds
	HeartbeatTimeout time.Duration `json:"heartbeatTimeout"` // seconds
	WaitForCancellation bool `json:"waitForCancellation"` // seconds
}

func LoadTempoConfig(ctx context.Context) map[string]TempoConfig {
	return mockTempoConfig()
	// TODO: create secret
	data, err := getSecret(ctx, SecretSuite{
		Region:     C.Secret.Region,
		SecretName: C.Secret.Namespace,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	// map [use-case]config
	ns := make(map[string]TempoConfig)
	err = json.Unmarshal(data, &ns)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return ns
}

func mockTempoConfig() map[string]TempoConfig {
	ns := make(map[string]TempoConfig)
	ns["tempo-example-app"] = TempoConfig{
		HostPort: C.Server.TempoHost,
		Namespace: &Namespace{
			Name: "tempo-example-app",
		},
		Workflow: &Workflow{
			TaskQueueName: "tempo-example-app",
		},
	}
	return ns
}