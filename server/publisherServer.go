package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/quic-go/quic-go"
)

const publisherAddr = "localhost:1111"
const separator = 0

var globalIDs int = 0
var publishers = make(map[string]*pubCon)

type pubCon struct {
	publisherID int
	connection  quic.Connection
	stream      quic.Stream
	lastPing    int64
}

func publisherServer(tlsConfig *tls.Config) {
	listener, err := quic.ListenAddr(publisherAddr, tlsConfig, nil)
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
			time.Now().Unix(), // valid 1st ping time, because AcceptStream was trigerred by 1st ping
		}

		publishers[connectionID(connection)] = publisher
		fmt.Println("Publishers:", len(publishers))

		go handlePublisher(publisher)
	}
}

func handlePublisher(publisher *pubCon) {
	if len(subscribers) > 0 {
		_, err := publisher.stream.Write(subsExistPayload)

		if err != nil {
			panic(err)
		}
	}

	buf1 := make([]byte, 1)
	message := make([]byte, 0, 10)

	for {
		n, err := publisher.stream.Read(buf1) // non-blocking; n = 0 or 1

		for n == 0 && err == nil {
			time.Sleep(time.Second)
			n, err = publisher.stream.Read(buf1)
		}

		if n == 1 {
			if buf1[0] == separator {
				processMessage(publisher.publisherID, message)
				publisher.lastPing = time.Now().Unix() // separator also serves as PING
				message = message[:0]                  // clear
			} else {
				message = append(message, buf1[0])
			}
		}

		if err != nil {
			return
		}
	}
}

func processMessage(publisherID int, message []byte) {
	if len(message) > 0 {
		processMessage2(publisherID, string(message))
	}
}

func processMessage2(publisherID int, message string) {
	fmt.Printf("Forwarding from publisher %v (to %v subscribers): %v\n", publisherID, len(subscribers), message)

	for _, sub := range subscribers {
		_, err := sub.stream.Write([]byte(fmt.Sprintf("%v#%v\x00", publisherID, message))) // char=0 as separator

		if err != nil {
			panic(err)
		}
	}
}
