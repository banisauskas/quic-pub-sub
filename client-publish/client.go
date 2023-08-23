package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"time"

	"github.com/quic-go/quic-go"
)

const addr = "localhost:1111"

var separator = []byte{0}
var subscribersExist = false

func main() {
	var tlsConf = &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"abc123"},
	}

	var connection, err1 = quic.DialAddr(context.Background(), addr, tlsConf, nil)
	if err1 != nil {
		panic(err1)
	}

	var stream, err2 = connection.OpenStreamSync(context.Background())
	if err2 != nil {
		panic(err2)
	}

	var _, err3 = stream.Write(separator) // to trigger server AcceptStream
	if err3 != nil {
		panic(err3)
	}

	go trackSubscribers(stream)
	sendingLoop(stream)
}

func trackSubscribers(stream quic.Stream) {
	var buf1 = make([]byte, 1)

	for {
		var n, err = stream.Read(buf1) // non-blocking; n = 0 or 1

		for n == 0 && err == nil {
			time.Sleep(time.Second)
			n, err = stream.Read(buf1)
		}

		if n == 1 {
			if buf1[0] == 0 {
				subscribersExist = false
				fmt.Println("Stop sending")
			} else if buf1[0] == 1 {
				subscribersExist = true
				fmt.Println("Start sending")
			} else {
				panic("Must be 1 or 0")
			}
		}

		if err == io.EOF {
			panic("Server not responding")
		}
	}
}

func sendingLoop(stream quic.Stream) {
	for {
		for !subscribersExist {
			time.Sleep(time.Second)
		}

		var msg = randomMessage()
		fmt.Println("Send:", msg)
		var _, err1 = stream.Write([]byte(msg))

		if err1 != nil {
			panic(err1)
		}

		var _, err2 = stream.Write(separator)

		if err2 != nil {
			panic(err2)
		}

		time.Sleep(time.Second)
	}
}
