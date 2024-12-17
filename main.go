package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	rtd := NewRealTimeData()

	// Example US stock symbols (NASDAQ)
	exchangeSymbols := []string{
		"NASDAQ:AAPL",  // Apple
		"NASDAQ:MSFT",  // Microsoft
		"NASDAQ:GOOGL", // Google
		"NASDAQ:AMZN",  // Amazon
		"NASDAQ:META",  // Meta (Facebook)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel for receiving data
	dataChan := make(chan map[string]interface{})

	// Start getting data in a goroutine
	go func() {
		if err := rtd.GetLatestTradeInfo(exchangeSymbols, dataChan); err != nil {
			log.Printf("Error getting trade info: %v", err)
		}
	}()

	// Main loop
	for {
		select {
		case data := <-dataChan:
			fmt.Println("----------------------------------------")
			fmt.Printf("%+v\n", data)
		case <-sigChan:
			fmt.Println("\nReceived interrupt signal. Shutting down...")
			rtd.Close()
			return
		}
	}
}
