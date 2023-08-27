package main

import (
	"fmt"

	"github.com/quic-go/quic-go"
)

const signalSubsExist = 1
const signalSubsNotExist = 0

// If there is at least 1 subscriber.
var subsExist = false

func reader(stream quic.Stream) {
	buf1 := make([]byte, 1)

	for {
		// 'Read' is blocking, waits until there is at least 1 byte to return.
		// Except when error occurs, then returns immediatelly with any number of bytes.
		// Conclusion: blocks until n>0 or err!=nil.
		// Might return both n>0 and err!=nil, therefore must read bytes before error.
		n, err := stream.Read(buf1)

		if n == 1 { // n = 0 or 1
			processReceived(buf1[0])
		}

		if err != nil { // not always 'io.EOF'
			panic("Disconnected")
		}
	}
}

func processReceived(input byte) {
	if input == signalSubsExist {
		subsExist = true
		fmt.Println("START PUBLISHING")
	} else if input == signalSubsNotExist {
		subsExist = false
		fmt.Println("STOP PUBLISHING")
	} else {
		panic("Must be 1 or 0")
	}
}
