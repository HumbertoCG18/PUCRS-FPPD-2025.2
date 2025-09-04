package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	shopCapacity   = 20
	sofaCapacity   = 4
	standingCap    = 16
	barberChairs   = 3
	totalCustomers = 30
)

type Barbershop struct {
	shopLimit     chan struct{} // limite total da loja
	standingRoom  chan struct{} // clientes em pé
	sofa          chan struct{} // sofá
	chairs        chan struct{} // cadeiras de barbeiro
	customerReady chan int      // cliente pronto para cortar
	barberReady   chan int      // barbeiro confirma
	payment       chan int      // cliente pronto para pagar
	receipt       chan int      // barbeiro confirma pagamento
}

func barber(id int, shop *Barbershop) {
	for {
		// espera cliente
		cid := <-shop.customerReady
		shop.barberReady <- cid

		// corta cabelo
		time.Sleep(time.Duration(rand.Intn(100)+100) * time.Millisecond)

		// espera pagamento
		cid = <-shop.payment
		time.Sleep(time.Duration(rand.Intn(50)+50) * time.Millisecond) // caixa
		shop.receipt <- cid
	}
}

func customer(id int, shop *Barbershop, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

	// tenta entrar na loja
	select {
	case shop.shopLimit <- struct{}{}:
		// entrou
	default:
		fmt.Printf("Cliente %d: loja cheia, indo embora.\n", id)
		return
	}

	// espera em pé → sofá → cadeira
	shop.standingRoom <- struct{}{}
	shop.sofa <- struct{}{}
	<-shop.standingRoom
	shop.chairs <- struct{}{}
	<-shop.sofa

	// corte
	shop.customerReady <- id
	<-shop.barberReady

	// pagamento
	shop.payment <- id
	<-shop.receipt

	<-shop.chairs    // libera cadeira
	<-shop.shopLimit // sai da loja

	fmt.Printf("Cliente %d saiu.\n", id)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	shop := &Barbershop{
		shopLimit:     make(chan struct{}, shopCapacity),
		standingRoom:  make(chan struct{}, standingCap),
		sofa:          make(chan struct{}, sofaCapacity),
		chairs:        make(chan struct{}, barberChairs),
		customerReady: make(chan int),
		barberReady:   make(chan int),
		payment:       make(chan int),
		receipt:       make(chan int),
	}

	// inicia barbeiros
	for i := 1; i <= barberChairs; i++ {
		go barber(i, shop)
	}

	var wg sync.WaitGroup
	for i := 1; i <= totalCustomers; i++ {
		wg.Add(1)
		go customer(i, shop, &wg)
	}

	wg.Wait()
	fmt.Println("Todos os clientes foram atendidos.")
}
