package tvwsclient

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var heartbeatRegex = regexp.MustCompile(`~h~\d+`)

// GetLatestTradeInfo starts streaming real-time trade information for the given symbols
func (c *Client) GetLatestTradeInfo(symbols []string, dataChan chan<- map[string]interface{}) error {
	// Initial protocol and auth messages
	initMessages := []string{
		`{"m":"set_auth_token","p":["unauthorized_user_token"]}`,
		`{"m":"set_locale","p":["en","US"]}`,
		`{"m":"chart_create_session","p":["cs_` + generateSession("") + `",""]}`}

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
	session := generateSession("qs_")
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
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(part), &data); err != nil {
					continue
				}
				// Only forward actual trade data
				if m, ok := data["m"].(string); ok && m == "qsd" {
					dataChan <- data
				}
			}
		}
	}
}

// generateSession generates a random session ID with the given prefix
func generateSession(prefix string) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 12)
	for i := range b {
		b[i] = letterBytes[i%len(letterBytes)]
	}
	return prefix + string(b)
}
