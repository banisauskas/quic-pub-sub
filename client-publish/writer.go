package main

import (
	"fmt"
	"time"

	"github.com/quic-go/quic-go"
)

// Separator between messages.
// Also used as PING signal.
const separator = 0

// Must immediately send at least 1 byte to trigger server 'AcceptStream'
func writer(stream quic.Stream) {
	for {
		sendMessages(stream)
		sendPings(stream)
	}
}

// Sends random messages every 1 sec,
// while there is at least 1 subscriber
func sendMessages(stream quic.Stream) {
	for subsExist {
		message := randomMessage()
		payload := append([]byte(message), separator)

		fmt.Println("Sending:", message)
		_, err := stream.Write(payload)

		if err != nil {
			panic(err)
		}

		time.Sleep(time.Second)
	}
}

// Sends PING every 2 sec,
// while there are no subscribers
func sendPings(stream quic.Stream) {
	for !subsExist {
		_, err := stream.Write([]byte{separator})

		if err != nil {
			panic(err)
		}

		time.Sleep(2 * time.Second)
	}
}
