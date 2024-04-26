package pcg

import (
	"testing"
)

func TestPCG_Advance(t *testing.T) {
	pcg := NewPCG(12345, 67890)

	testCases := []struct {
		delta           uint64
		expectedStateHi uint64
		expectedStateLo uint64
	}{
		{1, 16443432798917770532, 1294492316257287365},
		{10, 9073714748428748454, 9095006751169262415},
		{100, 1498360792142116778, 11040029025224029795},
		{1000, 7761321322648589714, 770061004744980459},
		{10000, 8930526547519973282, 18106490617456118331},
	}

	for _, tc := range testCases {
		pcg.Advance(tc.delta)
		if pcg.hi.state != tc.expectedStateHi {
			t.Errorf("Advance(%d) hi state = %d; expected %d", tc.delta, pcg.hi.state, tc.expectedStateHi)
		}
		if pcg.lo.state != tc.expectedStateLo {
			t.Errorf("Advance(%d) lo state = %d; expected %d", tc.delta, pcg.lo.state, tc.expectedStateLo)
		}
	}
}

func TestPCG(t *testing.T) {
	p := NewPCG(1, 2)
	want := []uint64{
		0x52addb9b0d4aa107,
		0xc5d5c81b8c97ff8f,
		0xcfa82191c9a86caa,
		0x76b48e618586fdfe,
		0x765ac4ba3e566855,
		0x1d6058a5dd7ab27,
		0x2b913f2f76e81329,
		0x74873f4e5348d32e,
		0xc4c940eb70248174,
		0xb5a1651a6627a924,
		0xc34174eb7f136d0a,
		0xe612b37df73df71c,
		0x884a2539ea7aa198,
		0x2976010a57986e59,
		0x1d0d522531d62a7d,
		0xa7da1ad05db25a75,
		0xdbee2df7bd6428be,
		0x598c54d1eb4abdd7,
		0x559ca964532a3777,
		0x6e64af73ece533b0,
	}

	for i, x := range want {
		u := p.Uint64()
		if u != x {
			t.Errorf("PCG #%d = %#x, want %#x", i, u, x)
		}
	}
}

func TestPCG_Retreat(t *testing.T) {
	pcg := NewPCG(12345, 67890)

	testCases := []struct {
		delta           uint64
		expectedStateHi uint64
		expectedStateLo uint64
	}{
		{1, 7265056988599925051, 16912344864586758584},
		{10, 15578097273240930873, 13711579158205810606},
		{100, 3761525201756208775, 6157393363865312820},
		{1000, 15336446625969592741, 13630190462364618442},
		{10000, 10106684222517973779, 4620269966716251888},
	}

	for _, tc := range testCases {
		pcg.Retreat(tc.delta)
		if pcg.hi.state != tc.expectedStateHi {
			t.Errorf("Retreat(%d) hi state = %d; expected %d", tc.delta, pcg.hi.state, tc.expectedStateHi)
		}
		if pcg.lo.state != tc.expectedStateLo {
			t.Errorf("Retreat(%d) lo state = %d; expected %d", tc.delta, pcg.lo.state, tc.expectedStateLo)
		}
	}
}

func TestPCG_MarshalBinaryUnsafe(t *testing.T) {
	pcg := NewPCG(12345, 67890)
	b, err := pcg.MarshalBinaryUnsafe()
	if err != nil {
		t.Fatalf("MarshalBinaryUnsafe() error = %v; want nil", err)
	}
	if len(b) != 20 {
		t.Errorf("MarshalBinaryUnsafe() len(b) = %d; want 20", len(b))
	}
	if string(b[:4]) != "pcg:" {
		t.Errorf("MarshalBinaryUnsafe() b[:4] = %s; want 'pcg:'", string(b[:4]))
	}
}

func BenchmarkPCG_Seed(b *testing.B) {
	pcg := NewPCG(0, 0)
	for i := 0; i < b.N; i++ {
		pcg.Seed(12345, 67890, 12345, 67890)
	}
}

func BenchmarkPCG_MarshalBinary(b *testing.B) {
	pcg := NewPCG(12345, 67890)
	for i := 0; i < b.N; i++ {
		_, _ = pcg.MarshalBinary()
	}
}

func BenchmarkPCG_MarshalBinary_Unsafe(b *testing.B) {
	pcg := NewPCG(12345, 67890)
	for i := 0; i < b.N; i++ {
		_, _ = pcg.MarshalBinaryUnsafe()
	}
}
