// por Fernando Dotti - PUCRS
// dado abaixo um exemplo de estrutura em arvore, uma arvore inicializada
// e uma operação de caminhamento, pede-se fazer:
//   1.a) a operação que soma todos elementos da arvore.
//        func soma(r *Nodo) int {...}
//   1.b) uma operação concorrente que soma todos elementos da arvore

//  Respostas possíveis abaixo.
//

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

// -------- SOMA SEQ ----------
// soma sequencial recursiva
func soma(r *Nodo) int {
	if r != nil {
		//fmt.Print(r.v, ", ")
		return r.v + soma(r.e) + soma(r.d)
	}
	return 0
}

// -------- SOMA CONC ----------

// funcao "wraper" retorna valor
// internamente dispara recursao com somaConcCh
// usando canais
func somaConc(r *Nodo) int {
	s := make(chan int)
	go somaConcCh(r, s)
	return <-s
}
func somaConcCh(r *Nodo, s chan int) { // recursiva com canais
	if r != nil {
		s1 := make(chan int)
		go somaConcCh(r.e, s1)
		go somaConcCh(r.d, s1)
		s <- (r.v + <-s1 + <-s1)
	} else {
		s <- 0
	}
}

// -------- BUSCA SEQ ----------
// busca sequencial recursiva
func busca(r *Nodo, val int) bool {
	if r == nil {
		return false
	}
	if r.v == val {
		return true
	}
	return busca(r.e, val) || busca(r.d, val)
}

// -------- BUSCA CONC ----------
// busca concorrente usando canais
func buscaConc(r *Nodo, val int) bool {
	c := make(chan bool)
	go buscaConcCh(r, val, c)
	return <-c
}

func buscaConcCh(r *Nodo, val int, c chan bool) {
	if r == nil {
		c <- false
		return
	}
	if r.v == val {
		c <- true
		return
	}
	c1 := make(chan bool)
	c2 := make(chan bool)
	go buscaConcCh(r.e, val, c1)
	go buscaConcCh(r.d, val, c2)
	c <- (<-c1 || <-c2)
}

// -------- SEPARA PARES E IMPARES ----------
func separaParesEImpares(r *Nodo) {
	pares := []int{}
	impares := []int{}
	collectValores(r, &pares, &impares)
	fmt.Println("Valores pares:", pares)
	fmt.Println("Valores ímpares:", impares)
}

func collectValores(r *Nodo, pares *[]int, impares *[]int) {
	if r != nil {
		if r.v%2 == 0 {
			*pares = append(*pares, r.v)
		} else {
			*impares = append(*impares, r.v)
		}
		collectValores(r.e, pares, impares)
		collectValores(r.d, pares, impares)
	}
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

	fmt.Println()
	fmt.Print("Valores na árvore: ")
	caminhaERD(root)
	fmt.Println()
	fmt.Println()

	fmt.Println("Soma: ", soma(root))
	fmt.Println("SomaConc: ", somaConc(root))
	fmt.Println()

	fmt.Println("Busca 19: ", busca(root, 19))
	fmt.Println("Busca Conc 19: ", buscaConc(root, 19))
	fmt.Println("Busca Conc 20: ", buscaConc(root, 20))
	fmt.Println()

	separaParesEImpares(root)

}
