package main

import (
	"fmt"
	"net/http"
	// "time"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World - Alô mundo  - Hallo Welt - ...")
}

func help(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is help - ...")
}

func f1(w http.ResponseWriter, r *http.Request) {
	// computa algo e devovle
	for i := 1; i < 5; i++ {
		fmt.Fprintf(w, "Qualquer coisa computada aqui - ... \n")
		// time.Sleep(1 * time.Second)
	}
}

func main() {
	http.HandleFunc("/", helloWorld) // declara ativacao de helloWorld a partir de http://esteIP:PORTA/
	http.HandleFunc("/help", help)   // declara ativacao de help a partir de http://esteIP:PORTA/help
	http.HandleFunc("/f1", f1)       // declara ativacao de help a partir de http://esteIP:PORTA/f1

	http.ListenAndServe(":8080", nil)
	//  declara servidor esperando na porta 8080, nesta máquina
	//  127.0.0.1 é o endereço de loopback.  quer dizer acessando esta máquina.
	//  se voce quer acessar o servico em outra máquina, tem que substituir esta parte pelo IP da outra
	//  se browser nesta maquina acessa http://127.0.0.1:8080/
	//       então passa o pedido para helloWorld
	//  se browser nesta maquina acessa http://127.0.0.1:8080/help
	//       então passa o pedido para help
}
