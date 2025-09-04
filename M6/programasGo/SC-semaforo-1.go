// Disciplina de Modelos de Computacao Concorrente
// Escola Politecnica - PUCRS
// Prof.  Fernando Dotti

// Uso de semaforo para Secao Critica, em GO
// definicao de semaforo ee a ja fornecida pelo professor
// agora em formato package
//
// coloque o pacote FPPDSemaforo dentro de um diretorio
// chamado FPPDSemaforo, no diretorio corrente (onde esta seu codigo)
// faca
// import (
//	"./FPPDSemaforo"
// )
// declare um semaforo por exemplo:
//      FPPDSemaforo.Semaphore s = FPPDSemaforo.NewSemaphore(1)

package main

import (
	"fmt"
	"programasGo/FPPDSemaforo"
)

func useSC(shared *int, s *FPPDSemaforo.Semaphore, fim chan int) {
	for i := 0; i < 100; i++ {
		s.Wait()                           // entrada na SC
		*shared = ((*shared) + 1) % 100000 // SC
		s.Signal()                         // saida da SC
	}
	fim <- 1
}

func main() {
	var ch_fim chan int = make(chan int)

	s := FPPDSemaforo.NewSemaphore(1) // semaforo iniciado em 1 para SC
	sh := 0                           // A VARIAVEL COMPARTILHADA!!
	// -----------------------
	for i := 0; i < 100; i++ { // inicia 100 processos que usam SC
		go useSC(&sh, s, ch_fim) // todos usam sh, o mesmo semaroro e sinalizam fim em ch_fim
	}
	for i := 0; i < 100; i++ { // espera os 100 processos acabarem
		<-ch_fim
	}
	fmt.Println("Resultado com semaforos para SC  ", sh)
}
