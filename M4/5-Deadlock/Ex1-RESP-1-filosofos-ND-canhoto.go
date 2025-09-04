// PUCRS - Fernando Dotti
// Exercício:
//   1) que condição de coffman é quebrada com a solução abaixo ?
//   2) argumente em portugues porque esta solução não tem deadlock.
//   3) voce poderia ter mais filósofos canhotos ?
//           implemente e teste.

package main

import (
	"fmt"
	"strconv"
)

const (
	PHILOSOPHERS = 5
	FORKS        = 5
)

func philosopher(id int, first_fork chan struct{}, second_fork chan struct{}) {
	for {
		fmt.Println(strconv.Itoa(id) + " senta\n")
		<-first_fork // pega
		<-second_fork
		fmt.Println(strconv.Itoa(id) + " come\n")
		first_fork <- struct{}{} // devolve
		second_fork <- struct{}{}
		fmt.Println(strconv.Itoa(id) + " levanta e pensa \n")
	}
}

func main() {
	var fork_channels [FORKS]chan struct{}

	for i := 0; i < FORKS; i++ {
		fork_channels[i] = make(chan struct{}, 1)
		fork_channels[i] <- struct{}{} // no inicio garfo esta livre
	}
	for i := 0; i < (PHILOSOPHERS - 1); i++ {
		fmt.Println("Filosofo " + strconv.Itoa(i) + " destro!")
		go philosopher(i, fork_channels[i], fork_channels[(i+1)])
	}
	fmt.Println("Filosofo " + strconv.Itoa(PHILOSOPHERS-1) + " canhoto!")
	go philosopher(PHILOSOPHERS-1, fork_channels[0], fork_channels[FORKS-1])

	<-make(chan struct{})
}
