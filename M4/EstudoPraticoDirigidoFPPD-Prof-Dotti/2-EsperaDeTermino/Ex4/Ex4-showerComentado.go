// Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// Exemplo da Internet
// EXERCICIOS:
//   1) rode o programa abaixo e interprete.
//      todos os valores escritos no canal são lidos?
//
//      veja que ch tem buffer 5.  assim que o for do main escrever o ultimo valor, o main segue para acabar
//      mas ainda podem existir itens ainda nao lidos no canal ch quando o main acaba (acabando o programa)

//   2) como isto poderia ser resolvido ?
//
//     uma possibilidade seria ter um canal de tamanho 0.   ainda assim nao garantiria o print do ultimo valor
//     pois main pode acabar antes de shower chegar ao print posterior aa leitura.
//
//     entao existe a necessidade de sincronizacao explicita do final dos processos.
//

package main

import "fmt"

func main() {
	ch := make(chan int, 5)
	go shower(ch)
	for i := 0; i < 10; i++ {
		ch <- i
	}
}

func shower(c chan int) {
	for {
		j := <-c
		fmt.Printf("%d\n", j)
	}
}
