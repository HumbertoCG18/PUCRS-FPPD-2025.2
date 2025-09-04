// -----------------------------------------
// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// Este programa apresenta uma arcabouço para implementar o problema da rede em anel (token ring)
// PROBLEMA:
//    vide descricao em Ex2-TokenRingExplanacao.pdf  nesta pasta
// OBSERVACAO
//    note que os processos usuario, nodo(ou estacao nos slides), e respectivos canais
//    ja estao montados.
// EXERCICIO:
//    1) entenda a formacao do anel, e que há um usuario em cada nodo ou estação
//    2) implemente o codigo do nodo, veja em comentarios abaixo
//
package main

import (
	"strconv"
	"time"
)

const N = 4
const nMsgs = 100

type Msg struct {
	sender   int
	receiver int
	message  string
}

type Packet struct {
	token bool // se true é token, e se falso  ...
	msg   Msg  // o pacote tem uma mensagem
}

func node(id int,
	hasToken bool,
	send chan Msg,
	receive chan Msg,
	ringMy chan Packet,
	ringNext chan Packet) {

	println("node ", id)
	for {

		pct := Packet{false, Msg{id, id + 1, "msg"}}

		// ------------------
		// COMPLETE O CODIGO!!
		// ------------------
		Packet{true, Msg{0, 0, "noMsg"}}
	}
}

func user(id int, send chan Msg, rec chan Msg) {
	// usuario pode mandar e receber concorrentemente
	go func() {
		for i := 0; i <= nMsgs; i++ {
			send <- Msg{id, i % N, "msg" + strconv.Itoa(i)}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	go func() {
		for {
			m := <-rec
			println("Node: ", id, " rcv from: ", m.sender, " - ", m.message)
		}
	}()
}

func main() {
	var chanRing [N]chan Packet
	var chanSend [N]chan Msg
	var chanRec [N]chan Msg

	for i := 0; i < N; i++ {
		chanRing[i] = make(chan Packet)
		chanSend[i] = make(chan Msg)
		chanRec[i] = make(chan Msg)
	}
	for i := 0; i < (N - 1); i++ {
		go node(i, false, chanSend[i], chanRec[i], chanRing[i], chanRing[i+1])
		go user(i, chanSend[i], chanRec[i]) // um usuario em cada nodo do anel
	}
	go node(N-1, true, chanSend[N-1], chanRec[N-1], chanRing[N-1], chanRing[0]) // fecha o ciclo
	go user(N-1, chanSend[N-1], chanRec[N-1])

	fin := make(chan struct{})
	<-fin
}
