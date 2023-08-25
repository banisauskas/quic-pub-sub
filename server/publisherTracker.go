package main

import (
	"fmt"
	"time"
)

func trackPublishers() {
	for {
		// Discard timed-out publishers

		nowTime := time.Now().Unix()

		for conID, pub := range publishers {
			if nowTime-pub.lastPing > 5 { // should be max 2-3 sec; timeout if over 5 sec.
				pub.stream.Close()
				delete(publishers, conID)
				fmt.Println("Publishers:", len(publishers))
			}
		}

		// Wait 1 sec.

		time.Sleep(time.Second)
	}
}
