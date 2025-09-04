// MCC - Fernando Dotti
package main

import "fmt"

func proc(s string, rx chan struct{}, ry chan struct{}) {
	for {
		<-rx
		<-ry
		rx <- struct{}{}
		ry <- struct{}{}
		fmt.Print(s)
	}
}

func main() {
	r1 := make(chan struct{}, 1)
	r2 := make(chan struct{}, 1)
	r1 <- struct{}{}
	r2 <- struct{}{}
	go proc("|", r1, r2) //  proc A
	go proc("-", r2, r1) //  proc B
	var blq chan struct{} = make(chan struct{})
	<-blq
}
