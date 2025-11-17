//go:build dimexf
// +build dimexf

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"SD/DIMEX"
)

// waitForPeers tenta conectar a cada endereço (exceto self) repetidamente
// até que todos respondam ao Dial (ou até timeout total).
func waitForPeers(selfIdx int, addrs []string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for {
		allUp := true
		for i, a := range addrs {
			if i == selfIdx {
				continue
			}
			conn, err := net.DialTimeout("tcp", a, 800*time.Millisecond)
			if err != nil {
				allUp = false
				// não precisa logar tudo, mas loga uma vez por ciclo
				//fmt.Println("peer not ready:", a, err)
				break
			}
			_ = conn.Close()
		}
		if allUp {
			return
		}
		if time.Now().After(deadline) {
			// timeout: ainda assim retorna (evita loop infinito). Caller decide.
			fmt.Println("waitForPeers: timeout reached — proceeding anyway")
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify id and addresses!")
		return
	}

	id, _ := strconv.Atoi(os.Args[1])
	addresses := os.Args[2:]

	dmx := DIMEX.NewDIMEX(addresses, id, true)
	fmt.Println(dmx)

	// captura CTRL+C para tentar sair limpo
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		fmt.Println("\nSIGINT recebido — encerrando com limpeza...")
		// aqui podemos fazer limpeza adicional se necessário
		os.Exit(0)
	}()

	// espera curta para dar tempo do listener interno do PP2PLink subir
	// (o NewDIMEX já cria e starta o PP2PLink), mas aguardamos também os peers:
	fmt.Println("Aguardando peers ficarem disponíveis (timeout 10s)...")
	waitForPeers(id, addresses, 10*time.Second)
	// espera extra
	time.Sleep(500 * time.Millisecond)

	// loop principal: solicita e usa MX, escreve no arquivo, libera
	for {
		fmt.Println("[ APP id:", id, " PEDINDO MX ]")
		dmx.Req <- DIMEX.ENTER

		// aguarda autorização do DIMEX
		<-dmx.Ind
		fmt.Println("[ APP id:", id, " ENTROU MX ]")

		// tenta abrir/escrever no arquivo mas NÃO encerra o programa em caso de erro
		func() {
			file, err := os.OpenFile("./mxOUT.txt",
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Erro abrindo arquivo (não fatal):", err)
				return
			}
			defer func() {
				_ = file.Sync()
				_ = file.Close()
			}()

			_, err = file.WriteString("|.")
			if err != nil {
				fmt.Println("Erro escrevendo no arquivo (não fatal):", err)
				return
			}
		}()

		// sinaliza saída da seção crítica
		dmx.Req <- DIMEX.EXIT
		fmt.Println("[ APP id:", id, " SAIU MX ]")

		// aguarda um pouco antes de pedir de novo (reduz bursts)
		time.Sleep(500 * time.Millisecond)
	}
}
