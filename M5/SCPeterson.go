// Utilizado por Fernando Dotti
// Filename: SCPeterson.go
// Use below command to run:
// go run SCPeterson.go

// Leia o programa, procure entender como o algoritmo de Peterson estaa mapeado.
// O protocolo de entrada estaa em lock.   O de saida em unlock.
//
// 1) a linha 35 abaixo apresenta uma atribuicao atomica.
//    substitua por uma atribuicao comum e verifique o resultado
// 2) implemente as operações lock e unlock com canais

package main

import (
	"fmt"
	"sync/atomic"
)

// -----------------------------------------------------------------------------------------------
// ------- variaveis e procedimentos do algoritmo de peterson para exclusao mutua por duas threads
// -----------------------------------------------------------------------------------------------

var flag [2]int // variaveis do algoritmo de peterson
var turn uint32

func lock_init() {
	flag[0] = 0
	flag[1] = 0 // inicia lock resetando o desejo de ambas threads de adquirirem o lock.
	turn = 0    // daa a vez a uma delas
}

func lock(self uint32) { // executado antes da secao critica
	flag[self] = 1                    // flag[self] = 1 diz que quer o lock
	atomic.StoreUint32(&turn, 1-self) // mas antes daa aa outra trhead a chance de adquirir o lock
	for flag[1-self] == 1 && turn == 1-self {
	}
}

func unlock(self uint32) { // executado depois da secao critica
	flag[self] = 0
}

// -----------------------------------------------------------------------------------------------
// ------- exemplo de uso do algoritmo de peterson em dois processos
// -----------------------------------------------------------------------------------------------

const MAX int = 2000000

var ans int = 0 // a variavel compartilhada!!   a ser protegida

func processo(self uint32, fin chan int) {
	fmt.Println("Processo entrou: ", self) // diz qual o identificador deste processo:   0 ou 1

	for i := 0; i < MAX; i++ { // entra e sai MAX vezes na SC
		lock(self)   // CODIGO DE ENTRADA NA SC
		ans++        // SECAO CRITICA
		unlock(self) // CODIGO DE SAIDA DA SECAO CRITICA
	}
	fin <- 1
}

func main() {
	lock_init() // inicia as variaves de peterson
	fin := make(chan int)
	go processo(0, fin) // cria os processos
	go processo(1, fin)
	<-fin // espera fim de ambos
	<-fin
	fmt.Println("Valor do contador: ", ans, " | Valor esperado: ", MAX*2)
}
