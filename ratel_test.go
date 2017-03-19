package ratel

import (
	"testing"
	"time"
)

func TestDoNilFallback(t *testing.T) {
	var tests = []struct {
		toReturn string
		sleep time.Duration
		expected interface{}
	} {
		//{"Test", 0, "Test"},
		{"Test", 1200, nil},
		{"Test", 1000, "Test"},
	}

	for _,test := range tests {
		result := Do(func(result chan<- interface{}){

			time.Sleep(test.sleep * time.Millisecond)

			result <- test.toReturn
		}, nil)

		actual := <- result

		if actual != test.expected {
			t.Errorf("Result not correct, want %s, got %s", actual, result)
		}
	}
}
