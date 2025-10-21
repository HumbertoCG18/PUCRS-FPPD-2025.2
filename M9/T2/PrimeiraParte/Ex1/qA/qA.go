package main

import (
	"fmt"
	"sync"
)

func qA(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(i)
}

func questaoA() {
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go qA(i, &wg)
	}

	wg.Wait() // espera as 10 goroutines terminarem
}

func main() {
	questaoA()
}
