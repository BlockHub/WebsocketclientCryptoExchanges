package WebsocketCryptoScraper

//ListenOut contains data from exchanges and the corresponding exchange
type ListenOut struct {
	exchange string
	message  interface{}
}

//Huobi models

//PingHuobi is a ping message to Huobi
type PingHuobi struct {
	Ping int64 `json:"ping"`
}

//PongHuobi a pong message to Huobi
type PongHuobi struct {
	Pong int64 `json:"pong"`
}

//HuobiSubscription is a subscription a Huobi channel
type HuobiSubscription struct {
	Sub   string `json:"sub"`
	ID    string `json:"id"`
	Unsub string `json:"unsub"`
}

//Binance models

//PingBinance is a ping to Binance
type PingBinance struct {
	Ping int `json:"ping"`
}

//Bitfinex models

//BitFinexSub is a subscription message to BitFinex
type BitFinexSub struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Symbol  string `json:"symbol"`
}

//BitFinexPing is a ping message to BitFinex
type BitFinexPing struct {
	Event string `json:"event"`
	Cid   int    `json:"cid"`
}

//HitBTC models

//HitBtcSubscription is a subscription message to HitBtc
type HitBtcSubscription struct {
	Method string       `json:"method"`
	Params HitBtcParams `json:"params"`
	ID     int          `json:"id"`
}

//HitBtcParams is used for HitBtcSubscription
type HitBtcParams struct {
	Symbol string `json:"symbol"`
}
