package main

import (
	"fmt"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"github.com/quic-go/quic-go"
)

func main() {
	var tlsConfig = generateTLSConfig()

	go subscriberServer(tlsConfig)
	publisherServer(tlsConfig)
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
		NextProtos: []string{"abc123"},
	}
}

// Generates unique connection ID
func connectionID(con quic.Connection) string {
	return fmt.Sprintf("R%vL%v", con.RemoteAddr().String(), con.LocalAddr().String())
}