package main

import (
	"fmt"
	"os"
	portalgateway "portal-gateway"
	"portal-gateway/config"
	"portal-gateway/log"
)

func main() {
	var logLevel string

	config, err := config.LoadConfig("portal.yaml")
	if err != nil {
		fmt.Printf("failed to load configuration: %v", err)
		os.Exit(1)
	}
	if config.LoggingConfig.Level != "" && logLevel == "" {
		logLevel = config.LoggingConfig.Level
	} else {
		logLevel = "info"
	}
	logger, err := log.NewDefaultLogger(log.ParseLevel(logLevel))
	if err != nil {
		fmt.Println("failed to initialize logger")
		os.Exit(1)
	}

	gateway, err := portalgateway.NewPortal(config, logger)
	if err != nil {
		logger.Fatalf("failed to create gateway: %v", err)
	}
	logger.Fatal(gateway.Start())
}
