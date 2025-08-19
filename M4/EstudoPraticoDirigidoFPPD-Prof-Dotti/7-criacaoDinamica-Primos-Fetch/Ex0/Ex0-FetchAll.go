// PUCRS - Fernando Dotti
// Fetchall fetches URLs in parallel and reports their times and sizes.
// Atencao: programa baixado da internet
// A sobreposicao do tempo de transmissao de diferentes conteúdos pela rede é importante pois
// permite que os mesmos sejam baixados em paralelo e assim disponibilizados logo ao usuário.
// Veja como baixar conteudos de URLs concorrentemente
//
// EXERCICIO:
// 		faça o fetch de cada endereco em sequencia e veja a diferença
// 		para a versao concorrente.   Execute varias vezes cada uma (aprox 5).
// 		Faça uma media dos valores sequenciais e concorrentes
// 		e veja se há ganho com a concorrencia (tempo conc / tempo seq)

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now() // START marcacao de tempo
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // baixa conteudo desta url em um a gorotina concorrente, avisa termino em ch
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // espera em ch todas rotinas acabarem
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds()) // tempo total que passou desde o START
}

func fetch(url string, ch chan<- string) {
	start := time.Now()        // este é outro start, para esta url
	resp, err := http.Get(url) // busca o conteúdo da url
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds() // fecha contagem de tempo
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

// go run Ex1-FetchAll.go https://golang.org http://gopl.io https://godoc.org
// exemplo de saida
// 0.14s     6852  https://godoc.org
// 0.16s     7261  https://golang.org
// 0.48s     2475  http://gopl.io
// 0.48s elapsed
