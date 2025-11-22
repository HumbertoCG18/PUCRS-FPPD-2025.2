package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type SpeedupPoint struct {
	P       int
	Speedup float64
}

func monteCarloSequential(N int) (int, float64) {
	count := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < N; i++ {
		x := r.Float64()
		y := r.Float64()
		if x*x+y*y <= 1.0 {
			count++
		}
	}
	pi := 4.0 * float64(count) / float64(N)
	return count, pi
}

func monteCarloParallel(N int, P int) (int, float64) {
	var wg sync.WaitGroup
	counts := make([]int, P)
	chunk := N / P
	rem := N % P

	for i := 0; i < P; i++ {
		n := chunk
		if i < rem {
			n++
		}
		wg.Add(1)
		idx := i
		go func(seed int64, points int, idx int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(seed))
			localCount := 0
			for j := 0; j < points; j++ {
				x := r.Float64()
				y := r.Float64()
				if x*x+y*y <= 1.0 {
					localCount++
				}
			}
			counts[idx] = localCount
		}(time.Now().UnixNano()+int64(i)*37, n, idx)
	}

	wg.Wait()

	total := 0
	for _, c := range counts {
		total += c
	}
	pi := 4.0 * float64(total) / float64(N)
	return total, pi
}

func plotSpeedup(points []SpeedupPoint, filename string, title string) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = "Processadores (P)"
	p.Y.Label.Text = "Speedup"

	pts := make(plotter.XYs, len(points))
	for i, v := range points {
		pts[i].X = float64(v.P)
		pts[i].Y = v.Speedup
	}

	line, _ := plotter.NewLine(pts)
	line.LineStyle.Width = vg.Points(2)
	p.Add(line)

	p.Add(plotter.NewGrid())

	if err := p.Save(8*vg.Inch, 4*vg.Inch, filename); err != nil {
		panic(err)
	}
}

func main() {
	Ns := []int{1_000_000, 5_000_000, 20_000_000}
	Ps := []int{1, 2, 4, 8}

	for _, N := range Ns {

		// -------------------
		// EXECUÇÃO SEQUENCIAL
		// -------------------
		runtime.GOMAXPROCS(1)
		start := time.Now()
		_, piSeq := monteCarloSequential(N)
		tSeq := time.Since(start).Seconds()

		fmt.Printf("SEQ,N=%d,tempo=%f,pi=%f\n", N, tSeq, piSeq)

		speedData := []SpeedupPoint{{P: 1, Speedup: 1.0}}

		// -------------------
		// EXECUÇÃO PARALELA
		// -------------------
		for _, P := range Ps {
			runtime.GOMAXPROCS(P)
			start2 := time.Now()
			_, piPar := monteCarloParallel(N, P)
			tPar := time.Since(start2).Seconds()

			speedup := tSeq / tPar
			speedData = append(speedData, SpeedupPoint{P: P, Speedup: speedup})

			fmt.Printf("CSV,N=%d,P=%d,tempo=%f,pi=%f,speedup=%.2f\n",
				N, P, tPar, piPar, speedup)
		}

		// -------------------
		// GERAR GRÁFICO
		// -------------------
		filename := fmt.Sprintf("speedup_N%d.png", N)
		title := fmt.Sprintf("Speedup Monte Carlo - N=%d", N)

		plotSpeedup(speedData, filename, title)

		fmt.Println("Gráfico gerado:", filename)
		fmt.Println("------------------------------")
	}
}