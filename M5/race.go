package main

import (
	"fmt"
)

var sharedTest int = 0

var ch_fim chan struct{} = make(chan struct{})

func MyFunc() {

	for k := 0; k < 100; k++ {
		//  ---
		sharedTest = sharedTest + 1 //  nao ee >> ATOMICO <<: nao ee indivisivel
		//  ---
	}

	ch_fim <- struct{}{} // avisa que acabou
}

func main() {
	for i := 0; i < 200; i++ { // lanca processos
		go MyFunc()
	}
	for i := 0; i < 200; i++ { // espera processos acabarem
		<-ch_fim
	}
	fmt.Println("Resultado  ", sharedTest)

	//  vai ser 200 * 100 = 20.000   ?
}
