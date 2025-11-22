// Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
//  Professor: Fernando Dotti  (https://fldotti.github.io/)

/*
LANCAR N PROCESSOS EM SHELL's DIFERENTES, UMA PARA CADA PROCESSO, O SEU PROPRIO ENDERECO EE O PRIMEIRO DA LISTA
go run chat.go 127.0.0.1:5001  127.0.0.1:6001    ...   // o processo na porta 5001
go run chat.go 127.0.0.1:6001  127.0.0.1:5001    ...   // o processo na porta 6001
go run chat.go ...  127.0.0.1:6001  127.0.0.1:5001     // o processo na porta ...
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "SD/BEB"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port!")
		fmt.Println("go run chatBEB.go 127.0.0.1:5001  127.0.0.1:6001   127.0.0.1:7001")
		fmt.Println("go run chatBEB.go 127.0.0.1:6001  127.0.0.1:5001   127.0.0.1:7001")
		fmt.Println("go run chatBEB.go 127.0.0.1:7001  127.0.0.1:6001   127.0.0.1:5001")
		return
	}

	var registro []string
	addresses := os.Args[1:]
	fmt.Println(addresses)

	beb := BestEffortBroadcast_Module{
		Req: make(chan BestEffortBroadcast_Req_Message),
		Ind: make(chan BestEffortBroadcast_Ind_Message)}

	//beb.Init(addresses[0])
	beb.InitD(addresses[0], false)

	// enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string

		for {
			if scanner.Scan() {
				msg = scanner.Text()
				msg += "§" + addresses[0]
			}
			req := BestEffortBroadcast_Req_Message{
				Addresses: addresses[0:],
				Message:   msg}
			beb.Req <- req // ENVIA PARA TODOS PROCESSOS ENDERECADOS NO INICIO
		}
	}()

	// receptor de broadcasts
	go func() {
		for {
			in := <-beb.Ind // RECEBE MENSAGEM DE QUALQUER PROCESSO
			message := strings.Split(in.Message, "§")
			in.From = message[1]
			registro = append(registro, in.Message)
			in.Message = message[0]

			// imprime a mensagem recebida na tela
			fmt.Printf("               Message from %v: %v\n", in.From, in.Message)
		}
	}()

	blq := make(chan int)
	<-blq
}
