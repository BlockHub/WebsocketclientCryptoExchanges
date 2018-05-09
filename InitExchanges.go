package main
/*
This contains examples of how to initialize an exchange. Printer is only used to consume values fromt the channels
A seperate client should be used per exchange. Each client should use multiples Ws connections for each exchange
 */



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


//use subscription as the channel(trades etc.) and id as Symbol, Binance does not allow for user provided ID
func initBitfinex(out chan ListenOut, stop chan bool) {
	c := NewClient(BitfinexHandler{})
	c.Start("wss://api.bitfinex.com/ws/2", "trades", "tBTCUSD", out, stop)
	go Printer(out, stop)
}

//subscription should be in the form channel@symbol, HitBTC handler takes care of the rest
func initHitBTC(out chan ListenOut, stop chan bool) {
	c := NewClient(hitBtcHandler{})
	c.Start("wss://api.hitbtc.com/api/2/ws", "subscribeTicker@ETHBTC", "2", out, stop)
	go Printer(out, stop)
}