// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// PROBLEMA:
//   o dorminhoco especificado no arquivo Ex1-ExplanacaoDoDorminhoco.pdf nesta pasta
// ESTE ARQUIVO
//   Um template para criar um anel generico.
//   Adapte para o problema do dorminhoco.
//   Nada está dito sobre como funciona a ordem de processos que batem.
//   O ultimo leva a rolhada ...
//   ESTE  PROGRAMA NAO FUNCIONA.    É UM RASCUNHO COM DICAS.


package main

import (
	"fmt"
)

const NJ = 5           // numero de jogadores
const M = 4            // numero de cartas na mao

type carta string      // carta é um strirng

var ch [NJ]chan carta  // NJ canais de itens tipo carta  

func jogador(id int, in chan carta, out chan carta, cartasIniciais []carta, ... ) {
	mao := cartasIniciais    // estado local - as cartas na mao do jogador
	nroDeCartas := M         // quantas cartas ele tem 
    cartaRecebida := " "     // carta recebida é vazia
    estado := jogando

	for {
		if estado==jogando
			{
			// cartaRecebida = <-in     recebe carta na entrada
			//                          e processa, escreve outra na saida,
			//                          fica ou nao pronto para bater
			// OU
			// algem bate antes ?
			// 
		} else { // estado é prontoParaBater  
			// bate
			// OU
			// algem bate antes ?
		}
	}
}

func main() {
	// cria canais de passagem de cartas
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan carta)
	}

	// cria canais para bater ? 

	// baralho = cria um baralho com NJ*M cartas

	for i := 0; i < NJ; i++ {   // cria os NJ jogadores

		// cartasEscolhidas = escolhe aleatoriamente (e tira) M cartas do baralho para o jogador i

		go jogador(i, ch[i], ch[(i+1)%NJ], cartasEscolhidas , ...)    // cria jogador i conectado com i-1 e i+1, e com as cartas
	}
	
	// escolhe um jogador j e escreve uma carta em seu canal de entrada

	// espera ate jogadores baterem no(s) canal(is) de batida
	// registra ordem de batida
}


