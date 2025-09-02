package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Constantes que definem a configuração da barbearia
const (
	shopCapacity    = 20 // Capacidade máxima de clientes na loja
	sofaCapacity    = 4  // Lugares no sofá
	standingRoomCap = 16 // Espaço para ficar em pé (Capacidade - Sofá)
	barberChairs    = 3  // Número de barbeiros / cadeiras
	totalCustomers  = 30 // Número total de clientes que tentarão visitar a barbearia na simulação
)

// Estrutura para manter o estado compartilhado da barbearia
type Barbershop struct {
	customers     int
	mutex         sync.Mutex
	standingRoom  chan struct{}
	sofa          chan struct{}
	chairs        chan struct{}
	customerReady chan struct{}
	barberReady   chan struct{}
	paymentReady  chan struct{}
	receiptReady  chan struct{}
}

// Função para o barbeiro
func barber(id int, shop *Barbershop) {
	fmt.Printf("Barbeiro %d está pronto para trabalhar.\n", id)

	for {
		// 1. Espera (dorme) até um cliente estar na cadeira e sinalizar
		fmt.Printf("Barbeiro %d está dormindo...\n", id)
		<-shop.customerReady

		// 2. Sinaliza ao cliente que está pronto para começar o corte
		fmt.Printf("Barbeiro %d acordou para atender um cliente.\n", id)
		shop.barberReady <- struct{}{}

		// 3. Corta o cabelo (simulado com um sleep)
		fmt.Printf("Barbeiro %d está cortando o cabelo...\n", id)
		time.Sleep(time.Duration(rand.Intn(100)+100) * time.Millisecond)

		// 4. Espera por um pagamento (de qualquer cliente)
		fmt.Printf("Barbeiro %d terminou o corte e está esperando para receber um pagamento.\n", id)
		<-shop.paymentReady

		// 5. Aceita o pagamento e libera o cliente com o "recibo"
		fmt.Printf("Barbeiro %d está aceitando o pagamento.\n", id)
		time.Sleep(time.Duration(rand.Intn(50)+50) * time.Millisecond) // Tempo no caixa
		shop.receiptReady <- struct{}{}
	}
}

// Função para o cliente
func customer(id int, shop *Barbershop, wg *sync.WaitGroup) {
	defer wg.Done()

	// Simula um tempo de chegada
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	fmt.Printf("-> Cliente %d chegou na barbearia.\n", id)

	// --- 1. ENTRANDO NA LOJA ---
	shop.mutex.Lock()
	if shop.customers == shopCapacity {
		// Loja cheia, cliente vai embora
		fmt.Printf("   Cliente %d viu que a loja está lotada (%d clientes) e foi embora!\n", id, shop.customers)
		shop.mutex.Unlock()
		return
	}
	shop.customers++
	fmt.Printf("   Cliente %d entrou. Total de clientes: %d.\n", id, shop.customers)
	shop.mutex.Unlock()

	// --- 2. EFEITO CASCATA: EM PÉ -> SOFÁ -> CADEIRA ---

	// Ocupa um lugar na "área em pé". Se estiver lotada, espera aqui.
	shop.standingRoom <- struct{}{}
	fmt.Printf("   Cliente %d está em pé, esperando.\n", id)

	// Ocupa um lugar no sofá. Se estiver lotado, espera aqui.
	shop.sofa <- struct{}{}
	fmt.Printf("   Cliente %d sentou no sofá.\n", id)
	<-shop.standingRoom // Libera o lugar na área em pé para o próximo cliente.

	// Pega uma das cadeiras de barbeiro. Se estiverem ocupadas, espera aqui.
	shop.chairs <- struct{}{}
	fmt.Printf("   Cliente %d sentou na cadeira do barbeiro.\n", id)
	<-shop.sofa // Libera o lugar no sofá para o próximo cliente.

	// --- 3. O CORTE DE CABELO (1º RENDEZVOUS) ---
	shop.customerReady <- struct{}{} // Sinaliza: "Barbeiro, estou pronto!"
	<-shop.barberReady               // Espera o barbeiro responder: "Ok, vamos começar."
	fmt.Printf("   Cliente %d está cortando o cabelo.\n", id)

	// --- 4. O PAGAMENTO (2º RENDEZVOUS) ---
	fmt.Printf("   Cliente %d terminou o corte e vai pagar.\n", id)
	shop.paymentReady <- struct{}{} // Sinaliza: "Caixa, estou pronto para pagar!"
	<-shop.receiptReady             // Espera o "recibo" para poder sair.

	// Libera a cadeira do barbeiro para o próximo cliente.
	// NOTA: Esta linha é uma adição lógica. O pseudocódigo original omite a liberação
	// da cadeira, o que faria o sistema travar após os 3 primeiros clientes.
	<-shop.chairs

	// --- 5. SAINDO DA LOJA ---
	shop.mutex.Lock()
	shop.customers--
	fmt.Printf("<- Cliente %d pagou e está saindo. Total de clientes: %d.\n", id, shop.customers)
	shop.mutex.Unlock()
}

func main() {
	// Inicializa a semente para números aleatórios
	rand.Seed(time.Now().UnixNano())

	// Cria o estado compartilhado da barbearia
	shop := &Barbershop{
		standingRoom:  make(chan struct{}, standingRoomCap),
		sofa:          make(chan struct{}, sofaCapacity),
		chairs:        make(chan struct{}, barberChairs),
		customerReady: make(chan struct{}), // Canal sem buffer para rendezvous
		barberReady:   make(chan struct{}), // Canal sem buffer para rendezvous
		paymentReady:  make(chan struct{}), // Canal sem buffer para rendezvous
		receiptReady:  make(chan struct{}), // Canal sem buffer para rendezvous
	}

	// Inicia as goroutines dos barbeiros
	for i := 1; i <= barberChairs; i++ {
		go barber(i, shop)
	}

	// Um WaitGroup para esperar todos os clientes terminarem antes de encerrar o programa
	var wg sync.WaitGroup

	// Inicia as goroutines dos clientes
	for i := 1; i <= totalCustomers; i++ {
		wg.Add(1)
		go customer(i, shop, &wg)
	}

	// Espera todos os clientes completarem seu ciclo
	wg.Wait()
	fmt.Println("\nTodos os clientes foram atendidos. Fechando a barbearia.")
}
