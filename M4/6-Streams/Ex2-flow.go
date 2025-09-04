// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
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

func separadorImparesPares(.....) {
	
}

func separadorPrimos(.....) {
	
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

	...  := make(chan int, sizeChan)
	...

	go consumer("pares ", ... )      					 // consome pares
	go consumer("               impares ", ... )		 // consome impares nao primos
	go consumer("                                      primos", ... )   // consome primos

	go separadorPrimos( .... )     // manda primos para canal de primos, e nao primos para canal de nao primos
	go separadorPrimos( ... )      // manda primos para canal de primos, e nao primos para canal de nao primos
	
	go separadorImparesPares( .... )  // manda primos para um separador de primos, e pares para canal de pares
	go separadorImparesPares( .... )  // manda primos para outro separador de primos, e pares para canal de pares

	go geraValores(... )    // manda para um separador de pares impares
	go geraValores(... )    // manda para o outro

	// bloqueia para nao acabar
	<-make(chan bool)

}
