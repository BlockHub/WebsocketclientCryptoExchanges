package main

import "github.com/gorilla/websocket"

//Data to go to out channel
//all data from exchanges should be in ListenOut before it leaves te scraper
type ListenOut struct{
	mt 		int
	message interface{}
}

//ping message
type PingData struct {
	Ping int64 `json:"ping"`
}

//Pong message
type PongData struct {
	Pong int64 `json:"pong"`
}

//Subscription to a Huobi channel
type SubReqSend struct {
	Sub   string `json:"sub"`
	ID    string `json:"id"`
	Unsub string `json:"unsub"`
}

//represents a client, one client can have multiple connections
type Client struct{
	url 		string
	handler 	GenericresHandler

}

//represents a websocket
type Ws struct {
	conn 		*websocket.Conn
	connected 	bool
}

type BinanceStreams struct {
	symbol 		string
	streamtype 	string
}