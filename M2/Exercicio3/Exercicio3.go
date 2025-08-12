package main

import (
	"fmt"
)

func A(turns int, c1, c2 chan struct{}, ch_fim chan struct{}) {
	for i := 0; i < turns; i++ {
		fmt.Print("A")
		c1 <- struct{}{}
		<-c2
	}
	ch_fim <- struct{}{}
}

func B(turns int, c1, c2 chan struct{}, ch_fim chan struct{}) {
	for i := 0; i < turns; i++ {
		<-c1
		fmt.Print("B")
		c2 <- struct{}{}
	}
	ch_fim <- struct{}{}
}

func C(turns int, c1, c2 chan struct{}, ch_fim chan struct{}) {
	for i := 0; i < turns; i++ {
		<-c2
		<-c2
		fmt.Println("C")
		c1 <- struct{}{}
		c1 <- struct{}{}
	}
	ch_fim <- struct{}{}
}

func main() {
	c1 := make(chan struct{})
	c2 := make(chan struct{})
	ch_fim := make(chan struct{})

	go A(50, c1, c2, ch_fim)
	go B(50, c1, c2, ch_fim)
	go C(50, c1, c2, ch_fim)

	fmt.Println("Criei")
	for i := 0; i < 2; i++ {
		<-ch_fim
	}
}
