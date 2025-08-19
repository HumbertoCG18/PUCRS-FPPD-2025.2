// por Fernando Dotti - PUCRS -
// 		o programa abaixo soma os numeros primos até um dado número N de entrada.
//      Note que ele avalia sequencialmente em um processo os numeros ate N e, cada vez que encontra
//      que encontra um primo, envia este para outro processo que faz a soma de todos.
//      Entao basicamente sao dois processos.
// EXERCICIO:
//     1)  como voce pode modificar este programa para que
//         o calculo de primos seja sobreposto temporalmente??

package main

import (
	"fmt"
)

// Is p prime?
func isPrime(p int) bool {
	if p%2 == 0 {
		return false
	}
	for i := 3; i*i <= p; i += 2 {
		if p%i == 0 {
			return false
		}
	}
	return true
}

// Produce prime numbers 2 to <=n on the channel.
func primesUpTo(n int, results chan int) {
	results <- 2
	for p := 3; p <= n; p += 2 {
		if isPrime(p) {
			results <- p
		}
	}
	close(results)
}

// Return the sum of primes from 2 up to <=n.
func addPrimesTo(n int) (total int) {
	primes := make(chan int)
	go primesUpTo(n, primes)
	for p := range primes {
		total += p
	}
	return total
}

func main() {
	fmt.Println(addPrimesTo(100))
}
