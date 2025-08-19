// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// Este é um teste da semantica quando processos bloqueiam em um canal.
//
// Suponha que voce tem muitos processos bloqueados tentando escrever em um canal.
// QUESTAO: Será que o processo que está bloqueado há mais tempo é o
//          que ganha o direito de escrever no canal antes dos demais bloqueados ?
//          ou seja, o atendimento dos bloqueados é FIFO ?
// A especificação da linguagem, tanto quanto pude avaliar, nao aborda este aspecto.
// Isto é um teste.   Não serve como *prova* de que o comportamento é este.
// EXERCICIO:
//     1) leia o programa, execute e veja o resultado.  Qual a sua conclusão ?
// OBSERVACAO:
//    Para "aumentar a probabilidade" de que os processos se bloqueiem na ordem de seu
//    identificador, entre a criação de um processo e outro, main faz sleep.
//    Veja mais comentarios na linha do sleep.

package main

import (
	"fmt"
	"time"
)

const N = 200

// uma rotina concorrente que tenta escrever em um canal sincronizante
func seBloqueiaNoCanal(i int, c chan int) {
	c <- i
}

func main() {
	var chanBloq chan int = make(chan int) // canal para bloquear as N goRotinas
	for i := 0; i <= N; i++ {
		go seBloqueiaNoCanal(i, chanBloq) // lança N goRotinas
		time.Sleep(20 * time.Microsecond) // espera um pouco :   microsec = (1/1.000.000) sec
		// para que a rotina lançada chegue ao bloqueio.  Não é garantido.  Apenas um teste.
		// Dependendo da maquina e SO, deve regular este valor para gerar 0 respostas fora de ordem
	}
	// chegando aqui, devemos ter N processos "seBloqueiaNoCanal" bloqueados no canal chanBloq
	// e com alta probabilidade resguardando a ordem de bloqueio conforme seu identificador i
	// agora vamos ver se as sincronizações acontecem mantendo a ordem de bloqueio ou não
	outOfOrder := 0 // contador de elementos fora de ordem
	v := 0          // valor lido do processo bloqueado
	for i := 0; i <= N; i++ {
		v = <-chanBloq // desbloqueia os N processos, obtendo seu id
		if i != v {    // se o id nao for i, nao está em ordem,
			fmt.Print(" .", v, " ") // avisa com um . qual id esta fora de ordem
			outOfOrder++            // computa quantos fora de ordem
		} else {
			fmt.Print(v, " ") // print o id na ordem
		}
	}
	fmt.Println(" ")
	fmt.Println(" -- ")
	fmt.Println(" --  Total de respostas fora de ordem: ", outOfOrder)
	fmt.Println(" --  Caso seja maior que 0, tente aumentar um pouco o tempo de sleep.")
	fmt.Println(" -- ")
}
