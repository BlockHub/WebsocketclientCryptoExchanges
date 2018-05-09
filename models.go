package main

//Data to go to out channel
//all data from exchanges should be in ListenOut before it leaves te scraper
type ListenOut struct{
	mt 		int
	message interface{}
}


//Huobi models
//PingHuobi is a ping message to Huobi
type PingHuobi struct {
	Ping int64 `json:"ping"`
}

//PongHuobiis a pong message to Huobi
type PongHuobi struct {
	Pong int64 `json:"pong"`
}

//Subscription to a Huobi channel
type SubReqSend struct {
	Sub   string `json:"sub"`
	ID    string `json:"id"`
	Unsub string `json:"unsub"`
}



//Bitfinex models
type BitFinexSub struct {
	Event 	string 	`json:"event"`
	Channel string 	`json:"channel"`
	Symbol 	string	`json:"symbol"`

}

type BitFinexPing struct {
	Event 	string	`json:"event"`
	Cid		int		`json:"cid"`
}