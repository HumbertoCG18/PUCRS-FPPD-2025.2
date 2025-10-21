package main

import (
	"fmt"
	"runtime"
	"time"
)

func qB(i int, c chan struct{}) {
	fmt.Println(i)
	c <- struct{}{} // envia sinal de "terminei"
}

func questaoB() {
	c := make(chan struct{}) // canal sem buffer (síncrono): envio bloqueia até alguém receber

	for i := 0; i < 10; i++ {
		go qB(i, c)
	}

	time.Sleep(20 * time.Millisecond)

	fmt.Printf("[DEBUG] Q.B.1 -> NumGoroutine (aprox): %d\n", runtime.NumGoroutine())

	for i := 0; i < 10; i++ {
		<-c
	}

	time.Sleep(10 * time.Millisecond)

	fmt.Printf("[DEBUG] Q.B.2 -> NumGoroutine (aprox): %d\n", runtime.NumGoroutine())
}

func main() {
	questaoB()
}
