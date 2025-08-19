// PUCRS - Fernando Dotti
// Filosofos.  Cada vez que filosofo acaba de comer, levanta da mesa
// e pensa caminhando.   Quando fica com fome, tenta sentar novamente.
// Mas um garçom (waiter) deixa somente 4 filósofos estarem na mesa.
// Exercicios:
//     Veja como o waiter foi modelado.
//     Note o uso de canais como sendo um conjunto de creditos e
//     como o credito é tomado e devolvido pelo filosofo.

package main

import (
	"fmt"
	"strconv"
)

const (
	PHILOSOPHERS = 5
	FORKS        = 5
	PLACES       = 4
)

func philosopher(id int, first_fork chan struct{}, second_fork chan struct{}, waiterPlaces chan struct{}) {
	for {
		<-waiterPlaces // posso sentar ?
		fmt.Println(strconv.Itoa(id) + " senta\n")
		<-first_fork  // pega um ...
		<-second_fork // e o outro
		fmt.Println(strconv.Itoa(id) + " come\n")
		first_fork <- struct{}{}  // devolve os ...
		second_fork <- struct{}{} // dois garfos
		fmt.Println(strconv.Itoa(id) + " levanta e pensa \n")
		waiterPlaces <- struct{}{} // levantei para caminhar
	}
}

func main() {
	var fork_channels [FORKS]chan struct{}
	for i := 0; i < FORKS; i++ {
		fork_channels[i] = make(chan struct{}, 1)
		fork_channels[i] <- struct{}{} // no inicio garfo esta livre
	}
	waiter := make(chan struct{}, PLACES)
	for i := 0; i < PLACES; i++ {
		waiter <- struct{}{} // inicia waiter com 4 lugares livres
	}
	for i := 0; i < (PHILOSOPHERS - 1); i++ {
		go philosopher(i, fork_channels[i], fork_channels[(i+1)%PHILOSOPHERS], waiter)
	}
	<-make(chan struct{})
}
