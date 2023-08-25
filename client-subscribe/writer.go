package main

import (
	"time"

	"github.com/quic-go/quic-go"
)

// Ping every 2 sec with value 77.
// Must send at least 1 byte to trigger server AcceptStream.
func writer(stream quic.Stream) {
	ping := []byte{77}

	for {
		_, err := stream.Write(ping)
		if err != nil {
			panic(err)
		}

		time.Sleep(2 * time.Second)
	}
}
