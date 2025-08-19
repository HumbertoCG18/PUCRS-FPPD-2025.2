// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// Problema:
//   resolver sort de N valores
//   abordagem: inserir cada valor em posicao correta com relacao aos demais
//              imagine inicio dos valores à esquerda.
//              cada valor eh comparado com cada outro, em ordem.
//              ao achar posicao, inserir valor e o que estava ali deve ser
//              deslocado para a direita, assim como todos os outros à direita.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const N = 200
const MAX = 999

func main() {
	var v [N + 1]int
	var j int
	fmt.Println("  ------ sequencial -------")
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < N-1; i++ {
		// gera um valor
		valor := rand.Intn(MAX) - rand.Intn(MAX)

		// acha posicao em relacao aos demais ja colocados
		for j = 0; j < i; j++ {
			if v[j] >= valor {
				break
			}
		}

		// desloca restante para a direita
		for k := i + 1; N > k && k >= j; k-- {
			v[k+1] = v[k]
		}
		v[j] = valor
	}
	fmt.Println(v)
}
