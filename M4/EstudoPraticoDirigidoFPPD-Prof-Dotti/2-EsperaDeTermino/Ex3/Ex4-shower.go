// Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// Exemplo da Internet
// EXERCICIOS:
//   1) rode o programa abaixo e interprete.
//      todos os valores escritos no canal são lidos?
//   2) como isto poderia ser resolvido ?

package main

import "fmt"

func main() {
	ch := make(chan int, 5)
	go shower(ch)
	for i := 0; i < 10; i++ {
		ch <- i
	}
}

func shower(c chan int) {
	for {
		j := <-c
		fmt.Printf("%d\n", j)
	}
}
