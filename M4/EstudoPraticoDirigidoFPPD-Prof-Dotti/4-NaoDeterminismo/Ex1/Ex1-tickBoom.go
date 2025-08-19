// Disciplina de Modelos de Computacao Concorrente
// Escola Politecnica - PUCRS
// Prof.  Fernando Dotti
// programa da internet

// a biblioteca time de go oferece a possibilidade de criar
// canais que mandam sinais temporizados conforme parametro passado na criacao.
//
// Atenção: veja o comando nao deterministico (select case) reagindo a diversos canais
//
// Exercício:
// 		crie outros canais para eventos temporizados e declare reacoes aos
// 		mesmos junto ao mesmo select

package main

import (
	"fmt"
	"time"
)

func main() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}
