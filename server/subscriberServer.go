package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/quic-go/quic-go"
)

const subscriberAddr = "localhost:2222"

var subscribers = make(map[string]*subCon)

type subCon struct {
	connection quic.Connection
	stream     quic.Stream
	lastPing   int64
}

func subscriberServer(tlsConfig *tls.Config) {
	listener, err := quic.ListenAddr(subscriberAddr, tlsConfig, nil)
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
			time.Now().Unix(), // valid 1st ping time, because AcceptStream was trigerred by 1st ping
		}

		subscribers[connectionID(connection)] = subscriber
		fmt.Println("Subscribers:", len(subscribers))

		go handleSubscriber(subscriber)
	}
}

func handleSubscriber(subscriber *subCon) {
	buf1 := make([]byte, 1)

	for {
		// blocks if nothing to read (n is always 1)
		n, err := subscriber.stream.Read(buf1)

		// error if blocked (subscriber disconnected) for >1 min.
		if err != nil {
			return
		}

		// expecting 1 ping byte=77
		if n != 1 || buf1[0] != 77 {
			return
		}

		subscriber.lastPing = time.Now().Unix()
	}
}
