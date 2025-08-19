// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// esta solucao utiliza workgroups de Golang.
// um wg é usado para esperar que um grupo de rotinas acabe.
// a solução é equivalente à anterior com canais.
// nesta disciplina, na parte que se refere uso de canais,
// o professor utilizará somente canais como forma de sincronizacao, e nao outras
// abstracoes da biblioteca da linguagem - todas elas podem ser reduzidas ao uso de canais.
// voce pode usar esta funcionalidade se desejar.

package main

import (
	"fmt"
	"sync"
)

func say(s string, wg *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		fmt.Println(s)
	}
	wg.Done()
}

func main() {
	var waitgroup sync.WaitGroup
	waitgroup.Add(2)
	go say("world", &waitgroup)
	go say("hello", &waitgroup)
	waitgroup.Wait()
}
