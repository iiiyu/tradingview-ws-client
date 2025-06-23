module basic-example

go 1.24

require github.com/iiiyu/tradingview-ws-client v0.1.0

require github.com/gorilla/websocket v1.5.0 // indirect

replace github.com/iiiyu/tradingview-ws-client => ../..
