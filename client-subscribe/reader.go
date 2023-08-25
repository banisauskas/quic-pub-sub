package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/quic-go/quic-go"
)

func reader(stream quic.Stream) {
	buf1 := make([]byte, 1)
	message := make([]byte, 0, 10)

	for {
		n, err := stream.Read(buf1) // non-blocking; n = 0 or 1

		for n == 0 && err == nil {
			time.Sleep(time.Second)
			n, err = stream.Read(buf1)
		}

		if n == 1 {
			if buf1[0] == 0 { // 0 is separator
				processMessage(message)
				message = message[:0] // clear
			} else {
				message = append(message, buf1[0])
			}
		}

		if err != nil { // must check after reading 1 byte
			panic("Disconnected")
		}
	}
}

func processMessage(message []byte) {
	if len(message) > 0 {
		parts := strings.Split(string(message), "#") // # is separator
		processMessage2(parts[0], parts[1])
	}
}

func processMessage2(publisherID string, message string) {
	fmt.Printf("Received from publisher %v: %v\n", publisherID, message)
}