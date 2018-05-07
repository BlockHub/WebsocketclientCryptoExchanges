package main

import (
	"io"
	"github.com/gorilla/websocket"
	"regexp"
	"strconv"
	"encoding/json"
	"bytes"
	"fmt"
	"net/http"
)

type GenericresHandler interface {
	handle(ws Ws ,reader io.Reader, out chan ListenOut)
	listener(ws Ws, out chan ListenOut, stop chan bool)
	EstablishConn(url string, subscription string,id string, out chan ListenOut, stop chan bool, d *websocket.Dialer) Ws
	reconnector(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer)
}


type HuobiHandler struct {}

func (h HuobiHandler) EstablishConn(url string, subscription string, id string,
									out chan ListenOut, stop chan bool, d *websocket.Dialer ) Ws {
	req := http.Header{}
	conn, res, err := d.Dial(url, req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if (err != nil){
		panic(err)
	}
	ws := Ws{conn, true,url, subscription,  id}
	h.subscribe(ws, ws.subscription, ws.id)
	go h.listener(ws, out, stop)
	go h.reconnector(ws, out, stop, d)
	return ws
}

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

func (h HuobiHandler) listener(ws Ws, out chan ListenOut, stop chan bool){
	defer ws.conn.Close()
	for {
		select {
		default:
			_, byte, err := ws.conn.ReadMessage()
			r := bytes.NewReader(byte)
			if (err != nil) {
					panic(err)
			} else {
				h.handle(ws, r, out)
			}
		case <-stop:
			stop <- false
			return
		}
	}
}

func (h HuobiHandler)reconnector(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer){
	for {
		if ws.connected == false {
			select {
			default:
				fmt.Println("reconnecting")
				req := http.Header{}
				conn, res, err := d.Dial(ws.url, req)
				if err != nil {
					panic(err)
				}
				defer res.Body.Close()
				if (err != nil) {
					panic(err)
				}
				ws.conn = conn
				ws.connected = true
				go h.listener(ws, out, stop)
			case <-stop:
				stop <- false
				return
			}
		}
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