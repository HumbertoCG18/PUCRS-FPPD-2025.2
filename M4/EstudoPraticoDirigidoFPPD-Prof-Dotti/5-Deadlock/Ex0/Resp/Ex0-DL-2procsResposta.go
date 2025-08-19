// PUCRS - Fernando Dotti
// Modelagem de dois processos que precisam acessar concorrentemente dois
// recursos (r1 e r2).
// Exercicios:
//     compreenda o programa
//     execute e observe o comportamento.
//     use a saida do runtime de go para identificar em que ponto cada processo para.
//     explique a razão da parada.
//
//     como voce resolveria este problema alterando uma linha de codigo apenas ?
//     (nao precisa acrescentar!!)

//     Resposta:  faça com que os processos bloqueiem recursos na mesma ordem.
//         assim, o primeiro a pegar o primeiro recurso conseguira pegar o segundo
//         e nao acontecera o bloqueio mutuo!!

package main

import "fmt"

func proc(s string, rx chan struct{}, ry chan struct{}) {
	for {
		<-rx
		<-ry
		rx <- struct{}{}
		ry <- struct{}{}
		fmt.Print(s)
	}
}

func main() {
	r1 := make(chan struct{}, 1)
	r2 := make(chan struct{}, 1)
	r1 <- struct{}{}
	r2 <- struct{}{}
	go proc("|", r1, r2) //  proc A
	go proc("-", r2, r1) //  proc B     <--- TROQUE A ORDEM DOS RECURSOS PARA SER A MESMA DO PROC A
	//for {}
	var blq chan struct{} = make(chan struct{})
	<-blq
}
