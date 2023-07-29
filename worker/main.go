package main

import (
	"context"
	"fmt"
	"github.com/pnnguyen58/go-temporal/config"
	"os"
	"os/signal"
)

func main() {
	config.ReadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan)
	go func() {
		// Wait for termination signal
		<-signalChan
		// Trigger cancellation of the context
		cancel()
		// Wait for goroutine to finish
		fmt.Println("The service terminated gracefully")
	}()
	tempoConf := config.LoadTempoConfig(ctx)

	for key, conf := range tempoConf {
		switch key {
		case "tempo-example-app":
			registerExample(conf)
		default:
			fmt.Println("app not defined")
		}
	}
}

