package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)


func NewClient(url string, handler GenericresHandler) Client {
	c := Client{url, handler}
	return c
}

func (c *Client) EstablishConn(url string, out chan ListenOut, stop chan bool) Ws {
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
	go c.listener(conn, out, stop)
	webSocket := Ws{conn, true}
	return webSocket
}

func (c *Client) listener( conn *websocket.Conn, out chan ListenOut, stop chan bool, ) {
	defer conn.Close()
	for {
		select {
		default:
			_, r, err := conn.NextReader()
			if (err != nil) {
				panic(err)
			}
			c.handler.handle(conn, r, out)
		case <-stop:
			return
		}
	}
}