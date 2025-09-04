// PUCRS - Fernando Dotti
// Canais são úteis para construir aplicacoes com correntes (streams) de dados.
// Normalmente estas aplicações rodam indefinidamente, sem fim.
//
// Veja o programa abaixo:
//    Entao, exceto o processo gerador, que pára em N valores, os outros processos permanecem.
//    Aqui temos dois geradores que mandam para um merger, e este para um consumidor unico.
//    Note o nao determinismo do merger.
//
// EXERCICIO:
//   1)  Com estes processos, instancie outras topologias.
//       Por exemplo, faça com que o resultado de um merger va para outro merger
//       que recebe também de um terceiro gerador.
//   Observe a facilidade de utilizar processos como blocos que podem ser combinados,
//   desde que respeitada a interface: ou seja, a tipagem e o comportamento dos
//   canais de comunicacao entre eles.

package main

const N = 200
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
		println("Em ", nome, " chegou ", dado)
	}
}

func main() {
	geradorUm := make(chan int, sizeChan)
	geradorDois := make(chan int, sizeChan)
	geradorSend := make(chan int, sizeChan)
	go consumer("geral", geradorSend)
	go merger(geradorUm, geradorDois, geradorSend)
	go geraValores(geradorUm)
	go geraValores(geradorDois)

	// bloqueia para nao acabar
	<-make(chan bool)

}
