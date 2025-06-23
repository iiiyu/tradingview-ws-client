package tvwsclient

import "fmt"

// func SendQuoteFastSymbolsMessage(c *Client, session string, symbol string) error {
// 	// Transform symbols into required format with both regular and extended sessions
// 	// formattedSymbols := make([]string, len(symbols)*2)
// 	// for i, symbol := range symbols {
// 	// 	// Regular session format
// 	// 	formattedSymbols[i*2] = fmt.Sprintf(`={"adjustment":"dividends","backadjustment":"default","symbol":"%s"}`, symbol)
// 	// 	// Extended session format
// 	// 	formattedSymbols[i*2+1] = fmt.Sprintf(`={"adjustment":"dividends","backadjustment":"default","session":"extended","symbol":"%s"}`, symbol)
// 	// }

// 	// var formattedSymbolsString string
// 	// for i, formattedString := range formattedSymbols {
// 	// 	if i > 0 {
// 	// 		formattedSymbolsString += `","`
// 	// 	}
// 	// 	formattedSymbolsString += formattedString
// 	// }
// 	message := fmt.Sprintf(`{"m":"quote_fast_symbols","p":["%s","={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"symbol\":\"%s\"}","%s","={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"session\":\"extended\",\"symbol\":\"%s\"}"]}`,
// 		session,
// 		symbol,
// 		symbol,
// 		symbol,
// 	)
// 	// message := fmt.Sprintf(`{"m":"quote_fast_symbols","p":["%s","%s"]}`,
// 	// 	session,
// 	// 	symbol,
// 	// )
// 	return sendWSMessage(c.ws, message, "quote fast symbols message")
// }

func SendQuoteFastSymbolsMessageWithType(c *Client, session string, symbol string, symbolType string) error {
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	
	params := getQuoteFastSymbolsMessageParams(symbol, symbolType)
	message := fmt.Sprintf(`{"m":"quote_fast_symbols","p":["%s","%s"]}`,
		session,
		params,
	)
	return sendWSMessage(ws, message, "quote fast symbols message")
}

func getQuoteFastSymbolsMessageParams(symbol string, symbolType string) string {
	switch symbolType {
	case OnlySymbol:
		return symbol
	case LessParameters:
		// "={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"session\":\"extended\",\"symbol\":\"NASDAQ:NTLA\"}","={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"symbol\":\"NASDAQ:NTLA\"}"
		return `={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"session\":\"extended\",\"symbol\":\"` + symbol + `\"}` + `,{\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"symbol\":\"` + symbol + `\"}`
	case MediumParameters:
		// "={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"session\":\"extended\",\"symbol\":\"NASDAQ:NTLA\"}","={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"symbol\":\"NASDAQ:NTLA\"}","NASDAQ:NTLA"
		return `={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"session\":\"extended\",\"symbol\":\"` + symbol + `\"}` + `,{\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"symbol\":\"` + symbol + `\"}` + `,` + symbol
	case MoreParameters:
		// "={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"symbol\":\"NASDAQ:NTLA\"}","={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"session\":\"extended\",\"symbol\":\"NASDAQ:NTLA\"}","={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"session\":\"extended\",\"symbol\":\"NASDAQ:NTLA\"}"
		return `={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"symbol\":\"` + symbol + `\"}` + `,{\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"session\":\"extended\",\"symbol\":\"` + symbol + `\"}` + `,{\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"session\":\"extended\",\"symbol\":\"` + symbol + `\"}`
	case MostParameters:
		// "={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"session\":\"extended\",\"symbol\":\"NASDAQ:NTLA\"}","={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"session\":\"extended\",\"symbol\":\"NASDAQ:NTLA\"}","={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"symbol\":\"NASDAQ:NTLA\"}"
		return `={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"session\":\"extended\",\"symbol\":\"` + symbol + `\"}` + `,{\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"session\":\"extended\",\"symbol\":\"` + symbol + `\"}` + `,{\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"symbol\":\"` + symbol + `\"}`
	}
	return symbol
}
