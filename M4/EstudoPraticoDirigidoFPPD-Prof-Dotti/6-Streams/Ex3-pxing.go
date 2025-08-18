// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// escreva o programa que fica fazendo ifinitamente pang-peng-ping-pong-pung-pang-peng ...
// quantas bolas podem ser colocadas neste sistema ?

package main

func pXng(s string, ...  ) {
	for {
		<-in
		println(s)
		out <- ... 
	}
}

func main() {
	...  := make(chan struct{})
	... 
	
	go pXng(" ", ...  )
	go pXng(" ", ... )
	... 
	// coloca diversas bolas no sistema
	// quantas podem ser colocadas ?

	<- make chan struct{}{}
}
