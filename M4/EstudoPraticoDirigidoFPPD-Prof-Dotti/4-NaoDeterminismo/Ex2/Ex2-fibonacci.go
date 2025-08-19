// Disciplina de Modelos de Computacao Concorrente
// Escola Politecnica - PUCRS
// Prof.  Fernando Dotti
// programa da internet - site google

// Exercício:
//		leia e entenda o funcionamento
// Atenção ao fato que o processo fibonacci fica à disposição para
// sincronizar com o main, e o main decide qual canal de fibonacci ele
// usa a cada momento.
// Se voce estuda calculo de processos, esta construção equivale
// à sincronização externa de CSP

package main

import "fmt"

// a cada leitura de um elemento da serie de fibonacci
// o processo mantem estado interno para avancar para o proximo
// isso acontece ate que diga para sair - quit
// entao este processo é como um servico concorrente que está sempre
// disponivel para dar o proximo numero da sequencia sempre que for solicitado
// com uma leitura em 'c' por outro processo
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
func main() {
	c := make(chan int)
	quit := make(chan int)
	// note que o trecho abaixo, ate ...}()  executa concorrentemente com a chamada da funcao fibonacci
	// o processo abaixo "dirige" as sincronizacoes do processo fibonacci.
	// ele oferece leitura de 'c' ou escrita em 'quit'  e o fibonacci apenas reage.
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
			// cada leitura acima é uma solicitacao de um valor novo
			// pois o processo fibonacci esta sempre pronto para sincronizar aqui ou em quit
		}
		quit <- 0
	}() // de go func(){ ...      até aqui é concorrente com o abaixo
	fibonacci(c, quit)
}

//   ** observe a forma de declarar um processo ativo, com "go func() { ... } ()"
//   note que valem as regras de escopo:
//   func esta definido em main, entao pode usar os canais c e quit.
//   ** observe a construção de escolha não determinística no corpo de fibonacci.
//   execute.   investigue.   pergunte!!!!!!
