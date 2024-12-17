package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func generateSession(prefix string) string {
	b := make([]byte, 12)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return prefix + string(b)
}

func (rtd *RealTimeData) prependHeader(message string) string {
	return fmt.Sprintf("~m~%d~m~%s", len(message), message)
}

func (rtd *RealTimeData) constructMessage(function string, params []interface{}) (string, error) {
	msg := struct {
		M string        `json:"m"`
		P []interface{} `json:"p"`
	}{
		M: function,
		P: params,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (rtd *RealTimeData) createMessage(function string, params []interface{}) (string, error) {
	msg, err := rtd.constructMessage(function, params)
	if err != nil {
		return "", err
	}
	return rtd.prependHeader(msg), nil
}

func (rtd *RealTimeData) sendMessage(function string, params []interface{}) error {
	msg, err := rtd.createMessage(function, params)
	if err != nil {
		return err
	}

	return rtd.ws.WriteMessage(websocket.TextMessage, []byte(msg))
}
