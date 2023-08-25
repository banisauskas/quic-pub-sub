package main

import (
	"context"
	"crypto/tls"

	"github.com/quic-go/quic-go"
)

const addr = "localhost:2222"

func main() {
	stream := createClient()

	go reader(stream)
	writer(stream)
}

func createClient() quic.Stream {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"abc123"},
	}

	connection, err := quic.DialAddr(context.Background(), addr, tlsConf, nil)
	if err != nil {
		panic(err)
	}

	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		panic(err)
	}

	return stream
}
