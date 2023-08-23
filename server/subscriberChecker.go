package main

import (
	"time"
)

// if there is at least 1 subscriber
var subsExist bool = false

func checkSubscribers() {
	for {
		var subCons = len(subscriberConnections)

		if subCons > 0 {
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

		time.Sleep(time.Second) // wait 1 sec.
	}
}

func notifyPublishers(subsExist bool) {
	var payload []byte

	if subsExist {
		payload = subscribersExist
	} else {
		payload = subscribersNotExist
	}

	for _, pubStream := range publisherStreams {
		var _, err = pubStream.Write(payload)

		if err != nil {
			panic(err)
		}
	}
}
