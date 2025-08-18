// por Fernando Dotti - PUCRS -
// 		o programa abaixo soma os numeros primos até um dado número N de entrada.
//      Note que ele avalia sequencialmente em um processo os numeros ate N e, cada vez que encontra
//      que encontra um primo, envia este para outro processo que faz a soma de todos.
//      Entao basicamente sao dois processos.
// EXERCICIO:
//     1)  como voce pode modificar este programa para que
//         o calculo de primos seja sobreposto temporalmente??
//
// SOLUCAO:  abaixo uma solucao possivel.

package main

import (
	"fmt"
)

// eh primo?  se for retorna em result, sempre avisa termino em fin
func isPrime(p int, result chan int, fin chan struct{}) {
	if p%2 == 0 {
		fin <- struct{}{}
		return
	}
	for i := 3; i*i <= p; i += 2 {
		if p%i == 0 {
			fin <- struct{}{}
			return
		}
	}
	result <- p
	fin <- struct{}{}
}

// Produz numeros primos de 2 ate <= n no cana results.
// lanca rotina para detectar se cada valor eh primo
// espera todas acabarem
// fecha canal
func primesUpTo(n int, results chan int) {
	fin := make(chan struct{})
	results <- 2
	for p := 3; p <= n; p += 2 {
		go isPrime(p, results, fin)
	}
	for p := 3; p <= n; p += 2 {
		<-fin
	}
	close(results)
}

// faz a soma dos primos de 2 ate <=n
// cada primo detectado é printed na tela
// depois a soma
func addPrimesTo(n int) (total int) {
	primes := make(chan int)
	go primesUpTo(n, primes)
	for p := range primes {
		fmt.Print(" ", p, " ")
		total += p
	}
	return total
}

func main() {
	fmt.Println(addPrimesTo(100))
}
