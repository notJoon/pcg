package main

import (
	"github.com/notJoon/pcg"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	pcg := pcg.NewPCG32().Seed(12345, 67890)
	numBins := 10
	numSamples := 1000000
	bins := make([]int, numBins)

	for i := 0; i < numSamples; i++ {
		r := pcg.NextUint32()
		binIndex := int(uint64(r) * uint64(numBins) >> 32)
		bins[binIndex]++
	}

	pl := plot.New()

	pl.Title.Text = "PCG32 Uniform Distribution"
	pl.X.Label.Text = "Bin"
	pl.Y.Label.Text = "Count"

	w := vg.Points(20)
	bars := make(plotter.Values, numBins)
	for i := range numBins {
		bars[i] = float64(bins[i])
	}

	bar, err := plotter.NewBarChart(bars, w)
	if err != nil {
		panic(err)
	}

	pl.Add(bar)

	if err := pl.Save(4*vg.Inch, 4*vg.Inch, "uniform_distribution.png"); err != nil {
		panic(err)
	}
}