package main

import (
	"fmt"
	"sync"
	"time"
)

// Constante N - tamanho do buffer dos canais
const N = 0

// Variáveis globais canais c1 e c2
var c1, c2 chan int

// Thread t1 - produz números ímpares em c1
func t1(wg *sync.WaitGroup) {
	defer wg.Done()
	i := 1
	for {
		c1 <- i
		i = i + 2
	}
}

// Thread t2 - produz números pares em c2
func t2(wg *sync.WaitGroup) {
	defer wg.Done()
	i := 2
	for {
		c2 <- i
		i = i + 2
	}
}

// Thread t3 - usa select para ler de c1 ou c2
func t3(wg *sync.WaitGroup) {
	defer wg.Done()
	var x int
	for {
		select {
		case x = <-c1: // x recebe valor do canal c1
			fmt.Print(x, ", ")
		case x = <-c2: // x recebe valor do canal c2
			fmt.Print(x, ", ")
		}
	}
}

func main() {
	// Cria os canais com buffer de tamanho N
	c1 = make(chan int, N)
	c2 = make(chan int, N)

	// WaitGroup para aguardar as goroutines
	var wg sync.WaitGroup

	// Inicia as 3 threads (goroutines)
	wg.Add(3)
	go t1(&wg)
	go t2(&wg)
	go t3(&wg)

	// Como as threads têm loops infinitos, executamos por tempo limitado
	time.Sleep(2 * time.Second)

	fmt.Println("\n\nPrograma 2 finalizado após 2 segundos de execução")
}