package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pnnguyen58/go-temporal/config"
	"github.com/pnnguyen58/go-temporal/core/app"
	"github.com/pnnguyen58/go-temporal/infra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sdk/proto/example"
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


	ap := fx.New(
		fx.Provide(
			infra.NewLogger,
			infra.NewTemporalClient,
			context.TODO,
			// TODO add all providers
			app.NewExampleApp,
			config.LoadTempoConfig,
		),
		fx.Invoke(
			listenAndServe,
		),
	)
	if err := ap.Start(ctx); err != nil {
		os.Exit(1)
	}
}

func listenAndServe(ctx context.Context, logger *zap.Logger, app app.ExampleApp) {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.C.Server.GRPCPort))
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Create a gRPC server instance
	grpcServer := grpc.NewServer()
	// Register our service with the gRPC server
	example.RegisterExampleServiceServer(grpcServer, &ExampleController{logger: logger, app: app})
	// TODO: register more service here

	// Serve gRPC server
	logger.Info(fmt.Sprintf("Serving gRPC on 0.0.0.0:%v", config.C.Server.GRPCPort))
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			logger.Fatal(err.Error())
		}
	}()

	maxMsgSize := 1024 * 1024 * 20
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("0.0.0.0:%v", config.C.Server.GRPCPort),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize), grpc.MaxCallSendMsgSize(maxMsgSize)),
	)
	if err != nil {
		logger.Fatal(err.Error())
	}

	gwMux := runtime.NewServeMux()
	// Register service handlers
	err = example.RegisterExampleServiceHandler(ctx, gwMux, conn)
	if err != nil {
		logger.Fatal(err.Error())
	}
	// TODO: register more service handlers here

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.C.Server.HTTPPort),
		Handler: gwMux,
	}

	logger.Info(fmt.Sprintf("Serving gRPC-Gateway on port %v", config.C.Server.HTTPPort))
	go func() {
		if err = gwServer.ListenAndServe(); err != nil {
			logger.Fatal(err.Error())
		}
	}()
	// Wait for a signal to shut down the server
	<-ctx.Done()

	// Gracefully stop the server
	grpcServer.GracefulStop()
}

