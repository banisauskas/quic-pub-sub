package main

import (
	"time"

	"github.com/quic-go/quic-go"
)

const pingByte = 0
const pingTime = 2

// Ping every 2 sec with value byte=0.
// Must send at least 1 byte to trigger server 'AcceptStream'.
func writer(stream quic.Stream) {
	ping := []byte{pingByte}

	for {
		_, err := stream.Write(ping)
		if err != nil {
			panic(err)
		}

		time.Sleep(pingTime * time.Second)
	}
}
