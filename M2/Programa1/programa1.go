package main

import (
	"fmt"
	"sync"
	"time"
)

// Constante N - tamanho do buffer do canal
const N = 0

// Variável global canal c
var c chan int

// Thread t1 - produz números ímpares
func t1(wg *sync.WaitGroup) {
	defer wg.Done()
	i := 1
	for {
		c <- i
		i = i + 2
	}
}

// Thread t2 - produz números pares
func t2(wg *sync.WaitGroup) {
	defer wg.Done()
	i := 2
	for {
		c <- i
		i = i + 2
	}
}

// Thread t3 - consome e imprime valores do canal
func t3(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		fmt.Print(<-c, ", ")
	}
}

func main() {
	// Cria o canal com buffer de tamanho N
	c = make(chan int, N)

	// WaitGroup para aguardar as goroutines
	var wg sync.WaitGroup

	// Inicia as 3 threads (goroutines)
	wg.Add(3)
	go t1(&wg)
	go t2(&wg)
	go t3(&wg)

	// Como as threads têm loops infinitos, vamos executar por um tempo limitado
	// para demonstração (senão o programa nunca terminaria)
	time.Sleep(2 * time.Second)

	fmt.Println("\n\nPrograma finalizado após 2 segundos de execução")
	// Note: Em um cenário real com loops infinitos, você precisaria de um
	// mecanismo de cancelamento (context) para finalizar as goroutines
}
