package pcg

import (
	"errors"
	"math/bits"
	"unsafe"
)

// A PCG is a PCG generator with 128 bits of internal state.
// A zero PCG is equivalent to one seeded with 0.
type PCG struct {
	hi, lo *PCG32
}

// NewPCG returns a new PCG seeded with the given values.
func NewPCG(seed1, seed2 uint64) *PCG {
	return &PCG{
		hi: NewPCG32().Seed(seed1, 0),
		lo: NewPCG32().Seed(seed2, 0),
	}
}

func (p *PCG) Seed(seed1, seed2, seq1, seq2 uint64) *PCG {
	mask := ^uint64(0) >> 1
	if seq1&mask == seq2&mask {
		seq2 = ^seq2
	}
	p.lo.Seed(seed1, seq1)
	p.hi.Seed(seed2, seq2)

	return p
}

func (p *PCG) Random() uint64 {
	return uint64(p.hi.Random())<<32 | uint64(p.lo.Random())
}

func (p *PCG) Bounded(bound uint64) uint64 {
	threshold := -bound % bound
	for {
		r := p.Random()
		if r >= threshold {
			return r % bound
		}
	}
}

func (p *PCG) Advance(delta uint64) *PCG {
	p.hi.Advance(delta)
	p.lo.Advance(delta)
	return p
}

func (p *PCG) Retreat(delta uint64) *PCG {
	safeDelta := ^uint64(0) - 1
	p.Advance(safeDelta)
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

func (p *PCG) MarshalBinary() ([]byte, error) {
	b := make([]byte, 20)
	copy(b, "pcg:")
	bePutUint64(b[4:], p.hi.state)
	bePutUint64(b[4+8:], p.lo.state)
	return b, nil
}

func bePutUint64Unsafe(b []byte, v uint64) {
	*(*uint64)(unsafe.Pointer(&b[0])) = v
}

func (p *PCG) MarshalBinaryUnsafe() ([]byte, error) {
	b := make([]byte, 20)
	*(*uint32)(unsafe.Pointer(&b[0])) = *(*uint32)(unsafe.Pointer(&[4]byte{'p', 'c', 'g', ':'}))
	bePutUint64Unsafe(b[4:], p.hi.state)
	bePutUint64Unsafe(b[4+8:], p.lo.state)
	return b, nil
}

var errUnmarshalPCG = errors.New("invalid PCG encoding")

func (p *PCG) UnmarshalBinary(b []byte) error {
	if len(b) != 20 || string(b[:4]) != "pcg:" {
		return errUnmarshalPCG
	}
	p.hi.state = beUint64(b[4:])
	p.lo.state = beUint64(b[4+8:])
	return nil
}

func (p *PCG) next() (uint64, uint64) {
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

func (p *PCG) Uint64() uint64 {
	hi, lo := p.next()

	// ref: https://www.pcg-random.org/posts/128-bit-mcg-passes-practrand.html (#64-bit Multiplier)
	const cheapMul = 0xda942042e4dd58b5 // 15750249268501108917
	hi ^= hi >> 22
	hi *= cheapMul
	hi ^= hi >> 48
	hi *= (lo | 1)

	return hi
}
