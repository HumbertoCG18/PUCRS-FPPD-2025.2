/*
  Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
  Professor: Fernando Dotti  (https://fldotti.github.io/)
  Modulo representando Berst Effort Broadcast tal como definido em:
    Introduction to Reliable and Secure Distributed Programming
    Christian Cachin, Rachid Gerraoui, Luis Rodrigues
  * Semestre 2018/2 - Primeira versao.  Estudantes:  Andre Antonitsch e Rafael Copstein
  Para uso vide ao final do arquivo, ou aplicacao chat.go que usa este
*/
package BestEffortBroadcast

import (
	"fmt"

	PP2PLink "SD/PP2PLink"
)

type BestEffortBroadcast_Req_Message struct {
	Addresses []string
	Message   string
}

type BestEffortBroadcast_Ind_Message struct {
	From    string
	Message string
}

type BestEffortBroadcast_Module struct {
	Ind      chan BestEffortBroadcast_Ind_Message
	Req      chan BestEffortBroadcast_Req_Message
	Pp2plink PP2PLink.PP2PLink
	dbg      bool
}

func (module *BestEffortBroadcast_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . [ BEB msg : " + s + " ]")
	}
}

func (module *BestEffortBroadcast_Module) Init(address string) {
	module.InitD(address, true)
}

func (module *BestEffortBroadcast_Module) InitD(address string, _dbg bool) {
	module.dbg = _dbg
	module.outDbg("Init BEB!")
	module.Pp2plink = PP2PLink.PP2PLink{
		Req: make(chan PP2PLink.PP2PLink_Req_Message),
		Ind: make(chan PP2PLink.PP2PLink_Ind_Message)}
	module.Pp2plink.InitD(address, _dbg)
	module.Start()
}

func (module *BestEffortBroadcast_Module) Start() {

	go func() {
		for {
			select {
			case y := <-module.Req:
				module.Broadcast(y)
			case y := <-module.Pp2plink.Ind:
				module.Deliver(PP2PLink2BEB(y))
			}
		}
	}()
}

func (module *BestEffortBroadcast_Module) Broadcast(message BestEffortBroadcast_Req_Message) {

	// aqui acontece o envio um opara um, para cada processo destinatario
	// em caso de injecao de falha no originador, no meio de um broadcast
	// este loop deve ser interrompido, tendo a mensagem ido para alguns mas nao para todos processos

	for i := 0; i < len(message.Addresses); i++ {
		msg := BEB2PP2PLink(message)
		msg.To = message.Addresses[i]
		module.Pp2plink.Req <- msg
		module.outDbg("Sent to " + message.Addresses[i])
	}
}

func (module *BestEffortBroadcast_Module) Deliver(message BestEffortBroadcast_Ind_Message) {

	// fmt.Println("Received '" + message.Message + "' from " + message.From)
	module.Ind <- message
	// fmt.Println("# End BEB Received")
}

func BEB2PP2PLink(message BestEffortBroadcast_Req_Message) PP2PLink.PP2PLink_Req_Message {

	return PP2PLink.PP2PLink_Req_Message{
		To:      message.Addresses[0],
		Message: message.Message}

}

func PP2PLink2BEB(message PP2PLink.PP2PLink_Ind_Message) BestEffortBroadcast_Ind_Message {

	return BestEffortBroadcast_Ind_Message{
		From:    message.From,
		Message: message.Message}
}

/*
func main() {

	if (len(os.Args) < 2) {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	addresses := os.Args[1:]
	fmt.Println(addresses)

	mod := BestEffortBroadcast_Module{
		Req: make(chan BestEffortBroadcast_Req_Message),
		Ind: make(chan BestEffortBroadcast_Ind_Message) }
	mod.Init(addresses[0])

	msg := BestEffortBroadcast_Req_Message{
		Addresses: addresses,
		Message: "BATATA!" }

	yy := make(chan string)
	mod.Req <- msg
	<- yy
}
*/
