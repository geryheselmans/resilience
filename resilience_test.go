package resilience

import (
	"testing"
	"time"
)

func TestDoNilFallback(t *testing.T) {
	var dummyStruct = struct{
		int
		string
	} {
		123,
		"test",
	}

	var tests = []struct {
		input interface {}
		sleep time.Duration
		want  interface{}
	} {
		{"Test", 0, "Test"},
		{"Test", 800, "Test"},
		{"Test", 1200, nil},
		{"Test", 1000, nil},
		{"Test", 998, "Test"},
		{123, 0, 123},
		{123, 800, 123},
		{123, 1200, nil},
		{123, 1000, nil},
		{123, 998, 123},
		{dummyStruct, 0, dummyStruct},
		{dummyStruct, 800, dummyStruct},
		{dummyStruct, 1200, nil},
		{dummyStruct, 1000, nil},
		{dummyStruct, 998, dummyStruct},
	}

	for _,test := range tests {
		result := Do(func(result chan<- interface{}){

			time.Sleep(test.sleep * time.Millisecond)

			result <- test.input
		}, nil)

		got := <- result

		if got != test.want {
			t.Errorf("%s @ %q: got %q, want %q,",
				test.input, test.sleep * time.Millisecond, got, test.want)
		}
	}
}
