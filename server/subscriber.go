package main

import (
	"context"
	"crypto/tls"

	"github.com/quic-go/quic-go"
)

const subscriberAddr = "localhost:2222"

var subscriberStreams = make(map[string]quic.Stream)
var subscribersExist = []byte{1}
var subscribersNotExist = []byte{0}

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

		var connectionID = connectionID(connection)
		subscriberStreams[connectionID] = stream

		if len(subscriberStreams) == 1 { // added first subscriber
			for _, pubStream := range publisherStreams {
				var _, err4 = pubStream.Write(subscribersExist)

				if err4 != nil {
					panic(err4)
				}
			}
		}
	}
}
