package pcg

import (
	"math/rand"
	"testing"
)

func TestPCG32_Bounded(t *testing.T) {
	pcg := NewPCG32().Seed(12345, 67890)

	testCases := []struct {
		bound uint32
	}{
		{0},
		{1},
		{10},
		{100},
		{1000},
		{10000},
	}

	for _, tc := range testCases {
		result := pcg.Uint32Range(tc.bound)
		if tc.bound != 0 && result >= tc.bound {
			t.Errorf("Bounded(%d) = %d; expected a value between 0 and %d", tc.bound, result, tc.bound)
		}
		if tc.bound == 0 && result != 0 {
			t.Errorf("Bounded(%d) = %d; expected 0", tc.bound, result)
		}
	}
}

func TestPCG32_UniformDistribution(t *testing.T) {
	pcg := NewPCG32().Seed(12345, 67890)
	numBins := 10
	numSamples := 1000000
	toleranceRatio := 10 // 10% tolerance
	bins := make([]int, numBins)

	for i := 0; i < numSamples; i++ {
		r := pcg.NextUint32()
		binIndex := int(uint64(r) * uint64(numBins) >> 32)
		bins[binIndex]++
	}

	expected := numSamples / numBins
	tolerance := expected / toleranceRatio

	for _, count := range bins {
		if abs(count-expected) > tolerance {
			t.Errorf("bin count %d is outside the expected range [%d, %d]", count, expected-tolerance, expected+tolerance)
		}
	}
}

func TestPCG32_lcg64(t *testing.T) {
	pcg := NewPCG32()

	testCases := []struct {
		state    uint64
		delta    uint64
		expected uint64
	}{
		{1, 1, 3643462645497912072},
		{1, 10, 15256603694110904427},
		{1, 100, 5234694153321213237},
		{1, 1000, 2323235076269450313},
		{1, 10000, 6143568259046921169},
	}

	for _, tc := range testCases {
		result := pcg.advancedLCG64(tc.state, tc.delta, multiplier, incrementStep)
		if result != tc.expected {
			t.Errorf("lcg64(%d, %d) = %d; expected %d", tc.state, tc.delta, result, tc.expected)
		}
	}
}

func TestPCG32_Advance(t *testing.T) {
	pcg := NewPCG32()

	testCases := []struct {
		initialState  uint64
		delta         uint64
		expectedState uint64
	}{
		{1, 1, 3643462645497912072},
		{1, 10, 15256603694110904427},
		{1, 100, 5234694153321213237},
		{1, 1000, 2323235076269450313},
		{1, 10000, 6143568259046921169},
	}

	for _, tc := range testCases {
		pcg.state = tc.initialState
		pcg.AdvancePCG32(tc.delta)
		if pcg.state != tc.expectedState {
			t.Errorf("Advance(%d) = %d; expected %d", tc.delta, pcg.state, tc.expectedState)
		}
	}
}

func TestPCG32_Retreat(t *testing.T) {
	pcg := NewPCG32()

	testCases := []struct {
		initialState  uint64
		delta         uint64
		expectedState uint64
	}{
		{3643462645497912072, 1, 1},
		{15256603694110904427, 10, 1},
		{5234694153321213237, 100, 1},
		{2323235076269450313, 1000, 1},
		{6143568259046921169, 10000, 1},
	}

	for _, tc := range testCases {
		pcg.state = tc.initialState
		pcg.RetreatPCG32(tc.delta)
		if pcg.state != tc.expectedState {
			t.Errorf("Retreat(%d) = %d; expected %d", tc.delta, pcg.state, tc.expectedState)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func BenchmarkPCG32(b *testing.B) {
	rng := NewPCG32()
	for i := 0; i < b.N; i++ {
		rng.NextUint32()
	}
}

func BenchmarkMathRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Uint32()
	}
}

func BenchmarkPCG32_InRange(b *testing.B) {
	rng := NewPCG32()
	for i := 0; i < b.N; i++ {
		rng.Uint32Range(100)
	}
}

func Benchmark_MathRanIntn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Intn(100)
	}
}
