package main

import (
	"fmt"
)

const N = 10

type Topology [N][N]int

type Message struct {
	id int
}

type inputChan [N]chan Message

type nodeStruct struct {
	id   int
	topo Topology
	inCh inputChan
}

const channelBufferSize = 1

func (n *nodeStruct) broadcast(m Message) {
	for j := 0; j < N; j++ {
		if n.topo[n.id][j] == 1 {
			n.inCh[j] <- m
		}
	}
}

func (n *nodeStruct) nodo() {
	fmt.Println(n.id, " ativo! ")
	for {
		m := <-n.inCh[n.id]
		fmt.Println(n.id, " tratando ", m)
		n.broadcast(m)
	}
}

func main() {
	var topo Topology
	topo = [N][N]int{
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 0           0 -> 1
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0}, // 1           1 -> 2
		{0, 0, 0, 1, 0, 0, 0, 0, 0, 0}, // 2           2 -> 3
		{0, 0, 0, 0, 1, 0, 0, 0, 1, 0}, // 3           3 -> 4 e  3 -> 7
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 1}, // 4           4 -> 5 e  4 -> 9
		{0, 0, 0, 0, 0, 0, 1, 0, 0, 0}, // 5           5 -> 6
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 0}, // 6           6 -> 7
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 0}, // 7           7 -> 8
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, // 8           8 -> 9
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}} // 9

	var inCh inputChan // cada nodo i tem um canal de entrada, chamado inCh[i]
	for i := 0; i < N; i++ {
		inCh[i] = make(chan Message, channelBufferSize) // criando cada um dos canais
	}

	// lanca todos os nodos
	for id := 0; id < N; id++ {
		n := nodeStruct{id, topo, inCh}
		go n.nodo()
	}

	for i := 1; i < 2; i++ { // gera mensagem de teste a cada segundo
		inCh[0] <- Message{i}
		//time.Sleep(time.Second)
	}
	<-make(chan struct{}) // bloqueia senao nodos acabam
}
