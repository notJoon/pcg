package pcg

import (
	"errors"
	"math/bits"
	"unsafe"
)

// A PCG64 is a PCG64 generator with 128 bits of internal state.
// A zero PCG64 is equivalent to one seeded with 0.
type PCG64 struct {
	hi, lo *PCG32
}

// NewPCG64 returns a new PCG64 generator seeded with thr given values.
// seed1 and seed2 are the initial state values for the generator.
func NewPCG64(seed1, seed2 uint64) *PCG64 {
	return &PCG64{
		hi: NewPCG32().Seed(seed1, 0),
		lo: NewPCG32().Seed(seed2, 0),
	}
}

// SeedPCG64 initializes the PCG64 generator with the given state and sequence values.
// seed1 and seed2 are the initial state values, and seq1 and seq2 are the sequence values.
func (p *PCG64) SeedPCG64(seed1, seed2, seq1, seq2 uint64) *PCG64 {
	mask := ^uint64(0) >> 1
	if seq1&mask == seq2&mask {
		seq2 = ^seq2
	}
	p.lo.Seed(seed1, seq1)
	p.hi.Seed(seed2, seq2)

	return p
}

// NextUint64 generates a pseudorandom 64-bit unsigned integer using the PCG64 algorithm.
func (p *PCG64) NextUint64() uint64 {
	return uint64(p.hi.NextUint32())<<32 | uint64(p.lo.NextUint32())
}

// NextUint64InRange generates a pseudorandom number in the range [0, bound) using the PCG64 algorithm.
func (p *PCG64) NextUint64InRange(bound uint64) uint64 {
	threshold := -bound % bound
	for {
		r := p.NextUint64()
		if r >= threshold {
			return r % bound
		}
	}
}

// AdvancePCG64 moves the PCG64 generator forward by `delta` steps.
// It updates the initial state of the generator.
func (p *PCG64) AdvancePCG64(delta uint64) *PCG64 {
	p.hi.AdvancePCG32(delta)
	p.lo.AdvancePCG32(delta)
	return p
}

// RetreatPCG64 moves the PCG64 generator backward by `delta` steps.
// it updates the initial state of the generator.
func (p *PCG64) RetreatPCG64(delta uint64) *PCG64 {
	safeDelta := ^uint64(0) - 1
	p.AdvancePCG64(safeDelta)
	return p
}

func beUint64(b []byte) uint64 {
	_ = b[7]
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

func bePutUint64(b []byte, v uint64) {
	_ = b[7]
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
}

// MarshalBinaryPCG64 serializes the state of the PCG64 generator to a binary format.
// It returns the serialized state as a byte slice.
func (p *PCG64) MarshalBinaryPCG64() ([]byte, error) {
	b := make([]byte, 20)
	copy(b, "pcg:")
	bePutUint64(b[4:], p.hi.state)
	bePutUint64(b[4+8:], p.lo.state)
	return b, nil
}

func bePutUint64Unsafe(b []byte, v uint64) {
	*(*uint64)(unsafe.Pointer(&b[0])) = v
}

// MarshalBinaryPCG64Unsafe serializes the state of the PCG64 generator to a binary format using unsafe operations.
// It returns the serialized state as a byte slice.
// This method does not allocate any memory and is about 30 times faster than the safe version.
// However, it should be used with caution as it relies on unsafe operations.
func (p *PCG64) MarshalBinaryUnsafe() ([]byte, error) {
	b := make([]byte, 20)
	*(*uint32)(unsafe.Pointer(&b[0])) = *(*uint32)(unsafe.Pointer(&[4]byte{'p', 'c', 'g', ':'}))
	bePutUint64Unsafe(b[4:], p.hi.state)
	bePutUint64Unsafe(b[4+8:], p.lo.state)
	return b, nil
}

var errUnmarshalPCG = errors.New("invalid PCG encoding")

// UnmarshalBinaryPCG64 deserializes the state of the PCG64 generator from a binary format.
// It takes the serialized state as a byte slice and updates the generator's state.
func (p *PCG64) UnmarshalBinary(b []byte) error {
	if len(b) != 20 || string(b[:4]) != "pcg:" {
		return errUnmarshalPCG
	}
	p.hi.state = beUint64(b[4:])
	p.lo.state = beUint64(b[4+8:])
	return nil
}

func (p *PCG64) next() (uint64, uint64) {
	const (
		mulHi = 2549297995355413924
		mulLo = 4865540595714422341
		incHi = 6364136223846793005
		incLo = 1442695040888963407
	)

	// state = state * mul + inc
	hi, lo := bits.Mul64(p.lo.state, mulLo)
	hi += p.hi.state*mulLo + p.lo.state*mulHi

	lo, c := bits.Add64(lo, incLo, 0)
	hi, _ = bits.Add64(hi, incHi, c)

	p.lo.state = lo
	p.hi.state = hi

	return hi, lo
}

// NextUInt64WithMCG generates a pseudorandom 64-bit unsigned integer using the PCG64 algorithm with Multiplier Congruential Generator (MCG).
// It updates the internal state of the generator and returns the generated value.
func (p *PCG64) NextUint64WithMCG() uint64 {
	hi, lo := p.next()

	// ref: https://www.pcg-random.org/posts/128-bit-mcg-passes-practrand.html (#64-bit Multiplier)
	const cheapMul = 0xda942042e4dd58b5 // 15750249268501108917
	hi ^= hi >> 22
	hi *= cheapMul
	hi ^= hi >> 48
	hi *= (lo | 1)

	return hi
}
