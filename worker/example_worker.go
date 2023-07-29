package main

import (
	"github.com/pnnguyen58/go-temporal/config"
	"github.com/pnnguyen58/go-temporal/core/activities"
	"github.com/pnnguyen58/go-temporal/core/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func registerExample(conf config.TempoConfig) {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort:  conf.HostPort,
		Namespace: conf.Namespace.Name,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	taskQueueName := conf.Workflow.TaskQueueName
	w := worker.New(c, taskQueueName, worker.Options{})

	w.RegisterWorkflow(workflows.CreateExampleWorkflow)
	w.RegisterActivity(activities.CreateExample)
	w.RegisterActivity(activities.CreateExampleCompensation)
	// TODO: add more workflows and activities

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
