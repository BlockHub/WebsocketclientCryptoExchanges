package main

import (
	"io"
	"github.com/gorilla/websocket"
	"regexp"
	"strconv"
	"encoding/json"
)

type GenericresHandler interface {
	handle(ws Ws ,reader io.Reader, out chan ListenOut)
}


type HuobiHandler struct {}


//handle deals with messages from huobi
func (h HuobiHandler) handle(ws Ws, reader io.Reader, out chan ListenOut)  {
	messageIn := Unzip(reader)
	if matched, _ := regexp.MatchString("ping*", messageIn); matched {
		message, err := strconv.ParseInt(messageIn[8:len(messageIn)-1], 10, 64)
		if err != nil {
			panic(err)
		}
		v := PongData{ message}
		messageOut, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		ws.conn.WriteMessage(2, messageOut)
		out <- ListenOut{2, messageIn}
	} else {
		//TODO replace messagetype
		out <- ListenOut{2, messageIn}
	}

}

//subscribe sends a subscription message to huobi
func (h HuobiHandler) subscribe(ws Ws, subMessage string, id string){
	err := ws.conn.WriteMessage(websocket.TextMessage, prepSubmessage(subMessage, id))
	if err != nil {
		panic(err)
	}
}


type BinanceHandler struct {}

//handle messages from binance
func (b BinanceHandler) handle (ws Ws, reader io.Reader, out chan ListenOut) {
	out <- ListenOut{2, (string(streamToByte(reader)))}
}