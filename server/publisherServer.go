package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"time"

	"github.com/quic-go/quic-go"
)

const publisherAddr = "localhost:1111"
const separator = 0

var publisherConnections = make(map[string]quic.Connection)
var publisherStreams = make(map[string]quic.Stream)
var publisherIDs int = 0

func publisherServer(tlsConfig *tls.Config) {
	var listener, err1 = quic.ListenAddr(publisherAddr, tlsConfig, nil)
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

		var conID = connectionID(connection)
		publisherConnections[conID] = connection
		publisherStreams[conID] = stream

		publisherIDs++
		go handlePublisher(publisherIDs, conID, stream)
	}
}

func handlePublisher(publisherID int, connectionID string, stream quic.Stream) {
	if len(subscriberStreams) > 0 {
		var _, err1 = stream.Write(subscribersExist)

		if err1 != nil {
			panic(err1)
		}
	}

	var online = true
	var buf1 = make([]byte, 1)
	var message = make([]byte, 0, 10)

	for online {
		for {
			var n, err = stream.Read(buf1) // non-blocking; n = 0 or 1

			for n == 0 && err == nil {
				time.Sleep(time.Second)
				n, err = stream.Read(buf1)
			}

			if n == 1 {
				if buf1[0] == separator {
					processMessage(publisherID, message)
					message = message[:0] // clear
				} else {
					message = append(message, buf1[0])
				}
			}

			if err == io.EOF {
				processMessage(publisherID, message)
				message = message[:0] // clear
				online = false
				break
			}
		}
	}

	delete(publisherStreams, connectionID)
	fmt.Printf("Publisher %v quit\n", publisherID)
}

func processMessage(publisherID int, message []byte) {
	if len(message) > 0 {
		processMessage2(publisherID, string(message))
	}
}

func processMessage2(publisherID int, message string) {
	fmt.Printf("Forwarding from publisher %v (to %v subscribers): %v\n", publisherID, len(subscriberStreams), message)

	for _, subStream := range subscriberStreams {
		var _, err1 = subStream.Write([]byte(fmt.Sprintf("%v#%v\x00", publisherID, message))) // char=0 as separator

		if err1 != nil {
			panic(err1)
		}
	}
}
