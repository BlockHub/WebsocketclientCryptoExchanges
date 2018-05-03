package main


//starts a connection with huobi
func initHuobi( out chan ListenOut, stop chan bool){
	c := NewClient("wss://api.huobi.pro/ws", HuobiHandler{})
	ws := c.EstablishConn(c.url, out, stop)
	h := HuobiHandler{}
	h.subscribe(ws.conn,"market.ethusdt.trade.detail", "id1" )
	h.subscribe(ws.conn, "market.ethbtc.trade.detail", "id2")
	go Printer(out, stop)

}

//start connection with binance
func initBinance(out chan ListenOut, stop chan bool, streams []BinanceStreams) {
	c := NewClient("wss://stream.binance.com:9443/ws/", BinanceHandler{})
	for i := 0; i < len(streams); i++ {
		subTo := c.url + streams[i].symbol + "@" + streams[i].streamtype
		c.EstablishConn(subTo, out, stop)
	}
}