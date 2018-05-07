package main

import (
	"time"
)



func main() {
	out := make(chan ListenOut, 1000000)
	stop := make(chan bool)
	initHuobi(out, stop)
	time.Sleep(500000000*time.Second)
	stop <- true
}


