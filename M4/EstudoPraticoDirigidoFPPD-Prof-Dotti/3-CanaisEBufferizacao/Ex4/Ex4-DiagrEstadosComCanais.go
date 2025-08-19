// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// ASSUNTO - Compreensão de concorrência e sincronizacao no espaco de estados
// EXERCÍCIOS:
//      1) considerando-se como estados os valores da dupla x,y
//         qual o diagrama de estados e transicoes que representa questaoStSp3()  ?
//      2) e qual representa questaoStSp4()  ?
// OBSERVACAO:
//      em canais de tamanho 0, a leitura e escrita sincronizam, ou seja, formam
//      uma mesma transicao do ponto de vista do diagrama de estados.
//      no caso de uma leitura sincronizada ser atribuida a uma variavel,
//      considere que a atribuicao e a sincronizacao de leitura sao um unico passo
//      Ex.:    x = <- c
//              a sincronizacao da leitura em c e
//              a atribuicao do valor lido a x são uma unica transicao
//
package main

var x, y int = 0, 0

func pxs(c chan int) {
	x = <-c
	x = 2
}
func pys(c chan int) {
	y = 1
	c <- y
	y = 2
}
func questaoStSp3() {
	c := make(chan int, 0)
	go pxs(c)
	go pys(c)
}

//---------------------------

func pxs2(c chan int) {
	x = 1
	<-c
	x = 2
}
func pys2(c chan int) {
	y = 1
	c <- y
	y = 2
}
func questaoStSp4() {
	c := make(chan int, 0)
	go pxs2(c)
	go pys2(c)
}

//---------------------------
//---------------------------

func main() {
	questaoStSp3()
	questaoStSp4()
	for {
	}
}
