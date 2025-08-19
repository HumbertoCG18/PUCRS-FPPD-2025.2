// PUCRS - Fernando Dotti
//  Um sistema tem um gerador de dados que solicita a um processo Fonte enviar dados
//  para um processo Destino.
//  A cada dado recebido, o Destino manda uma confirmação para a Fonte.
//  Note que isto é uma modelagem de um sistema onde a confirmacao seria importante.
//  Com canais não existe necessidade de confirmacao de recepcao pois nenhum dado é perdido.
//
//  Exercício:
//  Rode este sistema e veja se algum problema ocorre.
//  Corrija o problema.

package main

import (
	"fmt"
)

const (
	N = 100
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
		select {
		case x := <-solicitaEnvio:
			envia <- x
		case <-confirma:
			contConf++
		}
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
