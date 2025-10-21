package qd

import (
	"fmt"
)

func gera(c chan<- string, s string) {
	for {
		c <- s
	}
}

func questaoD() {
	c := make(chan string)
	go gera(c, "a")
	go gera(c, "b")
	for {
		fmt.Print(<-c)
	}
}

func main() {
	questaoD()
}
