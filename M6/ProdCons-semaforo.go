// Disciplina de Modelos de Computacao Concorrente
// Escola Politecnica - PUCRS
// Prof.  Fernando Dotti

// Uso de semaforo no modelo produtor consumidor

package main

import (
	"fmt"

	"programasGo/FPPDSemaforo"
)

// a estrutura para conter os valores
const N = 10

var buffer [N]int

// posicao para entrada e para saida
var in int = 0
var out int = 0

// os canais de sincronizacao
var mx *FPPDSemaforo.Semaphore = FPPDSemaforo.NewSemaphore(1)
var naoVazio *FPPDSemaforo.Semaphore = FPPDSemaforo.NewSemaphore(0) // sema para bloquear consumidor em caso de vazio
var naoCheio *FPPDSemaforo.Semaphore = FPPDSemaforo.NewSemaphore(N) // sema para bloquear produtor se estiver cheio

func insere(x int) {
	naoCheio.Wait()
	mx.Wait()
	buffer[in] = x
	in = (in + 1) % N // ajusta marcador de entrada
	fmt.Println("ins: ", x)
	mx.Signal()
	naoVazio.Signal()
}

func retira() int {
	naoVazio.Wait()
	mx.Wait()
	c := buffer[out]
	out = (out + 1) % N
	fmt.Println("                retira ", c)
	mx.Signal()
	naoCheio.Signal()
	return c
}

// --- processos produtores e consumidores

func produtor(f chan struct{}) {
	for i := 0; i < 100; i++ {
		insere(i)
	}
	f <- struct{}{}
}

func consumidor(f chan struct{}) {
	for i := 0; i < 100; i++ {
		retira()
	}
	f <- struct{}{}
}

func main() {

	fim := make(chan struct{})

	for i := 0; i < 10; i++ {
		go produtor(fim)
		go consumidor(fim)
	}

	for i := 0; i < 20; i++ {
		<-fim
	}
}
