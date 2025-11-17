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

func piSequential(N int) int {
	r := rand.New(rand.NewSource(1))
	count := 0
	for i := 0; i < N; i++ {
		x := r.Float64()
		y := r.Float64()
		if x*x+y*y <= 1 {
			count++
		}
	}
	return count
}

func piParallel(N, P int) int {
	each := N / P
	results := make(chan int, P)
	var wg sync.WaitGroup

	for i := 0; i < P; i++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			r := rand.New(rand.NewSource(seed))
			count := 0
			for j := 0; j < each; j++ {
				x := r.Float64()
				y := r.Float64()
				if x*x+y*y <= 1 {
					count++
				}
			}
			results <- count
		}(time.Now().UnixNano() + int64(i))
	}

	wg.Wait()
	close(results)

	total := 0
	for c := range results {
		total += c
	}
	return total
}

func PlotSpeedup(results []Result) {
	p := plot.New()
	p.Title.Text = "Speedup Monte Carlo Pi"
	p.X.Label.Text = "Processadores"
	p.Y.Label.Text = "Speedup"
	p.Add(plotter.NewGrid())

	pts := make(plotter.XYs, len(results))
	for i, r := range results {
		pts[i].X = float64(r.P)
		pts[i].Y = r.Speedup
	}

	line, _ := plotter.NewLine(pts)
	p.Add(line)

	p.Save(8*vg.Inch, 4*vg.Inch, "speedup.png")
}

func main() {
	N := 10_000_000
	processors := []int{1, 2, 4, 8}

	start := time.Now()
	c := piSequential(N)
	seqTime := time.Since(start).Seconds()
	fmt.Println("Sequencial pi ≈", 4*float64(c)/float64(N), "tempo", seqTime)

	var results []Result

	for _, P := range processors {
		start := time.Now()
		cpar := piParallel(N, P)
		parTime := time.Since(start).Seconds()
		speed := seqTime / parTime

		fmt.Printf("P=%d pi≈%.6f tempo=%.3f speedup=%.2f\n",
			P, 4*float64(cpar)/float64(N), parTime, speed)

		results = append(results, Result{P: P, Speedup: speed})
	}

	PlotSpeedup(results)
	fmt.Println("Gerado speedup.png")
}