// PUCRS - Fernando Dotti
//  Um sistema tem um gerador de dados que solicita a um processo Fonte enviar dados
//  para um processo Destino.
//  A cada dado recebido, o Destino manda uma confirmação para a Fonte.
//  Note que isto é uma modelagem de um sistema onde a confirmacao seria importante.
//  Com canais não existe necessidade de confirmacao de recepcao pois nenhum dado é perdido.
//
//  Resposta:
//     o processo Fonte pode escolher entre enviar e receber um confirmacao.
//     Se o processo destino recebe vários, ele manda confirmacoes, e o fonte
//     continua fazendo envios ao invés de ler confirmações, entao ocorrera que:
//     o canal de confirmações enche, o destino bloqueia tentando escrever uma
//     confirmacao, o fonte recebe mais itens do gerador e escreve para
//     o destino, enchendo o canal, na proxima escrita, o neste canal o fonte
//     bloqueia, não le a confirmacao do destino.
//
//  Novo Exercício:
//     faça uma solução que garanta que no máximo duas confirmacoes estao
//     aguardando leitura pela fonte.
//     nao adianta modificar tamanho de canal ...

package main

import (
	"fmt"
)

const (
	N = 1000
	T = 5
)

var solicitaEnvio = make(chan int, T)
var envia = make(chan int, T)
var confirma = make(chan struct{}, T)

func Gerador() {
	for i := 1; i < N; i++ {
		solicitaEnvio <- i
	}
}
func Fonte() {
	contConf := 0
	for {
		// select {
		// case x := <-solicitaEnvio:
		// 	envia <- x
		// case <-confirma:
		envia <- <-solicitaEnvio // leitura de solicitaEnvio é escrita em envia
		<-confirma
		contConf++
	}
}

func Destino() {
	for {
		rec := <-envia         // recebe valor
		confirma <- struct{}{} // confirma
		fmt.Print(rec, ", ")
	}
}

func main() {
	go Gerador()
	go Fonte()
	fmt.Println()
	go Destino()
	<-make(chan struct{}, 0)
}
