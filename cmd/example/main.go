package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"
)

func main() {
	// Create a new client
	client, err := tvwsclient.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Example symbols
	symbols := []string{
		"NASDAQ:AAPL",
		"NASDAQ:MSFT",
		"NASDAQ:GOOGL",
		"NASDAQ:AMZN",
		"NASDAQ:META",
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create data channel
	dataChan := make(chan map[string]interface{})

	// Start receiving data
	go func() {
		if err := client.GetLatestTradeInfo(symbols, dataChan); err != nil {
			log.Printf("Error: %v", err)
		}
	}()

	// Main loop
	for {
		select {
		case data := <-dataChan:
			fmt.Printf("Received: %+v\n", data)
		case <-sigChan:
			fmt.Println("\nShutting down...")
			return
		}
	}
}
