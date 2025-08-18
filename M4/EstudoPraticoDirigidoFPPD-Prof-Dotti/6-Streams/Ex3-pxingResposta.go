// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// escreva o programa que fica fazendo ifinitamente pang-peng-ping-pong-pung-pang-peng ...
// quantas bolas podem ser colocadas neste sistema ?
package main

func pXing(s string, entrada chan bool, saida chan bool) {
	for {
		<-entrada
		println(s)
		saida <- true
	}
}

func main() {
	c1 := make(chan bool)
	c2 := make(chan bool)
	c3 := make(chan bool)
	go pXing("pIng", c1, c2)
	go pXing("pOng", c2, c3)
	go pXing("pUng", c3, c1)

	// coloca bolas na mesa
	// note que deve ser menor que o nro de processo - qual a razao ?
	c1 <- true
	c2 <- true
	c3 <- true
	<-make(chan struct{}) // espera para sempre processos rodarem - acabe a exec com control C
}
