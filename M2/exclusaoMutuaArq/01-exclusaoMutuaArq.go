package main

import (
	"fmt"
	"os"
)

func fileAccess(file *os.File, lock chan struct{}, turns int, ch_fim chan struct{}){
	for i:=0; i<turns; i++ {
	   lock <- struct{}{} // ocupa lock
	       _, _  = file.WriteString("|") // marca entrada no arquivo
		   _, _  = file.WriteString(".") // marca entrada no arquivo
	   <- lock            // libera lock
	}
	ch_fim <- struct{}{}
}

func main() {
	file, _ := os.OpenFile("./mxOUT.txt", os.O_CREATE|os.O_WRONLY, 0644)

	lock    := make(chan struct{},1)
	ch_fim  := make(chan struct{})

	for i := 0; i < 100; i++ {     // cria 100 processos que acessam file
		go fileAccess(file,lock,10,ch_fim) 
	}

	fmt.Println("criei")
	for i := 0; i < 100; i++ {     // espera os 100 processos acabarem
		<-ch_fim
	}
}
