// por Fernando Dotti - PUCRS -
// 		assim como sobrepor tempo de comunicacao pode ser vantajoso,
// 		sobrepor tempos de processamento se a máquina dispõe de diversos núcleos
// 		pode também levar a ganhos.
//
// PROBLEMA:
//		encontre abaixo um programa sequencial que conta o numero de primos de um array.
// 		Se os números primos forem "grandes" o calculo se eles são primos torna-se CPU intensivo
// 		e o uso de diversos nucleos mostra o ganho de desempenho.
//
// EXERCICIO:
// 		1) torne este programa concorrente, sobrepondo temporalmente o trabalho de
//		   computar se um valor é primo
// OBSERVACAO: garanta com que seu gerador de numeros gere valores com diversas casas.   aprox 10.
//

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const N = 2000

func main() {
	fmt.Println("------ DIFFERENT prime count IMPLEMENTATIONS -------")
	slice := generateSlice(N)
	p := contaPrimosSeq(slice)
	fmt.Println("  ------ n primos :  ", p)
}

// Generates a slice of size, size filled with random numbers
func generateSlice(size int) []int {
	slice := make([]int, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(999999999)
	}
	return slice
}

func contaPrimosSeq(s []int) int {
	result := 0
	for i := 0; i < N; i++ {
		if isPrime(s[i]) {
			fmt.Println("  ------ primos :  ", s[i])
			result++
		}
	}
	return result
}

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
