package main

import (
	"fmt"
	"time"
)

// If there is at least 1 subscriber (initially 0)
var subsExist bool = false
var subsExistPayload = []byte{1}
var subsNotExistPayload = []byte{0}

func discardSubscribers() {
	for {
		// Discard timed-out subscribers

		nowTime := time.Now().Unix()

		for conID, sub := range subscribers {
			if nowTime-sub.lastPing > 5 { // should be max 2-3 sec; consider timeout if over 5 sec.
				sub.stream.Close()
				delete(subscribers, conID)
				fmt.Println("SUBSCRIBERS:", len(subscribers))
			}
		}

		// Notify publishers

		if len(subscribers) > 0 {
			if !subsExist {
				subsExist = true
				notifyPublishers(true)
			}
		} else {
			if subsExist {
				subsExist = false
				notifyPublishers(false)
			}
		}

		// Wait 1 sec.

		time.Sleep(time.Second)
	}
}

// Notifies ALL publisher is there is at least 1 subscriber
func notifyPublishers(subsExist bool) {
	var payload []byte

	if subsExist {
		payload = subsExistPayload
	} else {
		payload = subsNotExistPayload
	}

	for _, pub := range publishers {
		_, err := pub.stream.Write(payload)

		if err != nil {
			panic(err)
		}
	}
}

// Notifies ONE publisher is there is at least 1 subscriber
func notifyPublisher(publisher *pubCon, subsExist bool) {
	var payload []byte

	if subsExist {
		payload = subsExistPayload
	} else {
		payload = subsNotExistPayload
	}

	_, err := publisher.stream.Write(payload)

	if err != nil {
		panic(err)
	}
}
