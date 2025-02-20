package tvwsclient

import "fmt"

func SendQuoteFastSymbolsMessage(c *Client, session string, symbol string) error {
	// Transform symbols into required format with both regular and extended sessions
	// formattedSymbols := make([]string, len(symbols)*2)
	// for i, symbol := range symbols {
	// 	// Regular session format
	// 	formattedSymbols[i*2] = fmt.Sprintf(`={"adjustment":"dividends","backadjustment":"default","symbol":"%s"}`, symbol)
	// 	// Extended session format
	// 	formattedSymbols[i*2+1] = fmt.Sprintf(`={"adjustment":"dividends","backadjustment":"default","session":"extended","symbol":"%s"}`, symbol)
	// }

	// var formattedSymbolsString string
	// for i, formattedString := range formattedSymbols {
	// 	if i > 0 {
	// 		formattedSymbolsString += `","`
	// 	}
	// 	formattedSymbolsString += formattedString
	// }
	message := fmt.Sprintf(`{"m":"quote_fast_symbols","p":["%s","={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"symbol\":\"%s\"}","%s","={\"adjustment\":\"dividends\",\"backadjustment\":\"default\",\"session\":\"extended\",\"symbol\":\"%s\"}"]}`,
		session,
		symbol,
		symbol,
		symbol,
	)
	// message := fmt.Sprintf(`{"m":"quote_fast_symbols","p":["%s","%s"]}`,
	// 	session,
	// 	symbol,
	// )
	return sendWSMessage(c.ws, message, "quote fast symbols message")
}
