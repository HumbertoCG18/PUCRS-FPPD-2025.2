package main

import (
	"fmt"
	"math/rand"
)

const (
	NCL  = 100
	Pool = 10
)

type Request struct {
	v      int
	ch_ret chan int
}

// ------------------------------------ Cliente
func cliente(i int, req chan Request) {
	var v, r int
	my_ch := make(chan int)
	for {
		v = rand.Intn(1000)
		req <- Request{v, my_ch}
		r = <-my_ch
		fmt.Println("cli: ", i, " req: ", v, "  resp:", r)
	}
}

// ------------------------------------ Servidor
// thread de servico calcula a resposta e manda direto pelo canal de retorno informado pelo cliente
func trataReq(id int, req Request) {
	fmt.Println("                                 trataReq ", id)
	req.ch_ret <- req.v * 2
}

// servidor que dispara threads de servico
func servidorConc(in chan Request) {
	// servidor fica em loop eterno recebendo pedidos e criando um processo concorrente para tratar cada pedido
	var j int = 0
	for {
		j++
		req := <-in
		go trataReq(j, req)
	}
}

// ------------------------------------
// SOLUÇÃO 1: Pool de Workers (limitado a 10)
func worker(id int, jobs <-chan Request) {
	for req := range jobs {
		fmt.Println("worker ", id, " tratando req")
		req.ch_ret <- req.v * 2
	}
}

func servidorPool(in chan Request) {
	// Canal para distribuir jobs para os workers
	jobs := make(chan Request, 100)

	// Cria pool de workers fixo
	for w := 1; w <= Pool; w++ {
		go worker(w, jobs)
	}

	// Distribui requisições para os workers
	for {
		req := <-in
		jobs <- req
	}
}

// ------------------------------------
// SOLUÇÃO 2: Semáforo (limitado a 10)
func trataReqComSemaforo(id int, req Request, sem chan struct{}) {
	sem <- struct{}{}        // Adquire permissão
	defer func() { <-sem }() // Libera permissão ao final

	fmt.Println("                                 trataReq ", id)
	req.ch_ret <- req.v * 2
}

func servidorComSemaforo(in chan Request) {
	// Semáforo para limitar concorrência
	sem := make(chan struct{}, Pool)

	var j int = 0
	for {
		j++
		req := <-in
		go trataReqComSemaforo(j, req, sem)
	}
}

// ------------------------------------
// main
func main() {
	fmt.Println("------ Servidores - criacao dinamica -------")
	serv_chan := make(chan Request) // CANAL POR ONDE SERVIDOR RECEBE PEDIDOS
	go servidorConc(serv_chan)      // LANÇA PROCESSO SERVIDOR
	for i := 0; i < NCL; i++ {      // LANÇA DIVERSOS CLIENTES
		go cliente(i, serv_chan)
	}
	<-make(chan int)
}
