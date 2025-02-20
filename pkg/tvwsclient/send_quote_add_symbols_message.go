package tvwsclient

import (
	"fmt"
)

func SendQuoteAddSymbolsMessageWithType(c *Client, session string, symbol string, symbolType string) error {
	params := GetQuoteAddSymbolsMessageParams(symbol, symbolType)
	message := fmt.Sprintf(`{"m":"quote_add_symbols","p":["%s","%s"]}`,
		session,
		params,
	)
	return sendWSMessage(c.ws, message, "quote add symbols message")
}

func GetQuoteAddSymbolsMessageParams(symbol string, symbolType string) string {
	switch symbolType {
	case OnlySymbol:
		return symbol
	case LessParameters:
		return `={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"symbol\":\"` + symbol + `\"}`
	case MoreParameters:
		return `={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"session\":\"extended\",\"symbol\":\"` + symbol + `\"}`
	case MostParameters:
		return `={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"session\":\"extended\",\"symbol\":\"` + symbol + `\"}`
	}
	return symbol
}
