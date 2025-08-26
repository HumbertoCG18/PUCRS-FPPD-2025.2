package main

import (
	"fmt"
	"math/rand"
)

const (
	NCL = 10
)

type Request struct {
	v      int
	ch_ret chan int
}

// ------------------------------------
// cliente
func cliente(i int, req chan Request) {
	var v, r int
	my_ch := make(chan int)
	for {
		v = rand.Intn(1000)
		req <- Request{v, my_ch}
		r = <-my_ch
		fmt.Println("cli: ", i, " req: ", v, "  resp:", r)
	}
}

// ------------------------------------
// servidor sequencial
func servidorSeq(in chan Request) {
	for {
		req := <-in
		fmt.Println("                       trataReq ", req)
		req.ch_ret <- req.v * 2 // responde  ao cliente
	}
}

// ------------------------------------
// main
func main() {
	fmt.Println("------ Servidores Sequencial -------")
	serv_chan := make(chan Request)
	for i := 0; i < NCL; i++ {
		go cliente(i, serv_chan)
	}
	servidorSeq(serv_chan)
}
