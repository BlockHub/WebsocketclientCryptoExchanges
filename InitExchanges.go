package main


//starts a connection with huobi
func initHuobi( out chan ListenOut, stop chan bool){
	c := NewClient(HuobiHandler{})
	ws := c.EstablishConn("wss://api.huobi.pro/ws","market.ethusdt.trade.detail" , out, stop)
	h := HuobiHandler{}
	h.subscribe(ws,"market.ethusdt.trade.detail", "id1" )
	h.subscribe(ws, "market.ethusdt.trade.detail", "id2")
	go Printer(out, stop)

}
/*
//start connection with binance
func initBinance(out chan ListenOut, stop chan bool) {
	c := NewClient("wss://stream.binance.com:9443/ws/", BinanceHandler{},)
	ws := c.EstablishConn(c.url, out, stop)

}
*/