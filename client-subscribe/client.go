package main

import (
	"context"
	"crypto/tls"

	"github.com/quic-go/quic-go"
)

const addr = "localhost:2222"

func main() {
	var stream = createClient()
	go reader(stream)
	writer(stream)
}

func createClient() quic.Stream {
	var tlsConf = &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"abc123"},
	}

	var connection, err1 = quic.DialAddr(context.Background(), addr, tlsConf, nil)
	if err1 != nil {
		panic(err1)
	}

	var stream, err2 = connection.OpenStreamSync(context.Background())
	if err2 != nil {
		panic(err2)
	}

	return stream
}
