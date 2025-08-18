// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// RESPOSTA PARA EX1, abaixo, desta serie.
// EXERCICIO:
//   1)  Com estes processos, instancie outras topologias.
//       Por exemplo, faça com que o resultado de um merger va para outro merger
//       que recebe também de um terceiro gerador.
//   Observe a facilidade de utilizar processos como blocos que podem ser combinados,
//   desde que respeitada a interface: ou seja, a tipagem e o comportamento dos
//   canais de comunicacao entre eles.

package main

const N = 20
const sizeChan = 10

func geraValores(cOut chan int) {
	for i := 0; i < N; i++ {
		cOut <- int(i)
	}
}

func merger(rec1 chan int, rec2 chan int, cOut chan int) {
	for {
		select {
		case dado := <-rec1:
			println("Gerador 1, dado:", dado)
			cOut <- dado
		case dado := <-rec2:
			println("Gerador 2, dado:", dado)
			cOut <- dado
		}
	}
}

func consumer(nome string, rec chan int) {
	for {
		dado := <-rec
		println("                   Em ", nome, " chegou ", dado)
	}
}

func main() {
	geradorUm := make(chan int, sizeChan)
	geradorDois := make(chan int, sizeChan)
	geradorTres := make(chan int, sizeChan)

	merger1Send := make(chan int, sizeChan)
	merger2Send := make(chan int, sizeChan)

	go consumer("geral", merger2Send)

	go merger(geradorUm, geradorDois, merger1Send)
	go merger(geradorTres, merger1Send, merger2Send)

	go geraValores(geradorUm)
	go geraValores(geradorDois)
	go geraValores(geradorTres)

	// bloqueia para nao acabar
	<-make(chan bool)

}
