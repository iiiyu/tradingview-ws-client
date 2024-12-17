package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type RealTimeData struct {
	ws            *websocket.Conn
	requestHeader http.Header
	wsURL         string
	validateURL   string
}

func NewRealTimeData() *RealTimeData {
	header := http.Header{
		"Accept-Encoding": {"gzip, deflate, br, zstd"},
		"Accept-Language": {"en-US,en;q=0.9,fa;q=0.8"},
		"Cache-Control":   {"no-cache"},
		"Origin":          {"https://www.tradingview.com"},
		"Pragma":          {"no-cache"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"},
	}

	rtd := &RealTimeData{
		requestHeader: header,
		wsURL:         "wss://data.tradingview.com/socket.io/websocket?from=screener%2F",
		validateURL:   "https://scanner.tradingview.com/symbol?symbol=%s%%3A%s&fields=market&no_404=false",
	}

	conn, _, err := websocket.DefaultDialer.Dial(rtd.wsURL, header)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	rtd.ws = conn

	return rtd
}

func (rtd *RealTimeData) Close() {
	if rtd.ws != nil {
		rtd.ws.Close()
	}
}

func (rtd *RealTimeData) validateSymbols(exchangeSymbols []string) error {
	if len(exchangeSymbols) == 0 {
		return fmt.Errorf("exchange_symbol could not be empty")
	}

	for _, item := range exchangeSymbols {
		parts := strings.Split(item, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid symbol format '%s'. Must be like 'BINANCE:BTCUSDT'", item)
		}

		exchange, symbol := parts[0], parts[1]
		url := fmt.Sprintf(rtd.validateURL, exchange, symbol)

		retries := 3
		for attempt := 0; attempt < retries; attempt++ {
			resp, err := http.Get(url)
			if err != nil {
				if attempt == retries-1 {
					return fmt.Errorf("failed to validate symbol '%s' after %d attempts: %v", item, retries, err)
				}
				time.Sleep(time.Second)
				continue
			}
			resp.Body.Close()

			if resp.StatusCode == 404 {
				return fmt.Errorf("invalid symbol '%s'", item)
			}
			if resp.StatusCode != http.StatusOK {
				if attempt == retries-1 {
					return fmt.Errorf("failed to validate symbol '%s' after %d attempts", item, retries)
				}
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}
	return nil
}
