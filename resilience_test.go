package resilience

import (
	"testing"
	"time"
)

func TestOneDoSingleResponse(t *testing.T) {
	var dummyStruct = struct{
		int
		string
	} {
		123,
		"test",
	}

	stringFallBack := func(result chan<- interface{}){
		result <- "Fallback"
	}

	intFallBack := func(result chan<- interface{}){
		result <- 456
	}

	dummyStructFallback := func(result chan<- interface{}){
		result <- nil
	}

	var tests = []struct {
		input interface {}
		sleep time.Duration
		fallback func(result chan<- interface{})
		want  interface{}
	} {
		{"Test", 0, nil, "Test"},
		{"Test", 800, nil,"Test"},
		{"Test", 1200, nil,nil},
		{"Test", 1002, nil,nil},
		{"Test", 998, nil,"Test"},
		{"Test", 0, stringFallBack, "Test"},
		{"Test", 800, stringFallBack,"Test"},
		{"Test", 1200, stringFallBack,"Fallback"},
		{"Test", 1002, stringFallBack,"Fallback"},
		{"Test", 998, stringFallBack,"Test"},
		{123, 0, nil,123},
		{123, 800, nil,123},
		{123, 1200, nil,nil},
		{123, 1002, nil,nil},
		{123, 998, nil,123},
		{123, 0, intFallBack,123},
		{123, 800, intFallBack,123},
		{123, 1200, intFallBack,456},
		{123, 1002, intFallBack,456},
		{123, 998, intFallBack,123},
		{dummyStruct, 0, nil,dummyStruct},
		{dummyStruct, 800,nil, dummyStruct},
		{dummyStruct, 1200, nil,nil},
		{dummyStruct, 1002, nil,nil},
		{dummyStruct, 998, nil,dummyStruct},
		{dummyStruct, 0, dummyStructFallback,dummyStruct},
		{dummyStruct, 800,dummyStructFallback, dummyStruct},
		{dummyStruct, 1200, dummyStructFallback,nil},
		{dummyStruct, 1002, dummyStructFallback,nil},
		{dummyStruct, 998, dummyStructFallback,dummyStruct},
	}

	for _,test := range tests {
		start := time.Now()
		result := Do(func(result chan<- interface{}){

			time.Sleep(test.sleep * time.Millisecond)

			result <- test.input
		}, test.fallback)


		got := <- result
		stop:= time.Now()

		diff := stop.Sub(start)

		maxDiff := (test.sleep * time.Millisecond) + (5 * time.Millisecond)

		if diff > maxDiff {
			t.Errorf("%q @ %q: got %q must be max %q",
				test.input, test.sleep * time.Millisecond, diff , maxDiff)
		}

		if got != test.want {
			t.Errorf("%q @ %q: got %q, want %q,",
				test.input, test.sleep * time.Millisecond, got, test.want)
		}
	}
}
