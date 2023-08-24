package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/quic-go/quic-go"
)

func main() {
	var tlsConfig = generateTLSConfig()

	go subscriberServer(tlsConfig)
	go publisherServer(tlsConfig)
	checkSubscribers()
}

func generateTLSConfig() *tls.Config {
	var key, err1 = rsa.GenerateKey(rand.Reader, 1024)

	if err1 != nil {
		panic(err1)
	}

	var template = x509.Certificate{SerialNumber: big.NewInt(1)}
	var certDER, err2 = x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)

	if err2 != nil {
		panic(err2)
	}

	var keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	var certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	var tlsCert, err3 = tls.X509KeyPair(certPEM, keyPEM)

	if err3 != nil {
		panic(err3)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"abc123"},
	}
}

// Generates unique connection ID
func connectionID(connection quic.Connection) string {
	return fmt.Sprintf("R%vL%v", connection.RemoteAddr().String(), connection.LocalAddr().String())
}
