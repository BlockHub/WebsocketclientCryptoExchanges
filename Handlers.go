package main

import (
	"io"
	"github.com/gorilla/websocket"
	"regexp"
	"strconv"
	"encoding/json"
	"fmt"
	"net/http"
)

type GenericresHandler interface {
	handle(ws Ws ,reader io.Reader, out chan ListenOut)
	listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer)
	EstablishConn(url string, subscription string,id string, out chan ListenOut, stop chan bool, d *websocket.Dialer) Ws
	reconnector(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer)
	heartBeat(ws Ws)
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
	} else if matched, _ := regexp.MatchString("pong*", messageIn); !matched{
		//TODO replace messagetype
		out <- ListenOut{2, messageIn}
	}

}

func (h HuobiHandler) listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer){
	defer ws.conn.Close()
	for {
		select {
		default:
			_, r, err := ws.conn.NextReader()
			h.heartBeat(ws)
			if (err != nil) {
				/*if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
					h.reconnector(ws, out, stop, d)
					ws.connected=false
					return
				} else {*/
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
				go h.listener(ws, out, stop, d)
			case <-stop:
				stop <- false
				return
			}
		}
	}
}

func (h HuobiHandler) heartBeat(ws Ws) {
	ping, err := json.Marshal(PingData{18212558000})
	if err != nil {
		panic(err)
	}
	ws.conn.WriteMessage(websocket.TextMessage,ping)
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

func (b BinanceHandler) listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer){
	defer ws.conn.Close()
	for {
		_, r, err := ws.conn.NextReader()
		if err != nil {
			panic(err)
		}
		b.handle(ws, r, out)
	}

}

func (b BinanceHandler) EstablishConn(url string, subscription string, id string, out chan ListenOut, stop chan bool, d *websocket.Dialer) Ws{
	req := http.Header{}
	url = url  + subscription
	conn, res, err := d.Dial(url, req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if (err != nil){
		panic(err)
	}
	ws := Ws{conn, true,url, subscription,  id}
	return ws
}

func (b BinanceHandler) reconnector(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer){
	if (!ws.connected) {
		b.EstablishConn(ws.url, ws.subscription, ws.id, out, stop, d)
	}
}

func (b BinanceHandler) heartBeat(ws Ws){

}