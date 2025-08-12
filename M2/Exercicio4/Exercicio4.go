package main

import (
	"fmt"
)

const N = 10

// N processos fazem fase1 depois todos fazem fase 2, em loop
func proc(id int, turns int, init bool, c1, c2 chan struct{}, ch_fim chan struct{}) {
	for i := 0; i < turns; i++ {
		fmt.Println(" fase 1 ", id)
		for j := 0; j < 2; j++ {
			if init { // um processo escreve em c2 antes de ler em c1
				c2 <- struct{}{}
				<-c1
			} else {
				<-c1 // todos demais leem de c1 e escrevem em c2
				c2 <- struct{}{}
			}
		}
		fmt.Println(" fase 2 ", id)
		for j := 0; j < 2; j++ {
			if init {
				c2 <- struct{}{}
				<-c1
			} else {
				<-c1
				c2 <- struct{}{}
			}
		}
	}
	ch_fim <- struct{}{}
}

func main() {
	var c [N]chan struct{}
	for i := 0; i < N; i++ {
		c[i] = make(chan struct{})
	}
	ch_fim := make(chan struct{})
	for i := 0; i < (N); i++ {
		go proc(i, 10, (i == 0), c[i], c[(i+1)%N], ch_fim)
	}
	for i := 0; i < N; i++ {
		<-ch_fim
	}
}
