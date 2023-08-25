package main

import (
	"fmt"
	"time"

	"github.com/quic-go/quic-go"
)

// If there is at least 1 subscriber.
var subscribersExist = false

func trackSubscribers(stream quic.Stream) {
	buf1 := make([]byte, 1)

	for {
		// Sometimes 'blocking': waits until can return one byte with n=1.
		// Sometimes 'non-blocking': when error occurs immediately returns n=0.
		n, err := stream.Read(buf1) // n = 0 or 1

		for n == 0 && err == nil {
			time.Sleep(time.Second)
			n, err = stream.Read(buf1)
		}

		if n == 1 {
			// Expecting 0 from server when subscribers don't exist, and 1 if exist.
			if buf1[0] == 0 {
				subscribersExist = false
				fmt.Println("Stop sending messages")
			} else if buf1[0] == 1 {
				subscribersExist = true
				fmt.Println("Start sending messages")
			} else {
				panic("Must be 1 or 0")
			}
		}

		if err != nil { // not always 'io.EOF'
			panic("Disconnected")
		}
	}
}
