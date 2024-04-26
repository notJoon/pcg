package pcg

import (
	"runtime"
	"testing"
)

func TestPCG_MarshalBinary_Stress(t *testing.T) {
	p := NewPCG(12345, 67890)

	const iters = 100000000
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	initialAlloc := memStats.Alloc
	initialTotalAlloc := memStats.TotalAlloc

	for i := 0; i < iters; i++ {
		b, err := p.MarshalBinary()
		if err != nil {
			t.Fatalf("MarshalBinary failed: %v", err)
		}
		if len(b) != 20 {
			t.Errorf("MarshalBinary returned a slice of length %d; expected 20", len(b))
		}

		if b[0] != 'p' || b[1] != 'c' || b[2] != 'g' || b[3] != ':' {
			t.Errorf("MarshalBinary returned invalid byte slice: %v", b)
		}
	}

	runtime.ReadMemStats(&memStats)
	finalAlloc := memStats.Alloc
	finalTotalAlloc := memStats.TotalAlloc

	t.Logf("Memory usage:")
	t.Logf("  Initial allocated bytes: %d", initialAlloc)
	t.Logf("  Final allocated bytes: %d", finalAlloc)
	t.Logf("  Difference: %d", finalAlloc-initialAlloc)
	t.Logf("  Initial total allocated bytes: %d", initialTotalAlloc)
	t.Logf("  Final total allocated bytes: %d", finalTotalAlloc)
	t.Logf("  Difference: %d", finalTotalAlloc-initialTotalAlloc)
}
