package main

import (
	"io"
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

func prepSubmessage(subMessage string, id string) []byte {
	v := SubReqSend{subMessage, id, "false"}
	toSub, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return toSub
}

// convert reader to []byte
func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}


// Unzip a reader to string
func Unzip(reader io.Reader) string {
	r, err := gzip.NewReader(reader)
	if (err != nil) {
		panic(err)
	}
	defer r.Close()
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		panic(r)
	}
	return string(buf)
}
//Printer prints contents of a channel until a stop signal is given
func Printer(l chan ListenOut, stop chan bool){
	var i int = 0
	for {
		select {
		default:
			fmt.Println(i)
			fmt.Println(<-l)
			i++
		case <-stop:
			return
		}
	}
}