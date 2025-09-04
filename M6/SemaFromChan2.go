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
type Semaphore struct {             // este semáforo implementa quaquer numero de creditos em "v"
	v int                 			// valor do semaforo: negativo significa proc bloqueado
	fila chan struct{}    			// canal para bloquear os processos se v < 0
	sc chan struct{}      			// canal para atomicidade das operacoes wait e signal
}

func NewSemaphore(init int) *Semaphore {
	s := &Semaphore{
		v: init,                    // valor inicial de creditos
		fila: make(chan struct{}),  // canal sincrono para bloquear processos
		sc: make(chan struct{}, 1), // usaremos este como semaforo para SC, somente 0 ou 1
	}
	return s
}

func (s *Semaphore) Wait() {
	s.sc <- struct{}{}  			// SC do semaforo feita com canal
	s.v--               			// decrementa valor
	if s.v < 0 {        			// se negativo era 0 ou menor, tem que bloquear
		<- s.sc             		// antes de bloq, libera acesso
		s.fila <- struct{}{} 		// bloqueia proc
	} else { 
		<- s.sc              		// libera acesso
	}
}

func (s *Semaphore) Signal() {
	s.sc <- struct{}{}  			// entra sc
	s.v++
	if s.v <= 0 {        			// tem processo bloqueado ?
	   <- s.fila        			// desbloqueia
	}
	<- s.sc         				// libera SC para outra op
}
// ----------------------------------

func useSC(shared *int, s *Semaphore, fin chan struct{}, iter int) {
	for i:=0; i<iter; i++ {  		// cada processo faz iteracoes
		s.Wait()             		// entrada sc com nosso semaforo
		*shared = ((*shared) + 1)   // SC
		s.Signal()           		// saida SC com nosso semaforo
	}
	fin <- struct{}{}        		// sinaliza termino
}

func main() {
	n := 50             			// nro de processos
	iter := 10000       			// nro de iteracoes em cada um

	var fin chan struct{} = make(chan struct{})
	s := NewSemaphore(1)
	sh := 0

	for  i:=0;i<n;i++ {
		go useSC(&sh, s, fin, iter)
	}
	for  i:=0;i<n;i++ {
		<- fin
		fmt.Println("fin    ", i)
	}
	fmt.Println(sh)
}
