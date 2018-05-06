package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)


func NewClient(handler GenericresHandler) Client {
	c := Client{handler}
	return c
}

func (c *Client) EstablishConn(url string, subscription string,out chan ListenOut, stop chan bool ) Ws {
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
	webSocket := Ws{conn, true,url, subscription }
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