package WebsocketCryptoScraper

import "github.com/gorilla/websocket"

//Client can have multiple websocket connections to the same exchange
type Client struct {
	handler GenericresHandler
}

//NewClient returns a client object
func NewClient(handler GenericresHandler) Client {
	c := Client{handler}
	return c
}

//Start starts a websocket connection using the handler of that client
func (c *Client) Start(url string, subscription string, id string, out chan ListenOut, stop chan bool) {
	d := websocket.DefaultDialer
	ws := c.handler.EstablishConn(url, subscription, id, out, stop, d)
	go c.handler.listener(ws, out, stop, d)
}
