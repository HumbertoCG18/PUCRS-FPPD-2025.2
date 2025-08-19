// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// EXERCÍCIO:  dado o programa abaixo
//     considerando-se como estados os valores da tripla x,y,z
//     qual o diagrama de estados e transicoes que representa
//     1) a questaoStSp2()  ?
//     2) a questaoStSp3()  ?
// OBS: a execucao do programa abaixo nao mostra nada.   este serve como especificacao do problema.
//      note que como não há sincronizacao, todas combinacoes possiveis de estados acontecerao.

package main

//---------------------------

var x, y, z int = 0, 0, 0

func px() {
	x = 1
	x = 2
}

func py() {
	y = 1
	y = 2
}

func pz() {
	z = 1
	z = 2
}

func questaoStSp2() {
	go px()
	py()
	for {
	}
}

func questaoStSp3() {
	go px()
	go py()
	pz()
	for {
	}
}

func main() {
	questaoStSp2()
	questaoStSp3()
}
