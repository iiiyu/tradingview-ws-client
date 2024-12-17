package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

func (rtd *RealTimeData) GetLatestTradeInfo(exchangeSymbols []string, dataChan chan<- map[string]interface{}) error {
	if err := rtd.validateSymbols(exchangeSymbols); err != nil {
		return err
	}

	quoteSession := generateSession("qs_")
	chartSession := generateSession("cs_")
	log.Printf("Quote session: %s, Chart session: %s", quoteSession, chartSession)

	if err := rtd.initializeSessions(quoteSession, chartSession); err != nil {
		return err
	}

	if err := rtd.addMultipleSymbolsToSessions(quoteSession, exchangeSymbols); err != nil {
		return err
	}

	return rtd.getData(dataChan)
}

func (rtd *RealTimeData) initializeSessions(quoteSession, chartSession string) error {
	messages := []struct {
		function string
		params   []interface{}
	}{
		{"set_auth_token", []interface{}{"unauthorized_user_token"}},
		{"set_locale", []interface{}{"en", "US"}},
		{"chart_create_session", []interface{}{chartSession, ""}},
		{"quote_create_session", []interface{}{quoteSession}},
		{"quote_set_fields", append([]interface{}{quoteSession}, rtd.getQuoteFields()...)},
		{"quote_hibernate_all", []interface{}{quoteSession}},
	}

	for _, msg := range messages {
		if err := rtd.sendMessage(msg.function, msg.params); err != nil {
			return fmt.Errorf("failed to send message %s: %v", msg.function, err)
		}
	}
	return nil
}

func (rtd *RealTimeData) getQuoteFields() []interface{} {
	return []interface{}{
		"ch", "chp", "current_session", "description", "local_description",
		"language", "exchange", "fractional", "is_tradable", "lp",
		"lp_time", "minmov", "minmove2", "original_name", "pricescale",
		"pro_name", "short_name", "type", "update_mode", "volume",
		"currency_code", "rchp", "rtc",
	}
}

func (rtd *RealTimeData) addMultipleSymbolsToSessions(quoteSession string, exchangeSymbols []string) error {
	resolveSymbol := map[string]interface{}{
		"adjustment":  "splits",
		"currency-id": "USD",
		"session":     "regular",
		"symbol":      exchangeSymbols[0],
	}

	resolveSymbolJSON, err := json.Marshal(resolveSymbol)
	if err != nil {
		return err
	}

	messages := []struct {
		function string
		params   []interface{}
	}{
		{"quote_add_symbols", []interface{}{quoteSession, "=" + string(resolveSymbolJSON)}},
		{"quote_fast_symbols", []interface{}{quoteSession, "=" + string(resolveSymbolJSON)}},
		{"quote_add_symbols", append([]interface{}{quoteSession}, interfaceSlice(exchangeSymbols)...)},
		{"quote_fast_symbols", append([]interface{}{quoteSession}, interfaceSlice(exchangeSymbols)...)},
	}

	for _, msg := range messages {
		if err := rtd.sendMessage(msg.function, msg.params); err != nil {
			return fmt.Errorf("failed to send message %s: %v", msg.function, err)
		}
	}
	return nil
}

func (rtd *RealTimeData) getData(dataChan chan<- map[string]interface{}) error {
	heartbeatRegex := regexp.MustCompile(`~m~\d+~m~~h~\d+$`)

	for {
		_, message, err := rtd.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil
			}
			return fmt.Errorf("error reading message: %v", err)
		}

		if heartbeatRegex.Match(message) {
			if err := rtd.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				return fmt.Errorf("error sending heartbeat response: %v", err)
			}
			continue
		}

		parts := strings.Split(string(message), "~m~")
		for _, part := range parts {
			if strings.HasPrefix(part, "{") {
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(part), &data); err != nil {
					log.Printf("Error unmarshaling message: %v", err)
					continue
				}
				dataChan <- data
			}
		}
	}
}

// Helper function to convert []string to []interface{}
func interfaceSlice(strSlice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(strSlice))
	for i, s := range strSlice {
		interfaceSlice[i] = s
	}
	return interfaceSlice
}
