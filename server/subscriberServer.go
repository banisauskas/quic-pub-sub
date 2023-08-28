package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/quic-go/quic-go"
)

const subAddress = "localhost:2222"
const subPingByte = 0

var subscribers = make(map[string]*subCon)

type subCon struct {
	connection quic.Connection
	stream     quic.Stream
	lastPing   int64
}

func subscriberServer(tlsConfig *tls.Config) {
	listener, err := quic.ListenAddr(subAddress, tlsConfig, nil)
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

		subscriber := &subCon{
			connection,
			stream,
			time.Now().Unix(), // first ping time, because 'AcceptStream' was trigerred by first ping
		}

		subscribers[connectionID(connection)] = subscriber
		fmt.Println("SUBSCRIBERS:", len(subscribers))

		go handleSubscriber(subscriber)
	}
}

// If error occurs and this handler returns, last ping time won't be updated,
// subscriber becomes timed-out and later automatically discared.
func handleSubscriber(subscriber *subCon) {
	buf1 := make([]byte, 1)

	for {
		// blocks if nothing to read (n is always 1)
		n, err := subscriber.stream.Read(buf1)

		// error if blocked for >1 min (e.g. subscriber disconnected)
		if err != nil {
			return
		}

		// expecting 1 ping byte=0
		if n != 1 || buf1[0] != subPingByte {
			return
		}

		subscriber.lastPing = time.Now().Unix()
	}
}
