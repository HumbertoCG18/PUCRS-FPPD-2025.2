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
// Soma sequencial recursiva
func soma(r *Nodo) int {
	if r != nil {
		//fmt.Print(r.v, ", ")
		return r.v + soma(r.e) + soma(r.d)
	}
	return 0
}

// Funcao "wraper" retorna valor nternamente dispara recursao com somaConcCh usando canais.
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
// Busca sequencial recursiva:
func busca(r *Nodo, val int) bool {
	if r == nil {
		return false
	}
	if r.v == val {
		return true
	}
	return busca(r.e, val) || busca(r.d, val)
}

// Busca concorrente recursiva:
func buscaC(r *Nodo, val int) bool {
	resultado := make(chan bool, 1)
	go buscaConc(r, val, resultado)
	return <-resultado
}

func buscaConc(r *Nodo, val int, ret chan bool) {
	if r == nil {
		ret <- false
		return
	}
	if r.v == val {
		ret <- true
		return
	}

	retE := make(chan bool, 1)
	retD := make(chan bool, 1)

	go buscaConc(r.e, val, retE)
	go buscaConc(r.d, val, retD)

	// Se qualquer um encontrar, retorna true
	for i := 0; i < 2; i++ {
		select {
		case resE := <-retE:
			if resE {
				ret <- true
				return
			}
		case resD := <-retD:
			if resD {
				ret <- true
				return
			}
		}
	}
	ret <- false
}

// -------- SAIDAS PAR E IMPAR --------
// Sequencial

func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) {
	if r != nil {
		if r.v%2 == 0 {
			saidaP <- r.v
		} else {
			saidaI <- r.v
		}
		retornaParImpar(r.e, saidaP, saidaI, fin)
		retornaParImpar(r.d, saidaP, saidaI, fin)
	} else {
		fin <- struct{}{}
	}
}

// Versão concorrente de retornaParImpar
func retornaParImparConc(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) {
	if r == nil {
		fin <- struct{}{}
		return
	}

	// Processa o nó atual
	if r.v%2 == 0 {
		saidaP <- r.v
	} else {
		saidaI <- r.v
	}

	// Canais para sincronização das subárvores
	finE := make(chan struct{})
	finD := make(chan struct{})

	// Processa subárvores concorrentemente
	go retornaParImparConc(r.e, saidaP, saidaI, finE)
	go retornaParImparConc(r.d, saidaP, saidaI, finD)

	// Espera ambas as subárvores terminarem
	<-finE
	<-finD

	// Sinaliza que terminou
	fin <- struct{}{}
}

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
