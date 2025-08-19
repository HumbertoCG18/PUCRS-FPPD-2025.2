// Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// >>> Veja antes Ex4 desta série.
// EXERCICIOS:
//     1) esta é uma solução para a questão anterior ?
//     2) o que garante que todos valores serão lidos antes do programa acabar ?

package main

import "fmt"

func main() {
	ch := make(chan int)
	quit := make(chan bool)
	go shower(ch, quit)
	for i := 0; i < 1000; i++ {
		ch <- i
	}
	quit <- false // or true, does not matter
}

func shower(c chan int, quit chan bool) {
	for {
		select {
		case j := <-c:
			fmt.Printf("%d\n", j)
		case <-quit:
			break
		}
	}
}
