package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)


func NewClient(url string, handler GenericresHandler) Client {
	c := Client{url, handler}
	return c
}

func (c *Client) EstablishConn(url string, out chan ListenOut, stop chan bool, gs GenericStream) Ws {
	d := websocket.DefaultDialer
	req := http.Header{}
	conn, res, err := d.Dial(url, req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if (err != nil){
		panic(err)
	}
	webSocket := Ws{conn, true, gs}
	go c.listener(webSocket, out, stop)
	return webSocket
}

func (c *Client) listener( ws Ws, out chan ListenOut, stop chan bool, ) {
	defer ws.conn.Close()
	for {
		select {
		default:
			_, r, err := ws.conn.NextReader()
			if (err != nil) {
				panic(err)
			}
			c.handler.handle(ws, r, out)
		case <-stop:
			return
		}
	}
}