// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// >>> veja o Ex1 desta serie
// ASSUNTO - Compreensão de concorrência e canais com buffer
//
// EXERCÍCIO:
//     ...
//     3) Faça uma versão que tem vários processos destino
//        que podem consumir os dados de forma não determinística.
//        Ou seja, processos diferentes podem consumir quantidades
//        diferentes de itens,  conforme sua velocidade.
//        Como você coordenaria o término dos processos depois do
//        consumo dos N valores ?
// SOLUCAO:
//        possivel abaixo.
// EXERCÍCIO:
//     4) como seria a adição de um numero maior de destinos ?

package main

const N = 100
const tamBuff = 0

func fonteDeDados(saida chan int, n int, fin chan struct{}) {
	for i := 1; i < n; i++ {
		println(i, " -> ")
		saida <- i
	}
	fin <- struct{}{}
	fin <- struct{}{}
}

func destinoDosDados(entrada chan int, fin chan struct{}) {
	for {
		select {
		case v := <-entrada:
			println("                  -> ", v)
		case <-fin:
			return
		}
	}
}

func main() {
	c := make(chan int, tamBuff)
	fin := make(chan struct{})
	go fonteDeDados(c, N, fin)
	go destinoDosDados(c, fin)
	destinoDosDados(c, fin)
}
