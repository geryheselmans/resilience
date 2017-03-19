package ratel

import "time"

func Do(command func(<- chan interface{})) chan interface{} {
	returnChan := make(chan interface{})
	commandChan := make(chan interface{})

	timeOut := time.NewTimer(time.Second)

	select {
	case result := <- commandChan:
		returnChan <- result
	case timeOut.C:
		returnChan <- nil

	}

	return returnChan
}
