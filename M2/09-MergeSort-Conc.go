// Merge Sort in Golang
// extratos baixados da internet
// modificacoes feitas por Fernando Dotti - PUCRS
//
// aqui encontram-se 4 implementacoes de mergeSort
// duas sequenciais e duas concorrentes
// o programa avalia o tempo de execucao de cada uma
// go run MergeSort-1.go

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {

	slice := generateSlice(20)
	// fmt.Println("--- Unsorted ---------------------", slice)

	start := time.Now()
	// slice2 :=
	mergeSort(slice)
	// fmt.Println("--- Sorted -----------------------", slice2)
	fmt.Println("  -> Sorted ----------- secs: ", time.Since(start).Seconds())

	start1 := time.Now()
	// slice3 :=
	mergeSortGo(slice)
	// fmt.Println("--- Sorted with mergeSortGo ------", slice3)
	fmt.Println("  -> mergeSortGo ------ secs: ", time.Since(start1).Seconds())

	start2 := time.Now()
	// slice4 :=
	mergeSortGoPar(slice)
	// fmt.Println("--- Sorted with mergeSortGoPar ---", slice4)
	fmt.Println("  -> mergeSortGoPar --- secs: ", time.Since(start2).Seconds())

	start3 := time.Now()
	//slice5 :=
	mergeSortGoPar2(slice)
	// fmt.Println("--- Sorted with mergeSortGoPar ---", slice5)
	fmt.Println("  -> mergeSortGoParWG - secs: ", time.Since(start3).Seconds())
}

// Generates a slice of size, size filled with random numbers
func generateSlice(size int) []int {

	slice := make([]int, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(999) - rand.Intn(999)
	}
	return slice
}

// ---------------------------------------------------------------------
// mergeSortGo usa facilidades de slices

func mergeSortGo(s []int) []int {
	if len(s) > 1 {
		middle := len(s) / 2
		return merge(mergeSortGo(s[:middle]), mergeSortGo(s[middle:]))
	}
	return s
}

// ---------------------------------------------------------------------
// mergeSortGo usa facilidades de slices, recursivo e paralelo

func mergeSortGoPar(s []int) []int {

	if len(s) > 1 {
		middle := len(s) / 2

		var s1 []int
		var s2 []int
		c := make(chan struct{}, 2)

		go func() {
			s1 = mergeSortGoPar(s[middle:])
			c <- struct{}{}
		}()

		go func() {
			s2 = mergeSortGoPar(s[:middle])
			c <- struct{}{}
		}()

		<-c
		<-c

		return merge(s1, s2)

	}
	return s
}

// ---------------------------------------------------------------------
// mergeSortGo usa facilidades de slices, recursivo e paralelo, usando sincronizacao da biblioteca

func mergeSortGoPar2(s []int) []int {

	if len(s) > 1 {
		middle := len(s) / 2

		var s1 []int
		var s2 []int

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			s1 = mergeSortGoPar2(s[middle:])
		}()

		go func() {
			defer wg.Done()
			s2 = mergeSortGoPar2(s[:middle])
		}()

		wg.Wait()
		return merge(s1, s2)

	}
	return s
}

// ---------------------------------------------------------------------
// mergeSortGo ee uma implementacao tradicional

func mergeSort(items []int) []int {
	var num = len(items)

	if num == 1 {
		return items
	}

	middle := int(num / 2)
	var (
		left  = make([]int, middle)
		right = make([]int, num-middle)
	)
	for i := 0; i < num; i++ {
		if i < middle {
			left[i] = items[i]
		} else {
			right[i-middle] = items[i]
		}
	}

	return merge(mergeSort(left), mergeSort(right))
}

func merge(left, right []int) (result []int) {
	result = make([]int, len(left)+len(right))

	i := 0
	for len(left) > 0 && len(right) > 0 {
		if left[0] < right[0] {
			result[i] = left[0]
			left = left[1:]
		} else {
			result[i] = right[0]
			right = right[1:]
		}
		i++
	}

	for j := 0; j < len(left); j++ {
		result[i] = left[j]
		i++
	}
	for j := 0; j < len(right); j++ {
		result[i] = right[j]
		i++
	}

	return
}
