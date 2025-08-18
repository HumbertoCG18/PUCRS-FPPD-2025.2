// PUCRS - Fernando Dotti
// Exercício:
//   1) que condição de coffman é quebrada com a solução abaixo ?
//   2) que outro problema ela introduz ?

package main

import (
	"fmt"
)

const (
	PHILOSOPHERS = 5
	FORKS        = 5
)

func philosopher(id int, left chan struct{}, right chan struct{}) {
	// Made by Niccolas Maganeli
	for {
		select {
		case l := <-left:
			// Took left fork
			select {
			case r := <-right:
				// Took right fork
				fmt.Printf("Filoso %d pegou ambos garfos\n", id)
				// Time to eat
				fmt.Printf("Filosofo %d terminou de comer\n", id)
				// Releasing both forks
				left <- l
				right <- r
			default:
				// Right fork wasn't available
				// Release left fork
				left <- l
				break
			}
		default:
			// Left fork wasn't available
			// Philosopher will think for a moment
			fmt.Printf("Filosofo %d esta pensando\n", id)
			fmt.Printf("Filosofo %d terminou de pensar\n", id)
			break
		}
	}
}

func main() {
	var fork_channels [FORKS]chan struct{}
	for i := 0; i < FORKS; i++ {
		fork_channels[i] = make(chan struct{}, 1)
		fork_channels[i] <- struct{}{} // no inicio garfo esta livre
	}
	for i := 0; i < (PHILOSOPHERS - 1); i++ {
		go philosopher(i, fork_channels[i], fork_channels[(i+1)%PHILOSOPHERS])
	}
	var blq chan struct{} = make(chan struct{})
	<-blq
}
