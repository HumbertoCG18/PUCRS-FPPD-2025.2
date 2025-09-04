// ---------------------------------------
// PUCRS - Escola Politecnica
// Prof.  Fernando Dotti
// Concorrencia:  semaforos e barreiras em Go
// ---------------------------------------

package main

import (
	"fmt"
	"programasGo/FPPDSemaforo"
	// "golang.org/x/sync/semaphore"
)

// prompt>> go mod download golang ../FPPDSemaforo   ou  golang.org/x/sync/semaphore ou ...
// prompt>> go mod init ../FPPDSemaforo
// prompt>> go get ../FPPDSemaforo

// GO111MODULE=off go run ....go

// -----------------------------------------
// Definição de Barreira -------------------
// -----------------------------------------

// barreira:  N processos devem chegar em um ponto chamado barreira, para então todos seguirem adiante
// do￼￼ livro The Little book of Semaphores￼￼

// turnstile=Semaphore(0)
// turnstile2=Semaphore(1)
// mutex=Semaphore(1)

// Incrementa N e veja se todos chegaram. se chegaram, programa turnstile2 para bloquear processos depois de passar
// por turnstile, e sinaliza UM processo para prosseguir
// mutex.wait()
//    count += 1
// if count == n:
//    turnstile2.wait()       // turnstile2 permitia um processo, este passa e os demais bloqueiam
//    turnstile.signal()      // acorda um processo em turnstile
// mutex.signal()

// turnstile.wait()           // processo que acorda, libera outro:
// turnstile.signal()         // o padrao catraca de uso de semaforos

// #criticalpoint

// mutex.wait()
//    count -= 1
//    if count == 0:
//       turnstile.wait()    // ultimo processo da barreira gerou um signal a mais vide acima. aqui zera este semaforo para ninguem passar
//       turnstile2.signal() // ultimo processo libera a catraca turnstile2 e assim todos os processos podem reusar a barreira
// mutex.signal()
//
// turnstile2.wait()
// turnstile2.signal()
// mutex.wait()

type Barrier struct {
	mutex    *FPPDSemaforo.Semaphore
	catraca1 *FPPDSemaforo.Semaphore
	catraca2 *FPPDSemaforo.Semaphore
	count    int
	n        int
}

func NewBarrier(val int) *Barrier {
	b := &Barrier{
		mutex:    FPPDSemaforo.NewSemaphore(1),
		catraca1: FPPDSemaforo.NewSemaphore(0),
		catraca2: FPPDSemaforo.NewSemaphore(1),
		count:    0,
		n:        val,
	}

	return b
}

func (b *Barrier) Arrive() {
	b.mutex.Wait()
	b.count++
	if b.count == b.n {
		b.catraca2.Wait()
		b.catraca1.Signal()
		fmt.Println("   Ponto Critico") // para mostrar o ponto critico em relacao aos prints dos processos no exemplo abaixo
	}
	b.mutex.Signal()
	b.catraca1.Wait()
	b.catraca1.Signal()

}

func (b *Barrier) Leave() {
	b.mutex.Wait()
	b.count--
	if b.count == 0 {
		b.catraca1.Wait()
		b.catraca2.Signal()
		fmt.Println("                Ponto Não Critico")
	}
	b.mutex.Signal()
	b.catraca2.Wait()
	b.catraca2.Signal()
}

// -----------------------------------------
// Fim Definição de Barreira ---------------
// -----------------------------------------

// Exemplo de uso a barreira
// Os processos tem cada um um identificador
// eles escrevem o identificador antes e depois do ponto crítico
// como "   Ponto Critico" é escrito na operacao Arrive, pelo ultimo processo
// que chega no arrive, e como "   Ponto Nao Critico" pelo último que sai
// do ponto crítico, cada processo so imprime seu ID uma vez entre
// um e outro (ponto critico e nao critico)

func useBarrier(b *Barrier, i int, f chan struct{}) {
	for j := 0; j < 10; j++ {
		fmt.Print(i, " ")
		b.Arrive()
		fmt.Print(i, " ")
		b.Leave()
	}
	f <- struct{}{}
}

func procsComBarreiras() {
	b := NewBarrier(5)
	f := make(chan struct{})
	for i := 0; i < 5; i++ {
		go useBarrier(b, i, f)
	}
	for i := 0; i < 5; i++ {
		<-f
	}
}

// Fim do uso da barreira

func main() {
	procsComBarreiras()
}
