// Disciplina de Modelos de Computacao Concorrente
// Escola Politecnica - PUCRS
// Prof.  Fernando Dotti
// ATENCAO:   Aqui estamos discutindo como construir a
// semantica de semaforos a partir do uso de canais.
//
// As funcoes s.wait() e s.signal() tem o mesmo significado 
// de semaforos da literatura?

package main

import (
	"fmt"
)

// ---------------------------------
type Semaphore struct {
	sChan chan struct{}
}

func NewSemaphore(init int) *Semaphore {
	s := &Semaphore{
		sChan: make(chan struct{}, init),
	}
	return s
}

func (s *Semaphore) Wait() {
	s.sChan <- struct{}{}
}

func (s *Semaphore) Signal() {
	<- s.sChan 
}
// ----------------------------------

func useSC(shared *int, s *Semaphore, fin chan struct{}, iter int) {
	for i:=0; i<iter; i++ {
		s.Wait()
		*shared = ((*shared) + 1) 
		s.Signal()
	}
	fin <- struct{}{}
}

func main() {
	n := 50
	iter := 10000

	var fin chan struct{} = make(chan struct{})
	s := NewSemaphore(1)
	sh := 0

	for  i:=0;i<n;i++ {
		go useSC(&sh, s, fin, iter)
	}
	for  i:=0;i<n;i++ {
		<- fin
	}
	fmt.Println(sh)
}
