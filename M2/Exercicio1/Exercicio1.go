package main

import (
	"fmt"
	"os"
)

func fileAcess(file *os.File, lock chan struct{}, turns int, ch_fim chan struct{}) {
	for i := 0; i < turns; i++ {
		lock <- struct{}{}
		_, _ = file.WriteString("|")
		_, _ = file.WriteString(".")
		<-lock
	}
	ch_fim <- struct{}{}
}

func main() {
	file, _ := os.OpenFile("./mxOut.txt", os.O_CREATE|os.O_WRONLY, 0644)

	lock := make(chan struct{}, 1)
	ch_fim := make(chan struct{})

	for i := 0; i < 100; i++ {
		go fileAcess(file, lock, 10, ch_fim)
	}

	fmt.Println("Criei")
	for i := 0; i < 100; i++ {
		<-ch_fim
	}
}
