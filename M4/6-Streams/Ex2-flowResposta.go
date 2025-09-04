// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// suponha o seguinte sistema
// dois geradores mandam dados, cada um para um separador de pares e impares
// os separadores alimentam os mesmos canais: um de pares e um de impares
// sobre o canal de impares, um outro separa os numeros primos dos não primos
// ao final temos tres consumidores, o de pares, o de impares não primos, e os primos
//
// EXERCICIO:
//   1)  Complete o programa abaixo para se comportar como a descricao acima.
//       Note que foram fornecidas funcoes para calculo de par e primo.
//       Voce deve construir os processos concorrentes comunicando por canais,
//       em varios estagios.
//
//  SOLUCAO possivel abaixo.
//

package main

import (
	"math/rand"
	"time"
)

const N = 200
const sizeChan = 1

func geraValores(cOut chan int) {
	for i := 0; i < N; i++ {
		numero := rand.Int31n(2000) + 1
		cOut <- int(numero)
	}
}

func separadorImparesPares(in chan int, pares chan int, impares chan int) {
	v := 0
	for {
		v = <-in
		if ehPar(v) {
			pares <- v
		} else {
			impares <- v
		}
	}
}

func separadorPrimos(inImpares chan int, primos chan int, nPrimos chan int) {
	v := 0
	for {
		v = <-inImpares
		if ehPrimo(v) {
			primos <- v
		} else {
			nPrimos <- v
		}
	}
}

func consumer(nome string, rec chan int) {
	for {
		dado := <-rec
		println(nome, dado)
	}
}

func ehPar(v int) bool {
	return v%2 == 0
}

func ehPrimo(v int) bool {
	if ehPar(v) {
		return false
	}
	for i := 3; i*i <= v; i += 2 {
		if v%i == 0 {
			return false
		}
	}
	return true
}

func main() {

	rand.Seed(time.Now().UnixNano())

	separadorPIch1 := make(chan int, sizeChan)
	separadorPIch2 := make(chan int, sizeChan)

	separadorPrimosCh1 := make(chan int, sizeChan)
	separadorPrimosCh2 := make(chan int, sizeChan)

	pares := make(chan int, sizeChan)
	nprimos := make(chan int, sizeChan)
	primos := make(chan int, sizeChan)

	go consumer("pares ", pares)
	go consumer("               impares ", nprimos)
	go consumer("                                      primos", primos)

	go separadorImparesPares(separadorPIch1, pares, separadorPrimosCh1)
	go separadorImparesPares(separadorPIch2, pares, separadorPrimosCh2)

	go separadorPrimos(separadorPrimosCh1, primos, nprimos)
	go separadorPrimos(separadorPrimosCh2, primos, nprimos)

	go geraValores(separadorPIch1)
	go geraValores(separadorPIch2)

	// bloqueia para nao acabar
	<-make(chan bool)

}
