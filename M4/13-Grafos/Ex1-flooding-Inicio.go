// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// Este ee um exemplo muito simples de como modelar uma topologia com N nodos.
// A topologia ee modelada por uma matriz de incidencia.
// Basta adicionar 1s na matriz, na posicao [i][j] para criar arestas direcionadas do nodo i para o nodo j.
// Nesta modelagem cada nodo tem um canal de entrada.
// Com a funcao broadcast um nodo manda para todos vizinhos.
// j é dito vizinho de i se existe topology[i][j]=1
// ATENÇÃO:
//    1) a topologia criada no exemplo abaixo é um **grafo aciclico dirigido.**
//       cada aresta tem somente uma direcao (ex.: 0 manda para 1, mas 1 nao manda para 0)
//       Assim, um nodo pode receber mais de uma vez uma mensagem,
//       mas nesta topologia a mensagem nao entra em ciclo.
// EXERCÍCIO:
//    1) rode o exemplo.   note que cada mensagem é repassada mais de uma vez em alguns nodos
//    2) Implemente a eliminação de duplicatas.
//       Cada nodo deve repassar a mensagem apenas a primeira vez que a recebe.

package main

import (
	"fmt"
)

// nro de nodos
const N = 10

// Topologia é uma matriz NxN onde 1 em [i][j] indica presenca da aresta do nodo i para o j
type Topology [N][N]int

// O que é enviado entre nodos, pode adicionar campos nesta estrutura ...
type Message struct {
	id int // identificador da mensagem - um nro sequencial ...
}

// um canal de entrada para cada nodo i
type inputChan [N]chan Message

type nodeStruct struct {
	id   int
	topo Topology
	inCh inputChan
}

// tamanho do buffer de cada canal de entrada
const channelBufferSize = 1

// difusão ou broadcast - um nodo manda para TODOS seus vizinhos do grafo
// nodo origin, conforme topology, usando canais do vetor inCh, manda message para todos eles
func (n *nodeStruct) broadcast(m Message) { // broadcast(origin int, topo Topology, inCh inputChan, m Message) {
	for j := 0; j < N; j++ { // para todo vizinho j em N
		if n.topo[n.id][j] == 1 { //  a matriz em [origin][j] diz se origin conectado com j
			n.inCh[j] <- m // escreve m no canal de j
		}
	}
}

//  cada nodo recebe toda matriz de conectividade e os canais de entrada de todos processos
//  cada nodo le o seu canal de entrada e escreve a mensagem em todos canais de saida
//  (dele para outros nodos usando a funcao send)
func (n *nodeStruct) nodo() { //(id int, topo Topology, inCh inputChan) {
	fmt.Println(n.id, " ativo! ")
	for {
		m := <-n.inCh[n.id]                // espera entrada entrada, reage
		fmt.Println(n.id, " tratando ", m) // avisa
		n.broadcast(m)                     // repassa m em todas arestas de saida
	}
}

// ------------------------------------------------------------------------------------------------
// no main: montagem da topologia, criacao de canais, inicializacao de nodos e geracao de mensagens
// ------------------------------------------------------------------------------------------------

func main() {
	var topo Topology
	//  se [i,j]==1, entao o nodo i pode enviar para o nodo j pelo canal j
	//  para alterar a topologia basta adicionar 1s.  cada 1 é uma aresta direcional.
	//  para modelar comunicacao em ambas direcoes entre i e j, entao [i,j] e [j,i] devem ser 1
	topo = [N][N]int{
		// conforme algoritmo na funco "nodo"
		//  0  1  2  3  4  5  6  7  8  9       aresta de    para
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
		go n.nodo() // por simplicidade todos nodos tem acesso aa mesma topologia (que nao muda)
		// assim como todo nodo tem acesso ao canal de entrada de todos outros mas vai usar somente os que pode enviar
	}

	// carga de mensagens para que sejam "inundadas" na rede
	for i := 1; i < 2; i++ { // gera mensagem de teste a cada segundo
		inCh[0] <- Message{i}
		//time.Sleep(time.Second)
	}
	<-make(chan struct{}) // bloqueia senao nodos acabam
}
