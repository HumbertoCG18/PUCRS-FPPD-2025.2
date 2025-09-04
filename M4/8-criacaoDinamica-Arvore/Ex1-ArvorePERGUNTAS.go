// por Fernando Dotti - PUCRS
// dado abaixo um exemplo de estrutura em arvore, uma arvore inicializada
// e uma operação de caminhamento, pede-se fazer:
//   1.a) a operação que soma todos elementos da arvore.
//        func soma(r *Nodo) int {...}
//   1.b) uma operação concorrente que soma todos elementos da arvore
//   [OS ACIMA ESTAO RESOLVIDOS]
//   2.a) a operação de busca de um elemento v, dizendo true se encontrou v na árvore, ou falso
//        func busca(r* Nodo, v int) bool {}...}
//   2.b) a operação de busca concorrente de um elemento, que informa imediatamente
//        por um canal se encontrou o elemento (sem acabar a busca), ou informa
//        que nao encontrou ao final da busca
//   3.a) a operação que escreve todos pares em um canal de saidaPares e
//        todos impares em um canal saidaImpares, e ao final avisa que acabou em um canal fin
//        func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}){...}
//   3.b) a versao concorrente da operação acima, ou seja, os varios nodos sao testados
//        concorrentemente se pares ou impares, escrevendo o valor no canal adequado
//
//  ABAIXO: RESPOSTAS A QUESTOES 1a e b
//  APRESENTE A SOLUÇÃO PARA AS DEMAIS QUESTÕES

package main

import (
	"fmt"
)

type Nodo struct {
	v int
	e *Nodo
	d *Nodo
}

func caminhaERD(r *Nodo) {
	if r != nil {
		caminhaERD(r.e)
		fmt.Print(r.v, ", ")
		caminhaERD(r.d)
	}
}

// -------- SOMA ----------
// soma sequencial recursiva
func soma(r *Nodo) int {
	if r != nil {
		//fmt.Print(r.v, ", ")
		return r.v + soma(r.e) + soma(r.d)
	}
	return 0
}

// funcao "wraper" retorna valor
// internamente dispara recursao com somaConcCh
// usando canais
func somaConc(r *Nodo) int {
	s := make(chan int)
	go somaConcCh(r, s)
	return <-s
}
func somaConcCh(r *Nodo, s chan int) {
	if r != nil {
		s1 := make(chan int)
		go somaConcCh(r.e, s1)
		go somaConcCh(r.d, s1)
		s <- (r.v + <-s1 + <-s1)
	} else {
		s <- 0
	}
}

// -------- BUSCA ----------
// busca sequencial recursiva
// func busca(r *Nodo, val int) bool { }

// busca concorrente recursiva
// func buscaConc(r *Nodo, val int, ret chan bool) { }

// -------- SAIDAS PAR E IMPAR --------
// Sequencial

//func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) { }

// ---------   agora vamos criar a arvore e usar as funcoes acima

func main() {
	root := &Nodo{v: 10,
		e: &Nodo{v: 5,
			e: &Nodo{v: 3,
				e: &Nodo{v: 1, e: nil, d: nil},
				d: &Nodo{v: 4, e: nil, d: nil}},
			d: &Nodo{v: 7,
				e: &Nodo{v: 6, e: nil, d: nil},
				d: &Nodo{v: 8, e: nil, d: nil}}},
		d: &Nodo{v: 15,
			e: &Nodo{v: 13,
				e: &Nodo{v: 12, e: nil, d: nil},
				d: &Nodo{v: 14, e: nil, d: nil}},
			d: &Nodo{v: 18,
				e: &Nodo{v: 17, e: nil, d: nil},
				d: &Nodo{v: 19, e: nil, d: nil}}}}

	saidaP := make(chan int)
	saidaI := make(chan int)
	fin := make(chan struct{})

	fmt.Println()
	fmt.Println()
	fmt.Println("Valores na árvore: ")
	go retornaParImpar(root, saidaP, saidaI, fin)
	fim := false
	for count := 0; count < 20 && !fim; {
		select {
		case par := <-saidaP:
			fmt.Println("Par:", par)
			count++
		case impar := <-saidaI:
			fmt.Println("Impar:", impar)
			count++
		case <-fin:
			fim = true
		}
	}

	fmt.Println()
	fmt.Print("Valores na árvore: ")
	caminhaERD(root)
	fmt.Println()
	fmt.Println()

	fmt.Println("Soma: ", soma(root))
	fmt.Println("SomaConc: ", somaConc(root))
	fmt.Println()
	fmt.Println("Busca 17: ", busca(root, 17))
	fmt.Println("Busca 99: ", busca(root, 99))

	fmt.Println()
	fmt.Println("BuscaC 17: ", buscaC(root, 17))
	fmt.Println("BuscaC 99: ", buscaC(root, 99))
}
