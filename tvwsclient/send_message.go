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

const (
	OnlySymbol       = "only_symbol"
	LessParameters   = "less_parameters"
	MediumParameters = "medium_parameters"
	MoreParameters   = "more_parameters"
	MostParameters   = "most_parameters"
)

// sendWSMessage is a helper function that handles the common pattern of sending websocket messages
func sendWSMessage(ws *websocket.Conn, message string, operation string) error {
	if ws == nil {
		return fmt.Errorf("websocket connection is nil, cannot send %s", operation)
	}
	
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
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	message := fmt.Sprintf(`{"m":"set_auth_token","p":["%s"]}`, authToken)
	return sendWSMessage(ws, message, "set auth token message")
}

func SendSetLocalMessage(c *Client) error {
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	message := `{"m":"set_locale","p":["en","US"]}`
	return sendWSMessage(ws, message, "set local message")
}

// Chart Messages
func SendChartCreateSessionMessage(c *Client, session string) error {
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	message := fmt.Sprintf(`{"m":"chart_create_session","p":["%s",""]}`, session)
	return sendWSMessage(ws, message, "chart create session message")
}

func SendSwitchTimezoneMessage(c *Client, session string) error {
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	message := fmt.Sprintf(`{"m":"switch_timezone","p":["%s","Etc/UTC"]}`, session)
	return sendWSMessage(ws, message, "switch timezone message")
}

func SendResolveSymbolMessage(c *Client, session string, symbol string) error {
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	message := fmt.Sprintf(`{"m":"resolve_symbol","p":["%s","sds_sym_1","={\"adjustment\":\"splits\",\"session\":\"regular\",\"symbol\":\"%s\"}"]}`, session, symbol)
	return sendWSMessage(ws, message, "resolve symbol")
}

func SendCreateSeriesMessage(c *Client, session string, interval string, seriesNumber int64) error {
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	message := fmt.Sprintf(`{"m":"create_series","p":["%s","sds_1","s1","sds_sym_1","%s",%d,""]}`, session, interval, seriesNumber)
	return sendWSMessage(ws, message, "chart create session message")
}

// Chart Messages
func SendChartDeleteSessionMessage(c *Client, session string) error {
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	message := fmt.Sprintf(`{"m":"chart_delete_session","p":["%s",""]}`, session)
	return sendWSMessage(ws, message, "chart remove session message")
}

func SubscriptionChartSessionSymbol(client *Client, session string, symbol string, interval string, seriesNumber int64) error {
	if err := SendChartCreateSessionMessage(client, session); err != nil {
		slog.Error("failed to send chart create session message ", "error", err)
		return err
	}

	if err := SendSwitchTimezoneMessage(client, session); err != nil {
		slog.Error("failed to send switch timezone message ", "error", err)
		return err
	}

	if err := SendResolveSymbolMessage(client, session, symbol); err != nil {
		slog.Error("failed to send resolve symbol message ", "error", err)
		return err
	}

	if err := SendCreateSeriesMessage(client, session, interval, seriesNumber); err != nil {
		slog.Error("failed to send create series message ", "error", err)
		return err
	}
	return nil
}

// Quote Messages
func SendQuoteCreateSessionMessage(c *Client, session string) error {
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	message := fmt.Sprintf(`{"m":"quote_create_session","p":["%s"]}`, session)
	if err := sendWSMessage(ws, message, "quote create session message"); err != nil {
		return err
	}
	return nil
}

func SendQuoteSetFieldsMessage(c *Client, session string) error {
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	fields := strings.Split(defaultQuoteFields, ",")
	message := fmt.Sprintf(`{"m":"quote_set_fields","p":["%s",%s]}`,
		session,
		`"`+strings.Join(fields, `","`)+`"`,
	)
	return sendWSMessage(ws, message, "quote set fields message")
}

// func SendQuoteFastSymbolsMessage(c *Client, session string, symbols []string) error {
// 	// Transform symbols into required format with both regular and extended sessions
// 	formattedSymbols := make([]string, len(symbols)*2)
// 	for i, symbol := range symbols {
// 		// Regular session format
// 		formattedSymbols[i*2] = fmt.Sprintf(`={"adjustment":"dividends","backadjustment":"default","symbol":"%s"}`, symbol)
// 		// Extended session format
// 		formattedSymbols[i*2+1] = fmt.Sprintf(`={"adjustment":"dividends","backadjustment":"default","session":"extended","symbol":"%s"}`, symbol)
// 	}

// 	var formattedSymbolsString string
// 	for i, formattedString := range formattedSymbols {
// 		if i > 0 {
// 			formattedSymbolsString += `","`
// 		}
// 		formattedSymbolsString += formattedString
// 	}
// 	message := fmt.Sprintf(`{"m":"quote_fast_symbols","p":["%s","%s"]}`,
// 		session,
// 		formattedSymbolsString,
// 	)
// 	slog.Debug("Send Quote Fast Symbols Message", "message", message)
// 	return sendWSMessage(c.ws, message, "quote fast symbols message")
// }

func SendQuoteRemoveSymbolsMessage(c *Client, session string, symbols []string) error {
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	message := fmt.Sprintf(`{"m":"quote_remove_symbols","p":["%s","%s"]}`,
		session,
		strings.Join(symbols, `","`),
	)
	return sendWSMessage(ws, message, "quote remove symbols message")
}

func SendQuoteCompletedMessageAfterQuoteCompleted(c *Client, session string, receivedMessage string) error {
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	// Replace single backslash + quote with triple backslash + quote
	receivedMessage = strings.ReplaceAll(receivedMessage, `\`, `\\`)
	// Replace remaining quotes with escaped quotes
	receivedMessage = strings.ReplaceAll(receivedMessage, `"`, `\"`)
	message := fmt.Sprintf(`{"m":"quote_remove_symbols","p":["%s","%s"]}`,
		session,
		receivedMessage,
	)
	return sendWSMessage(ws, message, "remove quote message after quote completed message")
}

func SubscriptionQuoteSessionSymbol(client *Client, session string, symbol string) error {
	if err := SendQuoteCreateSessionMessage(client, session); err != nil {
		return err
	}

	// if err := SendQuoteSetFieldsMessage(client, session); err != nil {
	// 	return err
	// }

	// if err := SendQuoteAddSymbolsMessageWithType(client, session, symbol, MoreParameters); err != nil {
	// 	return err
	// }

	// if err := SendQuoteAddSymbolsMessageWithType(client, session, symbol, LessParameters); err != nil {
	// 	return err
	// }

	// if err := SendQuoteFastSymbolsMessageWithType(client, session, symbol, LessParameters); err != nil {
	// 	return err
	// }

	// if err := SendQuoteAddSymbolsMessageWithType(client, session, symbol, MostParameters); err != nil {
	// 	return err
	// }

	// if err := SendQuoteFastSymbolsMessageWithType(client, session, symbol, MoreParameters); err != nil {
	// 	return err
	// }

	// if err := SendQuoteAddSymbolsMessageWithType(client, session, symbol, MediumParameters); err != nil {
	// 	return err
	// }

	// if err := SendQuoteFastSymbolsMessageWithType(client, session, symbol, MostParameters); err != nil {
	// 	return err
	// }

	// if err := SendQuoteAddSymbolsMessageWithType(client, session, symbol, OnlySymbol); err != nil {
	// 	return err
	// }

	// if err := SendQuoteFastSymbolsMessageWithType(client, session, symbol, MediumParameters); err != nil {
	// 	return err
	// }

	// if err := SendQuoteFastSymbolsMessageWithType(client, session, symbol, OnlySymbol); err != nil {
	// 	return err
	// }

	if err := SendQuoteAddSymbolsMessageWithType(client, session, symbol, OnlySymbol); err != nil {
		return err
	}

	return nil
}
