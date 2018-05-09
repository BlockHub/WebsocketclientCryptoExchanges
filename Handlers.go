package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)
//GenericresHandler is an interface for all the functions each exchange needs to have at te least
type GenericresHandler interface {
	handle(ws Ws, reader io.Reader, out chan ListenOut)
	listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer)
	EstablishConn(url string, subscription string, id string, out chan ListenOut, stop chan bool, d *websocket.Dialer) Ws
	reconnecter(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer)
}

//HuobiHandler is the handler for Huobi
type HuobiHandler struct{}

//EstablishConn makes a connection to the Huobi endpoint and does not subscribe
func (h HuobiHandler) EstablishConn(url string, subscription string, id string,
	out chan ListenOut, stop chan bool, d *websocket.Dialer) Ws {
	req := http.Header{}
	conn, _, err := d.Dial(url, req)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	ws := Ws{conn, true, url, subscription, id}
	return ws
}

//handle deals with messages from huobi, responds correctly to pings and sends other messages to out
func (h HuobiHandler) handle(ws Ws, reader io.Reader, out chan ListenOut) {
	messageIn := Unzip(reader)
	if matched, _ := regexp.MatchString("ping*", messageIn); matched {
		message, err := strconv.ParseInt(messageIn[8:len(messageIn)-1], 10, 64)
		if err != nil {
			panic(err)
		}
		v := PongHuobi{message}
		messageOut, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		ws.conn.WriteMessage(2, messageOut)
		out <- ListenOut{"Huobi", messageIn}
	} else if matched, _ := regexp.MatchString("pong*", messageIn); !matched {
		out <- ListenOut{"Huobi", messageIn}
	}

}

func (h HuobiHandler) listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer) {
	defer ws.conn.Close()
	h.subscribe(ws, ws.subscription, ws.id)
	for {
		select {
		default:
			_, r, err := ws.conn.NextReader()
			ws.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			h.heartBeat(ws)
			if err != nil {
				if strings.Contains(err.Error(), "i/o timeout") || strings.Contains(err.Error(), "unexpected EOF") {
					fmt.Println("set ws connected to false")
					h.reconnecter(ws, out, stop, d)
					return
				}
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

func (h HuobiHandler) reconnecter(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer) {
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
	ws.conn.WriteMessage(websocket.TextMessage, ping)
}

//subscribe sends a subscription message to huobi
func (h HuobiHandler) subscribe(ws Ws, subMessage string, id string) {
	err := ws.conn.WriteMessage(websocket.TextMessage, prepSubmessage(subMessage, id))
	if err != nil {
		panic(err)
	}
}

//BinanceHandler is the handler for Huobi
type BinanceHandler struct{}

//handle messages from binance
func (b BinanceHandler) handle(ws Ws, reader io.Reader, out chan ListenOut) {
	out <- ListenOut{"Binance", (string(streamToByte(reader)))}
}

func (b BinanceHandler) listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer) {
	defer ws.conn.Close()
	for {
		_, r, err := ws.conn.NextReader()
		ws.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				fmt.Println("set ws connected to false")
				b.reconnecter(ws, out, stop, d)
				return
			}
			panic(err)
			}
		b.handle(ws, r, out)
	}

}

//EstablishConn makes a connection to the Binance endpoint and does subscribe
func (b BinanceHandler) EstablishConn(url string, subscription string, id string, out chan ListenOut, stop chan bool, d *websocket.Dialer) Ws {
	req := http.Header{}
	url = url + subscription
	conn, res, err := d.Dial(url, req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	ws := Ws{conn, true, url, subscription, id}
	return ws
}

func (b BinanceHandler) reconnecter(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer) {
	req := http.Header{}
	conn, _, err := d.Dial(ws.url, req)
	if err != nil {
		panic(err)
	}
	ws = Ws{conn, true, ws.url, ws.subscription, ws.id}
	go b.listener(ws, out, stop, d)
	out <- ListenOut{"Binance", "reconnected"}
}

func (b BinanceHandler) heartBeat(ws Ws) {
	ping, err := json.Marshal(PingBinance{})
	if err != nil {
		panic(err)
	}
	ws.conn.WriteMessage(9, ping)
}

//BitfinexHandler is the handler for Huobi
type BitfinexHandler struct {
}

func (bf BitfinexHandler) handle(ws Ws, reader io.Reader, out chan ListenOut) {
	//TODO filter out info events and corresponding codes
	out <- ListenOut{"Bitfinex", (string(streamToByte(reader)))}
}

func (bf BitfinexHandler) listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer) {
	defer ws.conn.Close()
	bf.subscribe(ws)
	for {
		_, r, err := ws.conn.NextReader()
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				fmt.Println("set ws connected to false")
				bf.reconnecter(ws, out, stop, d)
				return
			}
			panic(err)
		}
		bf.handle(ws, r, out)
	}
}

//EstablishConn makes a connection to the BitFinex endpoint and does not subscribe
func (bf BitfinexHandler) EstablishConn(url string, subscription string, id string, out chan ListenOut, stop chan bool, d *websocket.Dialer) Ws {
	req := http.Header{}
	conn, res, err := d.Dial(url, req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	ws := Ws{conn, true, url, subscription, id}
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

//currently unused, should be used if the user wants to check if the channel is open
func (bf BitfinexHandler) pinger(ws Ws) {
	ping, err := json.Marshal(BitFinexPing{"ping", 1234})
	if err != nil {
		panic(err)
	}
	ws.conn.WriteMessage(websocket.TextMessage, ping)
}

//HitBtcHandler is the handler for Huobi
type hitBtcHandler struct {
}

func (hi hitBtcHandler) handle(ws Ws, reader io.Reader, out chan ListenOut) {
	out <- ListenOut{"HitBTC", (string(streamToByte(reader)))}
}

//listener subcribes to a channel from ws when called, takes messages until it io times out and then reconnect or
//gives messages to handle
func (hi hitBtcHandler) listener(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer) {
	defer ws.conn.Close()
	hi.subscribe(ws)
	for {
		_, r, err := ws.conn.NextReader()
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				fmt.Println("set ws connected to false")
				hi.reconnecter(ws, out, stop, d)
				return
			}
			panic(err)
			}
		hi.handle(ws, r, out)
	}
}

//EstablishConn makes a connection to the HitBtc endpoint and does not subscribe
func (hi hitBtcHandler) EstablishConn(url string, subscription string, id string, out chan ListenOut, stop chan bool, d *websocket.Dialer) Ws {
	req := http.Header{}
	conn, res, err := d.Dial(url, req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	ws := Ws{conn, true, url, subscription, id}
	return ws
}

func (hi hitBtcHandler) reconnecter(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer) {
	req := http.Header{}
	conn, _, err := d.Dial(ws.url, req)
	if err != nil {
		panic(err)
	}
	ws = Ws{conn, true, ws.url, ws.subscription, ws.id}
	go hi.listener(ws, out, stop, d)
	out <- ListenOut{"HitBTC", "reconnected"}
}

func (hi hitBtcHandler) subscribe(ws Ws) {
	things := strings.Split(ws.subscription, "@")
	intId, err := strconv.Atoi(ws.id)
	if err != nil {
		panic(err)
	}
	params := HitBtcParams{things[1]}
	subscription, err := json.Marshal(HitBtcSubscription{things[0], params, intId})
	if err != nil {
		panic(err)
	}
	ws.conn.WriteMessage(websocket.TextMessage, subscription)

}
