package main

import (
	"github.com/gorilla/websocket"
)

//Ws represents a websocket connection
type Ws struct {
	conn 			*websocket.Conn
	connected 		bool
	url 			string
	subscription	string
	id				string
}
