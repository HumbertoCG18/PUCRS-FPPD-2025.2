// por Fernando Dotti - PUCRS
//
// go run PipeSort2.go

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const N = 2000
const MAX = 999

func main() {

	var v [N + 2]int
	var j int
	fmt.Println("  ------ sequencial -------")
	rand.Seed(time.Now().UnixNano())
	start1 := time.Now()
	for i := 0; i < N; i++ {
		valor := rand.Intn(MAX) - rand.Intn(MAX)
		//insereOrdenado()
		for j = 0; j < i; j++ {
			if v[j] >= valor {
				break
			}
		}
		//fmt.Println(v)

		for k := i + 1; k >= j; k-- {
			v[k+1] = v[k]
		}
		v[j] = valor
	}
	t1 := time.Since(start1).Microseconds()
	fmt.Println("  ------ tempo em microsec ------->>   ", t1)
	//fmt.Println(v)

	fmt.Println("------ Pipe Sort -------")

	var result chan int = make(chan int) // canal em que a main vai ler os resultados em ordem
	var canais [N + 1]chan int           // canais do pipe de ordenadores
	for i := 0; i <= N; i++ {            // aloca canais
		canais[i] = make(chan int, 100)
	}

	// Monta pipeline com N processos concorrentes.
	for i := 0; i < N; i++ {
		go cellSorter(i, canais[i], canais[i+1], result, MAX)
	}
	// Neste ponto temos N rotinas cellSorter concorrentes a esta linha de execucao main.
	// Elas estao paradas esperando valores em seus respectivos canais "in"

	// gera valores aleatorios para o pipeline
	//fmt.Println("  ------ pipeline -------")
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	for i := 0; i < N; i++ {
		valor := rand.Intn(MAX) - rand.Intn(MAX)
		canais[0] <- valor // manda valor para a primeira cellSorter
		//	fmt.Println("   in  ", i, " ", valor)
	}
	canais[0] <- MAX + 1 // depois de mandar N valores, insere sinal de final (MAX+1 significa fim)

	// le valores dos cellSorters (note que os cellSorters escrevem em ordem em result, como isso ee garantido ?)
	//fmt.Println("  ------ resultado -------")
	for i := 0; i < N; i++ {
		<-result
		//fmt.Println("   result  ", i, " ", <-result)
	}
	t := time.Since(start).Microseconds()
	fmt.Println("  ------ tempo em microsec ------->>   ", t)
	<-canais[N] // le sinal de fim do ultimo processo
}

// ---------------------------------------------------------------------
// cellSorter

func cellSorter(i int, in chan int, out chan int, result chan int, max int) {
	var myVal int
	var undef bool = true
	for {
		n := <-in       // rotina reage a uma entrada, altera estado e gera saida
		if n == max+1 { // sinal de final de stream de numeros
			result <- myVal // devolve valor guardado
			out <- n        // passa a diante sinal de fim
			break           // para
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
