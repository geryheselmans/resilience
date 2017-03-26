package resilience

import (
	"testing"
)

func TestGetSemaphore(t *testing.T) {
	sem1 := NewSemaphore("sem1")

	for i := 1; i <= 11; i++ {
		var index = i
		select {
		case <-sem1.getCounter():
			if index > 10 {
				t.Error("Got 11the token")
			}
		default:
			if index != 11 {
				t.Error("Not enough tokens")
			}
		}
	}
}
