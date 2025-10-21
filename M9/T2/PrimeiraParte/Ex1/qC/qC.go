package qc

import (
	"fmt"
)

const X = 40

var ch [X]chan struct{}

func fA(id int, in <-chan struct{}, out chan<- struct{}){
	for{
		<- in
		fmt.Println(id)
		out <- struct{}{}
	}
}

func questaoC(){
	for i := 0; i < X; i++ {
		ch[i] = make(chan struct{})
	}

	for i := 0; i < X; i++{
		go fA(i, ch[i], ch[(i+1)%X])
	}

	ch[0] <- struct{}{}

	blq := make(chan struct{})
	<- blq
}

func main(){
	questaoC()
}