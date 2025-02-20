package tvwsclient

import (
	"fmt"
)

func SendQuoteAddSymbolsMessageWithType(c *Client, session string, symbol string, symbolType string) error {
	params := getQuoteAddSymbolsMessageParams(symbol, symbolType)
	message := fmt.Sprintf(`{"m":"quote_add_symbols","p":["%s","%s"]}`,
		session,
		params,
	)
	return sendWSMessage(c.ws, message, "quote add symbols message")
}

func getQuoteAddSymbolsMessageParams(symbol string, symbolType string) string {
	switch symbolType {
	case OnlySymbol:
		return symbol
	case LessParameters:
		// "={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"symbol\":\"NASDAQ:NTLA\"}"
		return `={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"symbol\":\"` + symbol + `\"}`
	case MediumParameters:
		// "={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"symbol\":\"NASDAQ:NTLA\"}"
		return `={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"symbol\":\"` + symbol + `\"}`
	case MoreParameters:
		// "={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"session\":\"extended\",\"symbol\":\"NASDAQ:NTLA\"}"
		return `={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"session\":\"extended\",\"symbol\":\"` + symbol + `\"}`
	case MostParameters:
		// "={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"session\":\"extended\",\"symbol\":\"NASDAQ:NTLA\"}"
		return `={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"currency-id\":\"USD\",\"session\":\"extended\",\"symbol\":\"` + symbol + `\"}`
	}
	return symbol
}
