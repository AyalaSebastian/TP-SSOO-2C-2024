package utils

import (
	"sync"
)

var MutexPlanificador sync.Mutex
var Planificador = sync.NewCond(&MutexPlanificador)

type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore crea un nuevo semáforo con un contador inicial dado
func NewSemaphore(size int) *Semaphore {
	ch := make(chan struct{}, size)
	for i := 0; i < size; i++ {
		ch <- struct{}{} // Llena el canal (ocupado)
	}
	return &Semaphore{ch: ch}
}

func (s *Semaphore) Wait() {
	<-s.ch // Bloquea hasta que haya un permiso
}

func (s *Semaphore) Signal() {
	s.ch <- struct{}{} // Libera un permiso
}
