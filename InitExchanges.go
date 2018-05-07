package main


//starts a connection with huobi
func initHuobi( out chan ListenOut, stop chan bool){
	c := NewClient(HuobiHandler{})
	c.start("wss://api.huobi.pro/ws", "market.ethusdt.trade.detail", "id1", out, stop)
	go Printer(out, stop)

}
/*
//start connection with binance
func initBinance(out chan ListenOut, stop chan bool) {
	c := NewClient("wss://stream.binance.com:9443/ws/", BinanceHandler{},)
	ws := c.EstablishConn(c.url, out, stop)

}
*/