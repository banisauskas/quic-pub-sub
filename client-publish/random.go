package main

import (
	"math/rand"
	"time"
)

var randInit = false
var randLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// Generates random message; length = [1; 20]
func randomMessage() string {
	if !randInit {
		rand.Seed(time.Now().UnixNano())
		randInit = true
	}

	length := rand.Intn(20) + 1
	message := make([]rune, length)

	for i := range message {
		message[i] = randLetters[rand.Intn(len(randLetters))]
	}

	return string(message)
}
