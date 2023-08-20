package main

import (
	"math/rand"
	"time"
)

var initialized = false
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomMessage() string {
	if !initialized {
		rand.Seed(time.Now().UnixNano())
	}

	var length = rand.Intn(20) + 1 // 1-20
    var message = make([]rune, length)

    for i := range message {
        message[i] = letterRunes[rand.Intn(len(letterRunes))]
    }

    return string(message)
}