package main


//starts a connection with huobi
func initHuobi( out chan ListenOut, stop chan bool){
	c := NewClient(HuobiHandler{})
	c.Start("wss://api.huobi.pro/ws", "market.ethusdt.trade.detail", "id1", out, stop)
	go Printer(out, stop)
}


//start connection with binance
func initBinance(out chan ListenOut, stop chan bool) {
	c := NewClient(BinanceHandler{},)
	c.Start("wss://stream.binance.com:9443/ws/", "bnbbtc@trade","id1", out, stop)
	go Printer(out, stop)
}

func initBitfinex(out chan ListenOut, stop chan bool) {
	c := NewClient(BitfinexHandler{})
	c.Start("wss://api.bitfinex.com/ws/2", "tBTCUSD", "id1", out, stop)
	go Printer(out, stop)
}