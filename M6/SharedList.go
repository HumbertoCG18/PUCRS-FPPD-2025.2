package main

import (
	"fmt"
)

type Nodo struct {
	v   int
	nxt *Nodo
}

// -------- PRINT da lista ----------
// sequencial recursiva

func printLista(l *Nodo) {
	fmt.Println()
	fmt.Print("Valores na lista:    ")
	printListaR(l)
	fmt.Println()
}

func printListaR(l *Nodo) {
	if l != nil {
		fmt.Print(l.v, ", ")
		printListaR(l.nxt)
	} else {
		fmt.Print(" NIL ")
	}
}

// -------- Cont lista ----------
//
func contLista(l *Nodo) {
	fmt.Println()
	fmt.Print("# Elementos:    ")
	fmt.Println(contListaR(l))
}

func contListaR(l *Nodo) int {
	if l != nil {
		return 1 + contListaR(l.nxt)
	}
	return 0
}

// -------- SOMA ----------
// soma sequencial recursiva
func soma(l *Nodo) int {
	if l != nil {
		return l.v + soma(l.nxt)
	}
	return 0
}

// -------- INSERE val ----------
// insere no inicio
func insere(l *Nodo, val int) *Nodo {
	if l != nil {
		return &Nodo{v: val, nxt: l}
	}
	return &Nodo{v: val, nxt: nil}
}

// -------- RETIRA val ----------
// retira valor
func retira(l *Nodo, val int) *Nodo {
	if l != nil {
		if l.v == val {
			return l.nxt
		} else if l.nxt != nil {
			l.nxt = retira(l.nxt, val)
		}
	}
	return l
}

// ---------   agora vamos criar a arvore e usar as funcoes acima

func main() {
	l := insere(nil, 0)
	fin := make(chan struct{}, 0)

	go func() {
		for i := 1; i < 5000; i++ {
			l = insere(l, i)
		}
		fin <- struct{}{}
	}()

	go func() {
		for i := 1; i < 5000; i++ {
			l = insere(l, i)
		}
		fin <- struct{}{}
	}()

	<-fin
	<-fin

	printLista(l)
	contLista(l)
}
