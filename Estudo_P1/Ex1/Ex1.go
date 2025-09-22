package main

import (
	"fmt"
	"time"
)

func proc1() {
	for i := 0; i < 3; i++ {
		fmt.Println("P1 - passo", i)
		time.Sleep(time.Microsecond * 100)
	}
}

func proc2() {
	for i := 0; i < 3; i++ {
		fmt.Println("P2 - passo", i)
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	go proc1()
	go proc2()
	time.Sleep(time.Second)
}
