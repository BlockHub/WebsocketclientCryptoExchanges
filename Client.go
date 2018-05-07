package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
)

//represents a client, one client can have multiple connections
type Client struct{
	handler 	GenericresHandler

}

//NewClient returns a client object
func NewClient(handler GenericresHandler) Client {
	c := Client{handler}
	return c
}


//Establishconn makes a websocket connection to a url, then sets a listener and reconnector to that websocket
// and returns the websocket
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
	ws := Ws{conn, true,url, subscription }
	go c.listener(ws, out, stop)
	go c.reconnector(ws, out, stop, d)
	return ws
}

func (c *Client) listener( ws Ws, out chan ListenOut, stop chan bool, ) {
	defer ws.conn.Close()
	for {
		select {
		default:
			_, r, err := ws.conn.NextReader()
			if (err != nil) {
				if ce, ok := err.(*websocket.CloseError); ok {
					switch ce.Code {
					default :
						/*
						websocket.CloseNormalClosure,
						websocket.CloseGoingAway,
						websocket.CloseNoStatusReceived
						*/
						ws.connected = false
						fmt.Println("websocket closed")
					//default:
					//	panic(err)
					}
				}
				}
			c.handler.handle(ws, r, out)
		case <-stop:
			stop <- false
			return
		}
	}
}

func (c *Client) reconnector(ws Ws, out chan ListenOut, stop chan bool, d *websocket.Dialer ){
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
			case <-stop:
				stop <- false
				return
			}
		}
	}
}