package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"io"
	"time"
	"github.com/quic-go/quic-go"
)

const addr = "localhost:2222"
const separator = 0

func main() {
	var tlsConf = &tls.Config{
		InsecureSkipVerify: true,
		NextProtos: []string{"abc123"},
	}

	var connection, err1 = quic.DialAddr(context.Background(), addr, tlsConf, nil)
	if err1 != nil {
		panic(err1)
	}

	var stream, err2 = connection.OpenStreamSync(context.Background())
	if err2 != nil {
		panic(err2)
	}

	var _, err3 = stream.Write([]byte { 0 }) // to trigger server AcceptStream
	if err3 != nil {
		panic(err3)
	}

	var buf1 = make([]byte, 1)
	var message = make([]byte, 0, 10)

	for {
		for {
			var n, err = stream.Read(buf1) // non-blocking; n = 0 or 1 

			for n == 0 && err == nil {
				time.Sleep(time.Second)
				n, err = stream.Read(buf1)
			}

			if n == 1 {
				if buf1[0] == 0 { // char=0 as separator
					processMessage(message)
					message = message[:0] // clear
				} else {
					message = append(message, buf1[0])
				}
			}

			if err == io.EOF {
				panic("Server not responding")
			}
		}
	}
}

func processMessage(message []byte) {
	if len(message) > 0 {
		var response = strings.Split(string(message), "#") // # is separator
		processMessage2(response[0], response[1])
	}
}

func processMessage2(publisherID string, message string) {
	fmt.Printf("Received from publisher %v: %v\n", publisherID, message)
}