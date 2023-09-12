package main

import "testing"

func Test_randomMessage(t *testing.T) {
	t.Run("LengthCorrect", func(t *testing.T) {
		msg := randomMessage()
		if len(msg) == 0 || len(msg) > 20 {
			t.Errorf("randomMessage() = %v (length %v)", msg, len(msg))
		}
	})

	t.Run("MessagesDifferent", func(t *testing.T) {
		msg1 := randomMessage()
		msg2 := randomMessage()

		if msg1 == msg2 {
			t.Errorf("randomMessage() = %v (equal)", msg1)
		}
	})
}
