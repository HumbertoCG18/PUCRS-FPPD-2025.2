// -----------------------------------------
// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// Este programa apresenta uma arcabouço para implementar o problema da rede em anel (token ring)
// >>>> voce deve ter feito os anteriores desta serie.
// PROBLEMA:
//    1) seria possivel ter mais de um token circulando na rede para aumentar a vazao da mesma ?
//       como voce implementaria isto ?
//       é possivel simplesmente iniciar mais uma estacao com um token (em true) ?
//       se não é, qual o problema ?
//       como voce modificaria  ?
// ATENCAO
//       1) o codigo abaixo NAO apresenta a solucao.  é uma cópia do anterior.
//          VOCE DEVE DESENVOLVER ESTA SOLUCAO.
//          este desenvolvimento pode ser com seu grupo.
//          deverá ser entregue em data marcada.
//       2) NOTE que como temos mais de um token, uma estação deve poder tratar pacotes de
//          outras mesmo que um pacote seu esteja dando a volta no anel.
//          entao o trecho
//					ringNext <- Packet{false, msg}
//					// espero a msg dar a volta e tiro ela do anel
//					<-ringMy
//          nao funciona neste caso.
//          Dito isto, a modificacao pode ocorrer nesta e tambem em outras partes.

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
		if hasToken {
			select {
			case msg := <-send:
				// tenho token e mando a msg
				ringNext <- Packet{false, msg}
				// espero a msg dar a volta e tiro ela do anel
				<-ringMy
				// mando o token pq já usei
				ringNext <- Packet{true, Msg{}}
				hasToken = false
			default:
				// mando o token pq n sou egoista
				ringNext <- Packet{true, Msg{}}
				hasToken = false
			}
		} else {
			pk := <-ringMy
			// se tem token n faz mais nada, se n tem checa se eh pra mim
			// e sempre passa pra frente (a msg tem que rodar)
			if pk.token {
				hasToken = true
			} else if pk.msg.receiver == id {
				receive <- pk.msg
				ringNext <- pk
			} else {
				ringNext <- pk
			}

		}
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
		go user(i, chanSend[i], chanRec[i])
	}
	// este nodo inicia com o token
	go node(N-1, true, chanSend[N-1], chanRec[N-1], chanRing[N-1], chanRing[0])
	go user(N-1, chanSend[N-1], chanRec[N-1])

	fin := make(chan struct{})
	<-fin
}
