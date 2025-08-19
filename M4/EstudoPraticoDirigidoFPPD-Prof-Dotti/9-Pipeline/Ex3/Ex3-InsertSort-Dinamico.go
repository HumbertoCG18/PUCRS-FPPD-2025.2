// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// sort com pipeline de processos
// voce deve ter feito os exercicios anteriores desta serie
// PROBLEMA:
//        Como voce faria para que este pipe fosse dinâmico ?
//        ou seja, ele se ajusta a qualquer N que seja inserido.
// ATENCAO
//        o codigo abaixo NAO apresenta a solucao.  é uma cópia do anterior.
//        VOCE DEVE DESENVOLVER ESTA SOLUCAO.
//        este desenvolvimento pode ser com seu grupo.
//        deverá ser entregue em data marcada.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const N = 200
const MAX = 999

func main() {
	fmt.Println("------ Pipe Sort -------")
	// Define conjunto de canais
	var result chan int = make(chan int, N) // canal em que a main vai ler os resultados em ordem
	var canais [N + 1]chan int              // canais do pipe de ordenadores
	for i := 0; i <= N; i++ {               // aloca 0..N = N+1 canais.  canal N será lido pelo main
		canais[i] = make(chan int, 2)
	}
	// Monta pipeline com N processos concorrentes e seus canais.
	for i := 0; i < N; i++ {
		go cellSorter(i, canais[i], canais[i+1], result, MAX)
	}
	// Neste ponto temos N rotinas cellSorter concorrentes a esta linha de execucao main.
	// Elas estao paradas esperando valores em seus respectivos canais "in"
	// Agora vamos gerar valores aleatorios para o pipeline
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < N; i++ {
		valor := rand.Intn(MAX) - rand.Intn(MAX)
		canais[0] <- valor // manda valor para a primeira cellSorter
		fmt.Println("Entra ", i, " ", valor)
	}
	canais[0] <- MAX + 1 // depois de N valores, insere sinal de FIM (MAX+1 significa fim)
	// agora le valores dos cellSorters (note que os cellSorters escrevem em ordem em result, entenda como)
	for i := 0; i < N; i++ {
		v := <-result
		fmt.Println("   result  ", i, " ", v)
	}
	<-canais[N] // le sinal de fim do ultimo processo
}

// ---------------------------------------------------------------------
// processo cellSorter
func cellSorter(i int, in chan int, out chan int, result chan int, max int) {
	var myVal int
	var undef bool = true
	for {
		n := <-in // rotina bloqueia aa espera de uma entrada (modelo reativo)
		// so passa se um valor for lido, depois altera estado e gera saida
		if n == max+1 { // sinal de final de stream de numeros
			result <- myVal // devolve valor guardado
			out <- n        // passa a diante sinal de fim
			break           // PARA
		}
		if undef { // se primeiro valor
			myVal = n // guarda
			undef = false
		} else if n >= myVal { // se valor maior ou igual a este passa adiante senao fica
			out <- n
		} else {
			out <- myVal
			myVal = n
		}
	}
}
