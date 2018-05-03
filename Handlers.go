package main

import (
	"io"
	"github.com/gorilla/websocket"
	"fmt"
	"regexp"
	"strconv"
	"encoding/json"
)

type GenericresHandler interface {
	handle(conn *websocket.Conn ,reader io.Reader, out chan ListenOut)
}


type HuobiHandler struct {}


//handle messages from huobi
func (h HuobiHandler) handle(conn *websocket.Conn, reader io.Reader, out chan ListenOut)  {
	messageIn := Unzip(reader)
	fmt.Println(messageIn)
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
		conn.WriteMessage(2, messageOut)
		out <- ListenOut{2, messageIn}
	} else {
		//TODO replace messagetype
		out <- ListenOut{2, messageIn}
	}

}

// send a subscription message to huobi
func (h HuobiHandler) subscribe(conn *websocket.Conn, subMessage string, id string){
	err := conn.WriteMessage(websocket.TextMessage, prepSubmessage(subMessage, id))
	if err != nil {
		panic(err)
	}
}


type BinanceHandler struct {}

//handle messages from binance
func (b BinanceHandler) handle (conn *websocket.Conn, reader io.Reader, out chan ListenOut) {
	out <- ListenOut{2, (string(StreamToByte(reader)))}
}