package tvwsclient

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Add these fields as a constant since they're used in quote set fields
const defaultQuoteFields = "lp,ch,chp,current_session,description,local_description,language,exchange," +
	"fractional,is_tradable,lp_time,minmov,minmove2,original_name,pricescale,pro_name," +
	"short_name,type,update_mode,volume,currency_code,rchp,rtc,status"

// sendWSMessage is a helper function that handles the common pattern of sending websocket messages
func sendWSMessage(ws *websocket.Conn, message string, operation string) error {
	wrappedMsg := wrappedMessage(message)
	if err := ws.WriteMessage(websocket.TextMessage, []byte(wrappedMsg)); err != nil {
		return fmt.Errorf("error sending %s: %w", operation, err)
	}
	// Small delay between messages
	time.Sleep(100 * time.Millisecond)
	return nil
}

func SendSetAuthTokenMessage(c *Client, authToken string) error {
	message := fmt.Sprintf(`{"m":"set_auth_token","p":["%s"]}`, authToken)
	return sendWSMessage(c.ws, message, "set auth token message")
}

func SendSetLocalMessage(c *Client) error {
	message := `{"m":"set_locale","p":["en","US"]}`
	return sendWSMessage(c.ws, message, "set local message")
}

func SendChartCreateSessionMessage(c *Client, session string) error {
	message := fmt.Sprintf(`{"m":"set_auth_token","p":["%s",""]}`, session)
	return sendWSMessage(c.ws, message, "chart create session message")
}

func SendQuoteCreateSessionMessage(c *Client, session string) error {
	message := fmt.Sprintf(`{"m":"quote_create_session","p":["%s"]}`, session)
	if err := sendWSMessage(c.ws, message, "quote create session message"); err != nil {
		return err
	}
	return nil
}

func SendQuoteSetFieldsMessage(c *Client, session string) error {
	fields := strings.Split(defaultQuoteFields, ",")
	message := fmt.Sprintf(`{"m":"quote_set_fields","p":["%s",%s]}`,
		session,
		`"`+strings.Join(fields, `","`)+`"`,
	)
	return sendWSMessage(c.ws, message, "quote set fields message")
}

func SendQuoteAddSymbolsMessage(c *Client, session string, symbols []string) error {
	message := fmt.Sprintf(`{"m":"quote_add_symbols","p":["%s","%s"]}`,
		session,
		strings.Join(symbols, `","`),
	)
	return sendWSMessage(c.ws, message, "quote add symbols message")
}

func SendQuoteRemoveSymbolsMessage(c *Client, session string, symbols []string) error {
	message := fmt.Sprintf(`{"m":"quote_remove_symbols","p":["%s","%s"]}`,
		session,
		strings.Join(symbols, `","`),
	)
	return sendWSMessage(c.ws, message, "quote remove symbols message")
}
