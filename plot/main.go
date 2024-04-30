package main

import (
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/notJoon/pcg"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	numBins    = 100
	numSamples = 25000
)

func generateRandHistogram(bins []int, seed int64) {
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < numSamples; i++ {
		binIndex := r.Intn(numBins)
		bins[binIndex]++
	}
}

func generatePCG32Histogram(bins []int, seed int64) {
	pcg := pcg.NewPCG32().Seed(uint64(seed), uint64(seed))
	for i := 0; i < numSamples; i++ {
		r := pcg.Uint32()
		binIndex := int(uint64(r) * uint64(numBins) >> 32)
		bins[binIndex]++
	}
}

func main() {
	seed := time.Now().Unix()

	binsRand := make([]int, numBins)
	binsPCG := make([]int, numBins)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		generateRandHistogram(binsRand, seed)
	}()
	go func() {
		defer wg.Done()
		generatePCG32Histogram(binsPCG, seed)
	}()

	wg.Wait()

	// Creating plots
	pRand := plot.New()
	pRand.Title.Text = "math/rand Distribution"
	pRand.X.Label.Text = "Bin"
	pRand.Y.Label.Text = "Count"

	pPCG32 := plot.New()
	pPCG32.Title.Text = "PCG Distribution"
	pPCG32.X.Label.Text = "Bin"
	pPCG32.Y.Label.Text = "Count"

	pPCG64 := plot.New()
	pPCG64.Title.Text = "PCG64 Distribution"
	pPCG64.X.Label.Text = "Bin"
	pPCG64.Y.Label.Text = "Count"

	// Rand histogram
	randHistogram, err := plotter.NewBarChart(plotter.Values(makePlotterValues(binsRand)), vg.Points(5))
	if err != nil {
		panic(err)
	}
	randHistogram.Color = color.RGBA{R: 255, A: 255} // Red

	// PCG32 histogram
	pcg32Histogram, err := plotter.NewBarChart(plotter.Values(makePlotterValues(binsPCG)), vg.Points(5))
	if err != nil {
		panic(err)
	}
	pcg32Histogram.Color = color.RGBA{B: 255, A: 255} // Blue

	pRand.Add(randHistogram)
	pPCG32.Add(pcg32Histogram)

	if err := pRand.Save(8*vg.Inch, 4*vg.Inch, "rand_histogram.png"); err != nil {
		panic(err)
	}

	if err := pPCG32.Save(8*vg.Inch, 4*vg.Inch, "pcg32_histogram.png"); err != nil {
		panic(err)
	}
}

func makePlotterValues(bins []int) []float64 {
	values := make([]float64, len(bins))
	for i, count := range bins {
		values[i] = float64(count)
	}
	return values
}