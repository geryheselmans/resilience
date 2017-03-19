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
		maxDiff time.Duration
	} {
		{"Test", 0, nil, "Test", 5 * time.Millisecond},
		{"Test", 800, nil,"Test", 805 * time.Millisecond},
		{"Test", 1200, nil,nil,1005 * time.Millisecond},
		{"Test", 1002, nil,nil,1005 * time.Millisecond},
		{"Test", 998, nil,"Test",1005 * time.Millisecond},
		{"Test", 0, stringFallBack, "Test", 5 * time.Millisecond},
		{"Test", 800, stringFallBack,"Test", 805 * time.Millisecond},
		{"Test", 1200, stringFallBack,"Fallback",1005 * time.Millisecond},
		{"Test", 1002, stringFallBack,"Fallback",1005 * time.Millisecond},
		{"Test", 998, stringFallBack,"Test",1005 * time.Millisecond},
		{123, 0, nil,123, 5 * time.Millisecond},
		{123, 800, nil,123, 805 * time.Millisecond},
		{123, 1200, nil,nil,1005 * time.Millisecond},
		{123, 1002, nil,nil,1005 * time.Millisecond},
		{123, 998, nil,123,1005 * time.Millisecond},
		{123, 0, intFallBack,123, 5 * time.Millisecond},
		{123, 800, intFallBack,123, 805 * time.Millisecond},
		{123, 1200, intFallBack,456,1005 * time.Millisecond},
		{123, 1002, intFallBack,456,1005 * time.Millisecond},
		{123, 998, intFallBack,123,1005 * time.Millisecond},
		{dummyStruct, 0, nil,dummyStruct, 5 * time.Millisecond},
		{dummyStruct, 800,nil, dummyStruct, 805 * time.Millisecond},
		{dummyStruct, 1200, nil,nil,1005 * time.Millisecond},
		{dummyStruct, 1002, nil,nil,1005 * time.Millisecond},
		{dummyStruct, 998, nil,dummyStruct,1005 * time.Millisecond},
		{dummyStruct, 0, dummyStructFallback,dummyStruct, 5 * time.Millisecond},
		{dummyStruct, 800,dummyStructFallback, dummyStruct, 805 * time.Millisecond},
		{dummyStruct, 1200, dummyStructFallback,nil,1005 * time.Millisecond},
		{dummyStruct, 1002, dummyStructFallback,nil,1005 * time.Millisecond},
		{dummyStruct, 998, dummyStructFallback,dummyStruct,1005 * time.Millisecond},
	}

	for _,test := range tests {
		start := time.Now()
		result := Do(func(resultChan chan<- interface{}){

			time.Sleep(test.sleep * time.Millisecond)

			resultChan <- test.input
		}, test.fallback)

		got := <- result
		stop:= time.Now()

		diff := stop.Sub(start)

		if diff > test.maxDiff {
			t.Errorf("%q @ %q: got %q must be max %q",
				test.input, test.sleep * time.Millisecond, diff , test.maxDiff)
		}

		if got != test.want {
			t.Errorf("%q @ %q: got %q, want %q,",
				test.input, test.sleep * time.Millisecond, got, test.want)
		}
	}
}
