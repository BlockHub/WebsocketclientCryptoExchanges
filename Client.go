package main

import "github.com/gorilla/websocket"

//represents a client, one client can have multiple connections
type Client struct{
	handler 	GenericresHandler

}

//NewClient returns a client object
func NewClient(handler GenericresHandler) Client {
	c := Client{handler}
	return c
}

func (c *Client) Start(url string, subscription string, id string, out chan ListenOut, stop chan bool){
	d := websocket.DefaultDialer
	ws := c.handler.EstablishConn(url, subscription, id, out, stop, d)
	go c.handler.listener(ws, out, stop, d)
	go c.handler.reconnector(ws, out, stop, d)
}