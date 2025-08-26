package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NJ = 5 // Numero de jogadores
const M = 4  // Numero de cartas na mao

type carta string // carta é um string
type estado int

const (
	jogando = iota
	prontoParaBater
	jaBateu
)

var ch [NJ]chan carta     // Canais para passar cartas
var chBatida chan int     // Canal para registrar batidas
var chParar [NJ]chan bool // Canais para parar jogadores
var ordemBatida []int     // Ordem dos jogadores que bateram

func criarBaralho() []carta {
	// Cria baralho com NJ tipos de cartas, M+1 cartas de cada tipo
	var baralho []carta
	tipos := []string{"A", "B", "C", "D", "E", "F", "G", "H"}

	for i := 0; i < NJ; i++ {
		for j := 0; j < M+1; j++ {
			baralho = append(baralho, carta(tipos[i]))
		}
	}

	// Embaralha
	rand.Seed(time.Now().UnixNano())
	for i := len(baralho) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		baralho[i], baralho[j] = baralho[j], baralho[i]
	}

	return baralho
}

func escolherCartas(baralho *[]carta, n int) []carta {
	cartas := make([]carta, n)
	for i := 0; i < n; i++ {
		cartas[i] = (*baralho)[0]
		*baralho = (*baralho)[1:]
	}
	return cartas
}

func temQuatroIguais(mao []carta) bool {
	contador := make(map[carta]int)
	for _, c := range mao {
		contador[c]++
		if contador[c] == 4 {
			return true
		}
	}
	return false
}

func jogador(id int, in chan carta, out chan carta, cartasIniciais []carta) {
	mao := make([]carta, len(cartasIniciais))
	copy(mao, cartasIniciais)
	estadoAtual := jogando

	fmt.Printf("Jogador %d iniciado com cartas: %v\n", id, mao)

	for {
		select {
		case <-chParar[id]:
			fmt.Printf("Jogador %d parou\n", id)
			return

		default:
			if estadoAtual == jogando {
				// Verifica se já tem 4 cartas iguais
				if temQuatroIguais(mao) {
					estadoAtual = prontoParaBater
					fmt.Printf("Jogador %d está pronto para bater! Mão: %v\n", id, mao)

					// Tenta bater
					select {
					case chBatida <- id:
						estadoAtual = jaBateu
						fmt.Printf("Jogador %d BATEU!\n", id)
						continue
					case <-chParar[id]:
						return
					default:
						// Não conseguiu bater ainda, continua jogando
					}
				}

				// Recebe carta (com timeout)
				select {
				case cartaRecebida := <-in:
					// Escolhe carta aleatória para passar
					indiceParaPassar := rand.Intn(len(mao))
					cartaParaPassar := mao[indiceParaPassar]

					// Substitui carta na mão
					mao[indiceParaPassar] = cartaRecebida

					// Passa carta para próximo jogador
					select {
					case out <- cartaParaPassar:
					case <-chParar[id]:
						return
					case <-time.After(100 * time.Millisecond):
						// Timeout para evitar deadlock
					}

				case <-chParar[id]:
					return
				case <-time.After(10 * time.Millisecond):
					// Timeout para não bloquear indefinidamente
				}

			} else if estadoAtual == prontoParaBater {
				// Tenta bater
				select {
				case chBatida <- id:
					estadoAtual = jaBateu
					fmt.Printf("Jogador %d BATEU!\n", id)
				case <-chParar[id]:
					return
				default:
					// Continua tentando bater
				}
			} else {
				// Já bateu, só espera o jogo terminar
				select {
				case <-chParar[id]:
					return
				case <-time.After(10 * time.Millisecond):
				}
			}
		}
	}
}

func main() {
	fmt.Println("=== JOGO DO DORMINHOCO ===")

	// Inicializa canais
	chBatida = make(chan int, NJ)
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan carta, 10) // Buffer para evitar deadlocks
		chParar[i] = make(chan bool)
	}

	// Cria e distribui cartas
	baralho := criarBaralho()
	fmt.Printf("Baralho criado com %d cartas\n", len(baralho))

	// Inicia jogadores
	for i := 0; i < NJ; i++ {
		cartasIniciais := escolherCartas(&baralho, M)
		go jogador(i, ch[i], ch[(i+1)%NJ], cartasIniciais)
	}

	// Inicia o jogo dando a primeira carta para um jogador aleatório
	if len(baralho) > 0 {
		jogadorInicial := rand.Intn(NJ)
		primeiraCartaExtra := escolherCartas(&baralho, 1)[0]
		fmt.Printf("Dando carta %s para jogador %d para iniciar\n", primeiraCartaExtra, jogadorInicial)
		ch[jogadorInicial] <- primeiraCartaExtra
	}

	// Aguarda batidas dos jogadores
	fmt.Println("Aguardando batidas...")
	ordemBatida = make([]int, 0)

	// Declara variáveis antes do loop para evitar problema com goto
	var dorminhoco int
	jogadoresQueBateram := make(map[int]bool)

	// Coleta batidas até que NJ-1 jogadores tenham batido
	for len(ordemBatida) < NJ-1 {
		select {
		case jogadorQueBateu := <-chBatida:
			ordemBatida = append(ordemBatida, jogadorQueBateu)
			fmt.Printf("Batida #%d: Jogador %d\n", len(ordemBatida), jogadorQueBateu)
		case <-time.After(5 * time.Second):
			fmt.Println("Timeout - Encerrando jogo")
			goto fim
		}
	}

	// Determina o dorminhoco (quem não bateu)
	for _, j := range ordemBatida {
		jogadoresQueBateram[j] = true
	}

	for i := 0; i < NJ; i++ {
		if !jogadoresQueBateram[i] {
			dorminhoco = i
			break
		}
	}

fim:
	// Para todos os jogadores
	for i := 0; i < NJ; i++ {
		select {
		case chParar[i] <- true:
		default:
		}
	}

	// Mostra resultados
	fmt.Println("\n=== RESULTADO ===")
	fmt.Println("Ordem de batidas:")
	for i, j := range ordemBatida {
		fmt.Printf("%d° lugar: Jogador %d\n", i+1, j)
	}
	if len(ordemBatida) == NJ-1 {
		fmt.Printf("DORMINHOCO: Jogador %d (perdeu!)\n", dorminhoco)
	}

	time.Sleep(100 * time.Millisecond) // Aguarda goroutines terminarem
}
