package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Result struct {
	P       int
	Speedup float64
}

func merge(a, b []int) []int {
	out := make([]int, 0, len(a)+len(b))
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			out = append(out, a[i])
			i++
		} else {
			out = append(out, b[j])
			j++
		}
	}

	out = append(out, a[i:]...)
	out = append(out, b[j:]...)
	return out
}

func mergeSort(a []int) []int {
	if len(a) <= 1 {
		return a
	}

	mid := len(a) / 2
	left := mergeSort(a[:mid])
	right := mergeSort(a[mid:])

	return merge(left, right)
}

func mergeSortParallel(a []int, T int) []int {
	if len(a) <= T {
		return mergeSort(a)
	}

	mid := len(a) / 2

	var left, right []int
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		left = mergeSortParallel(a[:mid], T)
	}()
	go func() {
		defer wg.Done()
		right = mergeSortParallel(a[mid:], T)
	}()
	wg.Wait()

	return merge(left, right)
}

func PlotSpeedup(results []Result) {
	p := plot.New()
	p.Title.Text = "Speedup Mergesort"
	p.X.Label.Text = "Processadores"
	p.Y.Label.Text = "Speedup"
	p.Add(plotter.NewGrid())

	pts := make(plotter.XYs, len(results))
	for i, r := range results {
		pts[i].X = float64(r.P)
		pts[i].Y = r.Speedup
	}

	line, _ := plotter.NewLine(pts)
	line.LineStyle.Width = vg.Points(2)
	p.Add(line)

	p.Save(8*vg.Inch, 4*vg.Inch, "speedup.png")
}

func main() {
	N := 2_000_000
	processors := []int{1, 2, 4, 8}
	T := N // substituído abaixo para cada P

	a := make([]int, N)
	for i := range a {
		a[i] = rand.Int()
	}

	// sequencial
	cpy := make([]int, len(a))
	copy(cpy, a)
	start := time.Now()
	mergeSort(cpy)
	seqTime := time.Since(start).Seconds()
	fmt.Println("Tempo sequencial:", seqTime)

	var results []Result

	for _, P := range processors {
		T = N / P
		copy(cpy, a)

		start := time.Now()
		mergeSortParallel(cpy, T)
		parTime := time.Since(start).Seconds()
		speed := seqTime / parTime

		fmt.Printf("P=%d  tempo=%.3fs speedup=%.2f\n", P, parTime, speed)
		results = append(results, Result{P: P, Speedup: speed})
	}

	PlotSpeedup(results)
	fmt.Println("Gerado speedup.png")
}