package tvwsclient

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// GetLatestTradeInfo starts streaming real-time trade information for the given symbols
func (c *Client) GetLatestTradeInfo(symbols []string, dataChan chan<- map[string]interface{}) error {
	// Initial protocol and auth messages
	initMessages := []string{
		`{"m":"set_auth_token","p":["unauthorized_user_token"]}`,
		`{"m":"set_locale","p":["en","US"]}`,
		`{"m":"chart_create_session","p":["cs_` + GenerateSession("") + `",""]}`}

	// Send initial protocol messages
	for _, msg := range initMessages {
		wrappedMsg := fmt.Sprintf("~m~%d~m~%s", len(msg), msg)
		if err := c.ws.WriteMessage(websocket.TextMessage, []byte(wrappedMsg)); err != nil {
			return fmt.Errorf("error sending init message: %w", err)
		}
	}

	// Wait a bit for initialization
	time.Sleep(1 * time.Second)

	// Create quote session
	session := GenerateSession("qs_")
	quoteMessages := []string{
		fmt.Sprintf(`{"m":"quote_create_session","p":["%s"]}`, session),
		fmt.Sprintf(`{"m":"quote_set_fields","p":["%s","lp","ch","chp","current_session","description","local_description","language","exchange","fractional","is_tradable","lp_time","minmov","minmove2","original_name","pricescale","pro_name","short_name","type","update_mode","volume","currency_code","rchp","rtc","status"]}`, session),
		fmt.Sprintf(`{"m":"quote_add_symbols","p":["%s","%s"]}`, session, strings.Join(symbols, `","`)),
	}

	// Send quote session messages
	for _, msg := range quoteMessages {
		wrappedMsg := fmt.Sprintf("~m~%d~m~%s", len(msg), msg)
		if err := c.ws.WriteMessage(websocket.TextMessage, []byte(wrappedMsg)); err != nil {
			return fmt.Errorf("error sending quote message: %w", err)
		}
		// Small delay between messages
		time.Sleep(100 * time.Millisecond)
	}

	// Read messages in a loop
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil
			}
			return fmt.Errorf("error reading message: %w", err)
		}

		// Handle heartbeat messages
		if heartbeatRegex.Match(message) {
			if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				return fmt.Errorf("error sending heartbeat response: %w", err)
			}
			continue
		}

		// Process data messages
		parts := strings.Split(string(message), "~m~")
		for _, part := range parts {
			if strings.HasPrefix(part, "{") {
				var response TVResponse
				if err := json.Unmarshal([]byte(part), &response); err != nil {
					continue
				}

				// Only process quote data messages
				if response.Method == "qsd" && len(response.Params) >= 2 {
					// Extract the quote data from params
					if quoteDataRaw, err := json.Marshal(response.Params[1]); err == nil {
						var quote QuoteData
						if err := json.Unmarshal(quoteDataRaw, &quote); err == nil {
							// Convert to map for compatibility with existing channel
							dataMap := map[string]interface{}{
								"m": response.Method,
								"p": []interface{}{response.Params[0], quote},
							}
							dataChan <- dataMap
						}
					}
				}
			}
		}
	}
}
