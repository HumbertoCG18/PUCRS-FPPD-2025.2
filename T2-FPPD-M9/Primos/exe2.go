package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Result struct {
	P int
	Speedup float64
}

func isPrime(n uint64) bool {
	if n < 2 {
		return false
	}
	lim := uint64(math.Sqrt(float64(n)))
	for i := uint64(2); i <= lim; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func countPrimesSequential(nums []uint64) int {
	total := 0
	for _, n := range nums {
		if isPrime(n) {
			total++
		}
	}
	return total
}

func countPrimesParallel(nums []uint64, P int) int {
	bag := make(chan uint64)
	results := make(chan int)
	for i := 0; i < P; i++ {
		go func() {
			count := 0
			for n := range bag {
				if isPrime(n) {
					count++
				}
			}
			results <- count
		}()
	}

	go func() {
		for _, n := range nums {
			bag <- n
		}
		close(bag)
	}()

	total := 0
	for i := 0; i < P; i++ {
		total += <-results
	}
	return total
}

func PlotSpeedup(results []Result) {
	p := plot.New()
	p.Title.Text = "Speedup Bag of Tasks"
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
	N := 20000
	processors := []int{1, 2, 4, 8}

	nums := make([]uint64, N)
	for i := range nums {
		nums[i] = uint64(10_000_000 + rand.Intn(5000))
	}

	start := time.Now()
	seq := countPrimesSequential(nums)
	seqTime := time.Since(start).Seconds()
	fmt.Println("Sequencial:", seqTime, "Primos:", seq)

	var results []Result

	for _, P := range processors {
		start := time.Now()
		par := countPrimesParallel(nums, P)
		parTime := time.Since(start).Seconds()
		speed := seqTime / parTime

		fmt.Printf("P=%d -> tempo=%.3f speedup=%.2f encontrados=%d\n", P, parTime, speed, par)
		results = append(results, Result{P: P, Speedup: speed})
	}

	PlotSpeedup(results)
	fmt.Println("Gerado speedup.png")
}