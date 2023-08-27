package main

import (
	"fmt"
	"strings"

	"github.com/quic-go/quic-go"
)

const separatorByte = 0

var receivedMessage = make([]byte, 0, 1)

func reader(stream quic.Stream) {
	// Possible 1 message won't fit into 1 'buf10'.
	// Also possible several messages fit into 1 'buf10'.
	buf10 := make([]byte, 10)

	for {
		// 'Read' is blocking, waits until there is at least 1 byte to return.
		// Except when error occurs, then returns immediatelly with any number of bytes.
		// Conclusion: blocks until n>0 or err!=nil.
		// Might return both n>0 and err!=nil, therefore must read bytes before error.
		n, err := stream.Read(buf10)

		if n > 0 {
			appendReceived(buf10, n)
		}

		if err != nil { // not always 'io.EOF'
			panic("Disconnected")
		}
	}
}

func appendReceived(buf []byte, n int) {
	for i := 0; i < n; i++ {
		if buf[i] == separatorByte {
			// Several consucutive separators are allowed,
			// then messages between them are 0 bytes long.
			if len(receivedMessage) > 0 {
				processMessage(string(receivedMessage))
				receivedMessage = receivedMessage[:0] // clear
			}
		} else {
			receivedMessage = append(receivedMessage, buf[i])
		}
	}
}

func processMessage(message string) {
	parts := strings.Split(message, "#") // using # to separate message parts
	fmt.Printf("From publisher #%v: %v\n", parts[0], parts[1])
}
