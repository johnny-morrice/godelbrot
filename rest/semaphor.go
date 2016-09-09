package rest

type sem chan bool

func semaphor(concurrent uint) sem {
	return sem(make(chan bool, concurrent))
}

func (s sem) acquire(n uint) {
	// Should the whole action be atomic?
	for i := uint(0); i < n; i++ {
		s <- true
	}
}

func (s sem) release(n uint) {
	for i := uint(0); i < n; i++ {
		<-s
	}
}
