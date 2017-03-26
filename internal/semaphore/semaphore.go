package semaphore

import (
	"sync"
)

type Token struct{}
type Semaphore chan Token

var initMapOnce sync.Once
var mu sync.RWMutex
var keyToSemaphoreMap map[string]Semaphore

func initMap() {
	keyToSemaphoreMap = make(map[string]Semaphore)
}

func GetSemaphore(semaphoreKey string) Semaphore {
	initMapOnce.Do(initMap)

	mu.RLock()
	if _, ok := keyToSemaphoreMap[semaphoreKey]; ok {
		mu.RUnlock()
		return keyToSemaphoreMap[semaphoreKey]
	}
	mu.RUnlock()

	mu.Lock()
	if _, ok := keyToSemaphoreMap[semaphoreKey]; !ok {
		keyToSemaphoreMap[semaphoreKey] = createSemaphore()
	}
	result := keyToSemaphoreMap[semaphoreKey]
	mu.Unlock()

	return result
}
func createSemaphore() Semaphore {
	created := make(Semaphore, 10)

	for i := 0; i < 10; i++ {
		created <- Token{}
	}

	return created
}
