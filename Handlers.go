package main

import (
	"io"
	"github.com/gorilla/websocket"
	"regexp"
	"strconv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type GenericresHandler interface {
	handle(ws Ws ,reader io.Reader, out chan ListenOut)
	listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer)
	EstablishConn(url string, subscription string,id string, out chan ListenOut, stop chan bool, d *websocket.Dialer) Ws
	reconnecter(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer)
}


type HuobiHandler struct {}

func (h HuobiHandler) EstablishConn(url string, subscription string, id string,
									out chan ListenOut, stop chan bool, d *websocket.Dialer ) Ws {
	req := http.Header{}
	conn, _, err := d.Dial(url, req)
	if err != nil {
		panic(err)
	}
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
		v := PongHuobi{ message}
		messageOut, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		ws.conn.WriteMessage(2, messageOut)
		out <- ListenOut{"Huobi", messageIn}
	} else if matched, _ := regexp.MatchString("pong*", messageIn); !matched{
		out <- ListenOut{"Huobi", messageIn}
	}

}

func (h HuobiHandler) listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer){
	defer ws.conn.Close()
	for {
		select {
		default:
			_, r, err := ws.conn.NextReader()
			ws.conn.SetReadDeadline(time.Now().Add(5*time.Second))
			h.heartBeat(ws)
			if (err != nil) {
				if strings.Contains(err.Error(), "i/o timeout") || strings.Contains(err.Error(), "unexpected EOF") {
					fmt.Println("set ws connected to false")
					h.reconnecter(ws, out, stop, d)
					return
				} else {
					panic(err)
				}
			} else {
				h.handle(ws, r, out)
			}
		case <-stop:
			stop <- false
			return
		}
	}
}

func (h HuobiHandler)reconnecter(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer){
	req := http.Header{}
	conn, _, err := d.Dial(ws.url, req)
	if err != nil {
		panic(err)
	}
	ws = Ws{conn, true, ws.url, ws.subscription, ws.id}
	h.subscribe(ws, ws.subscription, ws.id)
	go h.listener(ws, out, stop, d)
	out <- ListenOut{"Huobi", "reconnected"}
}

func (h HuobiHandler) heartBeat(ws Ws) {
	ping, err := json.Marshal(PingHuobi{18212558000})
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
	out <- ListenOut{"Binance", (string(streamToByte(reader)))}
}

func (b BinanceHandler) listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer){
	defer ws.conn.Close()
	for {
		_, r, err := ws.conn.NextReader()
		ws.conn.SetReadDeadline(time.Now().Add(30*time.Second))
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				fmt.Println("set ws connected to false")
				b.reconnecter(ws, out, stop, d)
				return
			} else {
				panic(err)
			}
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

func (b BinanceHandler) reconnecter(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer){
	req := http.Header{}
	conn, _, err := d.Dial(ws.url, req)
	if err != nil {
		panic(err)
	}
	ws = Ws{conn, true, ws.url, ws.subscription, ws.id}
	go b.listener(ws, out, stop, d)
	out <- ListenOut{"Binance", "reconnected"}
}

func (b BinanceHandler) heartBeat(ws Ws){
	ping, err := json.Marshal(PingBinance{})
	if err != nil {
		panic(err)
	}
	ws.conn.WriteMessage(9, ping)
}

type BitfinexHandler struct {

}

func (bf BitfinexHandler) handle(ws Ws ,reader io.Reader, out chan ListenOut){
	//TODO filter out info events and corresponding codes
	out <- ListenOut{"Bitfinex", (string(streamToByte(reader)))}
}

func (bf BitfinexHandler) listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer){
	defer ws.conn.Close()
	bf.subscribe(ws)
	for {
		_, r, err := ws.conn.NextReader()
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				fmt.Println("set ws connected to false")
				bf.reconnecter(ws, out, stop, d)
				return
			} else {
				panic(err)
			}
		}
		bf.handle(ws, r, out)
	}
}

func (bf BitfinexHandler) EstablishConn(url string, subscription string,id string, out chan ListenOut, stop chan bool, d *websocket.Dialer) Ws {
	req := http.Header{}
	conn, res, err := d.Dial(url, req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	ws := Ws{conn,true, url, subscription, id}
	return ws
}

//subscribe uses ws.id as symbol and ws.subscription as the channel
func (bf BitfinexHandler) subscribe(ws Ws) {
	subMessage, err := json.Marshal(BitFinexSub{"subscribe", ws.subscription, ws.id})
	if err != nil {
		panic(err)
	}
	ws.conn.WriteMessage(websocket.TextMessage, subMessage)
}


func (bf BitfinexHandler) reconnecter(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer) {
	req := http.Header{}
	conn, _, err := d.Dial(ws.url, req)
	if err != nil {
		panic(err)
	}
	ws = Ws{conn, true, ws.url, ws.subscription, ws.id}
	go bf.listener(ws, out, stop, d)
	out <- ListenOut{"Bitfinex", "reconnected"}
}

func (bf BitfinexHandler) pinger(ws Ws){
	ping, err := json.Marshal(BitFinexPing{"ping", 1234})
	if err != nil {
		panic(err)
	}
	ws.conn.WriteMessage(websocket.TextMessage,ping)
}