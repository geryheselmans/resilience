package resilience

type token struct{}

type Semaphore struct {
	name    string
	counter chan token
}

func NewSemaphore(name string) *Semaphore {
	created := make(chan token, 10)

	for i := 0; i < 10; i++ {
		created <- token{}
	}

	return &Semaphore{
		name,
		created,
	}
}

func (s *Semaphore) getCounter() chan token {
	return s.counter
}
