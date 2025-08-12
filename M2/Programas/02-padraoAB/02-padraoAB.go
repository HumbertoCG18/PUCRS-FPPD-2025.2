package main

import (
	"fmt"
)

// tentativa de programa com saida ABABABABABAB ...

func A(turns int, c1 chan struct{}, ch_fim chan struct{}){
	for i:=0; i<turns; i++ {
		fmt.Print("A")
		c1 <- struct{}{}
	}
	ch_fim <- struct{}{}
}

func B_AB(turns int, c1 chan struct{}, ch_fim chan struct{}){
	for i:=0; i<turns; i++ {
		<- c1 
		fmt.Print("B")
	}
	ch_fim <- struct{}{}
}

func main() {
	c1      := make(chan struct{})
	ch_fim  := make(chan struct{})

	go A(50,c1,ch_fim)
	go B_AB(50,c1,ch_fim)

	fmt.Println("criei")
	for i := 0; i < 1; i++ {     // espera os 100 processos acabarem
		<-ch_fim
	}
}
