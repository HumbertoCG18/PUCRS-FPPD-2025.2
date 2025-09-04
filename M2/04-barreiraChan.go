package main

import (
	"fmt"
)

// tentativa de programa com saidas concorrentes A e B, depois C

func proc(id int, turns int, init bool, c1,c2 chan struct{}, ch_fim chan struct{}){
	for i:=0; i<turns; i++ {
		fmt.Println(" fase 1 ", id)
		for j:=0; j<2; j++ {
			if init{            // um processo escreve em c2 antes de ler em c1
				c2 <- struct{}{}
				<-c1
			} else {
				<- c1           // todos demais leem de c1 e escrevem em c2
				c2 <- struct{}{}		
			}
		}
		fmt.Println("              fase 2 ", id)
		for j:=0; j<2; j++ {
			if init{
				c2 <- struct{}{}
				<-c1
			} else {
				<- c1 
				c2 <- struct{}{}		
			}
		}
	}
	ch_fim <- struct{}{}
}

const N = 10

func main() {
	var c[N] chan struct{}
	for i := 0; i < N; i++ {
		c[i] = make(chan struct{})
	}
	ch_fim  := make(chan struct{})
	for i := 0; i < (N); i++ {
		go proc(i, 10, (i==0), c[i], c[(i+1)%N],ch_fim)
	}                          // processos e canais formam um anel lógico

	for i := 0; i < N; i++ {   // espera os 100 processos acabarem
		<-ch_fim
	}
}

