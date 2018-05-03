package main

import (
	"time"
)



func main() {
	/*
	a := BinanceStreams{"bnbbtc", "depth"}
	b := BinanceStreams{"bnbbtc", "ticker"}
	c := []BinanceStreams{a, b}
	*/
	out := make(chan ListenOut)
	stop := make(chan bool)
	initHuobi(out, stop)
	//initBinance(out, stop, c)
	time.Sleep(50000*time.Second)
	stop <- true
}


