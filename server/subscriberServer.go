package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/quic-go/quic-go"
)

const subscriberAddr = "localhost:2222"

var subscribers = make(map[string]*subscriber)

type subscriber struct {
	connection quic.Connection
	stream     quic.Stream
	lastPing   int64
}

func subscriberServer(tlsConfig *tls.Config) {
	var listener, err1 = quic.ListenAddr(subscriberAddr, tlsConfig, nil)
	if err1 != nil {
		panic(err1)
	}

	for {
		var connection, err2 = listener.Accept(context.Background())
		if err2 != nil {
			panic(err2)
		}

		var stream, err3 = connection.AcceptStream(context.Background())
		if err3 != nil {
			panic(err3)
		}

		var subscriber = &subscriber{
			connection,
			stream,
			time.Now().Unix(), // valid 1st ping time, because AcceptStream was trigerred by 1st ping
		}

		subscribers[connectionID(connection)] = subscriber
		fmt.Println("Subscribers:", len(subscribers))
		go handleSubscriber(subscriber)
	}
}

func handleSubscriber(subscriber *subscriber) {
	var buf1 = make([]byte, 1)

	for {
		// blocks if nothing to read (n is always 1)
		var n, err = subscriber.stream.Read(buf1)

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
