package resilience

import (
	"time"
)

func Do(command func(chan<- interface{}), fallback func(chan<- interface{})) chan interface{} {
	returnChan := make(chan interface{})

	go executeCommand(returnChan, command, fallback)

	return returnChan
}

func checkAndExecute(commandKey string, semaphoreKey string, returnChan chan<- interface{}, command func(chan<- interface{}), fallback func(chan<- interface{})) {

}

func executeCommand(returnChan chan<- interface{}, command func(chan<- interface{}), fallback func(chan<- interface{})) {
	commandChan := make(chan interface{})
	go command(commandChan)

	timeOut := time.NewTimer(time.Second)

	select {
	case result := <-commandChan:
		returnChan <- result
	case <-timeOut.C:
		executeFallBack(returnChan, fallback)
	}
}

func executeFallBack(returnChan chan<- interface{}, fallback func(chan<- interface{})) {
	if fallback != nil {
		fallbackChan := make(chan interface{})

		go fallback(fallbackChan)

		returnChan <- <-fallbackChan
	} else {
		returnChan <- nil
	}
}
