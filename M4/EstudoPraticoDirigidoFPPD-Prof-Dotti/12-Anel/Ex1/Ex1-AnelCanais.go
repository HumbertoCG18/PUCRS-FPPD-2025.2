// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// PROBLEMA:
//   tente montar uma topologia de processos e canais formando um anel
//   onde os processos passam circularmente um sinal(token)
//   cada vez que o processo recebe o sinal, ele faz print do seu identificador
// RESPOSTA
//   Veja uma solucao abaixo.
//   Cada processo printa numa coluna diferente da tela
// EXERCICIO:
//   1)	 Experimente colocar entrada em mais de um processo no sistema.
//       Vide comentarios no main.    O que ocorre ?

package main

import (
	"fmt"
)

const N = 40

var ch [N]chan struct{}

func nodoDoAnel(id int, s string, in chan struct{}, out chan struct{}) {
	for {
		<-in               // sincroniza com anterior
		fmt.Println(s, id) //  escreve
		out <- struct{}{}  // sincroniza com posterior
	}
}

func geraEspaco(tam int) string {
	s := "  "
	for j := 0; j < tam; j++ {
		s = s + "   "
	}
	return s
}

func main() {
	for i := 0; i < N; i++ {
		ch[i] = make(chan struct{})
	}
	for i := 0; i < N; i++ {
		go nodoDoAnel(i, geraEspaco(i), ch[i], ch[(i+1)%N]) // cria processos conectados circularmente
	}
	// coloca um sinal de entrada no primeiro processo
	// o que acontece ?
	ch[0] <- struct{}{} // libera
	// e se voce colocar mais sinais de entrada ?
	// ch[20] <- struct{}{} //  <<< tente usar mais liberacoes
	//    <<< e veja o resultado
	<-make(chan struct{}) // bloqueia
}
