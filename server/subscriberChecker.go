package main

import (
	"fmt"
	"time"
)

// If there is at least 1 subscriber.
var subsExist bool = false
var subsExistPayload = []byte{1}
var subsNotExistPayload = []byte{0}

func checkSubscribers() {
	for {
		// Discard timed-out subscribers

		var nowTime = time.Now().Unix()

		for conID, sub := range subscribers {
			if nowTime-sub.lastPing > 5 { // should be max 2-3 sec; timeout if over 5 sec.
				sub.stream.Close()
				delete(subscribers, conID)
				fmt.Println("Subscribers:", len(subscribers))
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

func notifyPublishers(subsExist bool) {
	var payload []byte

	if subsExist {
		payload = subsExistPayload
	} else {
		payload = subsNotExistPayload
	}

	for _, pubStream := range publisherStreams {
		var _, err = pubStream.Write(payload)

		if err != nil {
			panic(err)
		}
	}
}
