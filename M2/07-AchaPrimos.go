// por Fernando Dotti - PUCRS -
// 		assim como sobrepor tempo de comunicacao pode ser vantajoso,
// 		sobrepor tempos de processamento se a máquina dispõe de diversos núcleos
// 		pode também levar a ganhos.
// PROBLEMA:
//		encontre abaixo um programa sequencial que conta o numero de primos de um array.
// 		Se os números primos forem "grandes" o calculo se eles são primos torna-se CPU intensivo
// 		e o uso de diversos nucleos mostra o ganho de desempenho.
// EXERCICIO:
// 		1) torne este programa concorrente, sobrepondo temporalmente o trabalho de
//		   computar se um valor é primo
// OBSERVACAO: garanta com que seu gerador de numeros gere valores com diversas casas.   aprox 10.
//
// SOLUCAO:  abaixo uma solucao para o problema acima.
// OBSERVACAO: veja a adaptacao da funcao isPrime para mandar resposta em um canal,
//             e nao como retorno.   Este tipo de modificacao é comum pois ao tornar um trecho de codigo
//             concorrente, não há uma espera pelo retorno do controle e do resultado juntos.
//             E o resultado deve ser passado de alguma forma.

package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

const (
	N         = 2000
	INTERVALO = 999999999999
)

func main() {
	runtime.GOMAXPROCS(8)
	//fmt.Printf("GOMAXPROCS is %d\n", runtime.GOMAXPROCS(4))

	fmt.Println("------ conta primos -------")

	slice := generateSlice(N)
	fmt.Println(slice)

	start := time.Now()
	p := contaPrimosSeq(slice)
	fmt.Println("  -> sequencial  ------ secs: ", time.Since(start).Seconds())
	fmt.Println("  ------ n primos :  ", p)

	start1 := time.Now()
	p = contaPrimosConc(slice)
	fmt.Println("  -> concorrente ------ secs: ", time.Since(start1).Seconds())
	fmt.Println("  ------ n primos :  ", p)
}

// Generates a slice of size, size filled with random numbers
func generateSlice(size int) []int {
	slice := make([]int, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(INTERVALO) // - rand.Intn(INTERVALO)
	}
	return slice
}

func contaPrimosSeq(s []int) int {
	result := 0
	for i := 0; i < N; i++ {
		if isPrime(s[i]) {
			result++
		}
	}
	return result
}

func contaPrimosConc(s []int) int {
	result := 0
	ret := make(chan bool, N)
	for i := 0; i < N; i++ {
		go isPrimeConc(s[i], ret)
	}
	for i := 0; i < N; i++ {
		if <-ret {
			result++
		}
	}
	return result
}

func isPrimeConc(p int, ret chan bool) {
	ret <- isPrime(p)
}

// Is p prime?
func isPrime(p int) bool {
	if p%2 == 0 {
		return false
	}
	for i := 3; i*i <= p; i += 2 {
		if p%i == 0 {
			return false
		}
	}
	//fmt.Print(p, "  ")
	return true
}
