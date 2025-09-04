// MCC - Fernando Dotti

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
		fmt.Println(strconv.Itoa(id) + " senta")
		<-first_fork // pega
		fmt.Println(strconv.Itoa(id) + " pegou direita")
		<-second_fork
		fmt.Println(strconv.Itoa(id) + " come")
		first_fork <- struct{}{} // devolve
		second_fork <- struct{}{}
		fmt.Println(strconv.Itoa(id) + " levanta e pensa")
	}
}

func main() {
	var fork_channels [FORKS]chan struct{}
	for i := 0; i < FORKS; i++ {
		fork_channels[i] = make(chan struct{}, 1)
		fork_channels[i] <- struct{}{} // no inicio garfo esta livre
	}
	for i := 0; i < (PHILOSOPHERS); i++ {
		fmt.Println("Filosofo " + strconv.Itoa(i))
		go philosopher(i, fork_channels[i], fork_channels[(i+1)%PHILOSOPHERS])
	}
	var blq chan struct{} = make(chan struct{})
	<-blq
}
