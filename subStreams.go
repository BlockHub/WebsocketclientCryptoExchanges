package main

import "fmt"

type GenericStream interface {
	reconnect()
	init()
}

//TODO make sure handlers detect connection loss and reconnect using the correct subscription data
type BinanceStream struct {
	symbol 		string
	streamtype 	string
}
func (b BinanceStream) init(){}

func (b BinanceStream) reconnect(){
	fmt.Println("test good")
}

type HuobiStream struct {
	dataType	string
	topic 		string
	description	string
}

func (h HuobiStream) init(){}

func (h HuobiStream) reconnect(){}

