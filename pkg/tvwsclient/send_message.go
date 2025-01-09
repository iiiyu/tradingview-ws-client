package tvwsclient

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Add these fields as a constant since they're used in quote set fields
const defaultQuoteFields = "base-currency-logoid,ch,chp,currency-logoid,currency_code,currency_id," +
	"base_currency_id,current_session,description,exchange,format,fractional,is_tradable,language," +
	"local_description,listed_exchange,logoid,lp,lp_time,minmov,minmove2,original_name,pricescale," +
	"pro_name,short_name,type,typespecs,update_mode,volume,variable_tick_size,value_unit_id," +
	"unit_id,measure"

// sendWSMessage is a helper function that handles the common pattern of sending websocket messages
func sendWSMessage(ws *websocket.Conn, message string, operation string) error {
	wrappedMsg := wrappedMessage(message)
	slog.Debug("Send Message", "message", wrappedMsg)
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

// Chart Messages
func SendChartCreateSessionMessage(c *Client, session string) error {
	message := fmt.Sprintf(`{"m":"chart_create_session","p":["%s",""]}`, session)
	return sendWSMessage(c.ws, message, "chart create session message")
}

func SendSwitchTimezone(c *Client, session string) error {
	message := fmt.Sprintf(`{"m":"switch_timezone","p":["%s","Etc/UTC"]}`, session)
	return sendWSMessage(c.ws, message, "switch timezone message")
}

func SendResolveSymbol(c *Client, session string, symbol string) error {
	message := fmt.Sprintf(`{"m":"resolve_symbol","p":["%s","sds_sym_1","={\"adjustment\":\"splits\",\"session\":\"regular\",\"symbol\":\"%s\"}"]}`, session, symbol)
	return sendWSMessage(c.ws, message, "resolve symbol")
}

func SendCreateSeries(c *Client, session string, interval string) error {
	message := fmt.Sprintf(`{"m":"create_series","p":["%s","sds_1","s1","sds_sym_1","%s",10,""]}`, session, interval)
	return sendWSMessage(c.ws, message, "chart create session message")
}

// Quote Messages
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
