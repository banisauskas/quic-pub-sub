package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/quic-go/quic-go"
)

const pubAddress = "localhost:1111"
const separator = 0

var globalIDs int = 0
var publishers = make(map[string]*pubCon)

var receivedMessage = make([]byte, 0, 100)

type pubCon struct {
	publisherID int
	connection  quic.Connection
	stream      quic.Stream
	lastPing    int64
}

func publisherServer(tlsConfig *tls.Config) {
	listener, err := quic.ListenAddr(pubAddress, tlsConfig, nil)
	if err != nil {
		panic(err)
	}

	for {
		connection, err := listener.Accept(context.Background())
		if err != nil {
			panic(err)
		}

		stream, err := connection.AcceptStream(context.Background())
		if err != nil {
			panic(err)
		}

		globalIDs++

		publisher := &pubCon{
			globalIDs,
			connection,
			stream,
			time.Now().Unix(), // first ping time, because 'AcceptStream' was trigerred by first ping
		}

		publishers[connectionID(connection)] = publisher
		fmt.Println("PUBLISHERS:", len(publishers))

		go handlePublisher(publisher)
	}
}

func handlePublisher(publisher *pubCon) {
	// Manually notify

	if len(subscribers) > 0 {
		notifyPublisher(publisher, true)
	}

	// Accept messages

	// Possible 1 message won't fit into 1 'buf10'.
	// Also possible several messages fit into 1 'buf10'.
	buf10 := make([]byte, 10)

	for {
		// 'Read' is blocking, waits until there is at least 1 byte to return.
		// Except when error occurs, then returns immediatelly with any number of bytes.
		// Conclusion: blocks until n>0 or err!=nil.
		// Might return both n>0 and err!=nil, therefore must read bytes before error.
		n, err := publisher.stream.Read(buf10)

		if n > 0 {
			appendReceived(buf10, n, publisher.publisherID)
			publisher.lastPing = time.Now().Unix()
		}

		if err != nil { // not always 'io.EOF'
			return
		}
	}
}

func appendReceived(buf []byte, n int, publisherID int) {
	for i := 0; i < n; i++ {
		if buf[i] == separator {
			// Several consucutive separators are allowed,
			// then messages between them are 0 bytes long.
			if len(receivedMessage) > 0 {
				processMessage(string(receivedMessage), publisherID)
				receivedMessage = receivedMessage[:0] // clear
			}
		} else {
			receivedMessage = append(receivedMessage, buf[i])
		}
	}
}

func processMessage(message string, publisherID int) {
	fmt.Printf("Forward from publisher #%v to %v subscribers: %v\n", publisherID, len(subscribers), message)

	// Message format "pubid#message\0"
	payload := []byte(fmt.Sprint(publisherID))
	payload = append(payload, '#')
	payload = append(payload, []byte(message)...)
	payload = append(payload, separator)

	for _, sub := range subscribers {
		_, err := sub.stream.Write(payload)

		if err != nil {
			panic(err)
		}
	}
}
