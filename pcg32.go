package pcg

import "encoding/binary"

// ref: https://gist.github.com/ivan-pi/060e38d5f9a86c57923a61fbf18d095c
const (
	defaultState  = 0x853c49e6748fea9b //  9600629759793949339
	multiplier    = 0x5851f42d4c957f2d //  6364136223846793005
	incrementStep = 0x9e3779b97f4a7c15 //  https://www.pcg-random.org/posts/bugs-in-splitmix.html
)

// PCG32 is a 32-bit pseudorandom number generator based on the PCG family of algorithms.
type PCG32 struct {
	state, increment uint64
}

// NewPCG32 creates a new PCG32 generator with the default state and sequence values.
func NewPCG32() *PCG32 {
	return &PCG32{
		state:     defaultState,
		increment: incrementStep,
	}
}

// Seed initializes the PCG32 generator with the given state and sequence values.
func (p *PCG32) Seed(state, sequence uint64) *PCG32 {
	p.increment = (sequence << 1) | 1
	p.state = (state+p.increment)*multiplier + incrementStep
	return p
}

// neg_mask is a mask to extract the lower 5 bits of a number.
const neg_mask = 31

// Uint32 generates a pseudorandom 32-bit unsigned integer using the PCG32 algorithm.
// It updates the internal state of the generator using the PCG32 formula:
//
//	state = state * multiplier + increment
//
// It then applies a series of bitwise operations to the old state to produce the random number:
//  1. XOR the old state with (old state >> 18).
//  2. Shift the result right by 27 bits to obtain `xorshifted`.
//  3. Calculate the rotation amount `rot` by shifting the old state right by 59 bits.
//  4. Rotate `xorshifted` right by `rot` bits and OR it with `xorshifted` rotated left by `((-rot) & 31)` bits.
//
// The resulting value is returned as the random number.
func (p *PCG32) Uint32() uint32 {
	old := p.state
	p.state = old*multiplier + p.increment

	xorshifted := uint32(((old >> 18) ^ old) >> 27)
	rot := uint32(old >> 59)

	return (xorshifted >> rot) | (xorshifted << (neg_mask - rot))
}

// Uintn32 generates a pseudorandom number in the range [0, bound) using the PCG32 algorithm.
func (p *PCG32) Uintn32(bound uint32) uint32 {
	if bound == 0 {
		return 0
	}

	threshold := -bound % bound
	for {
		r := p.Uint32()
		if r >= threshold {
			return r % bound
		}
	}
}

// Uint63 generates a pseudorandom 63-bit integer using two 32-bit numbers.
// The function ensures that the returned number is within the range of 0 to 2^63-1.
func (p *PCG32) Uint63() int64 {
	upper := int64(p.Uint32()) & 0x7FFFFFFF // Use only the lower 31 bits of the upper half
	lower := int64(p.Uint32())              // Use all 32 bits of the lower half
	return (upper << 32) | lower            // Combine the two halves to form a 63-bit integer
}

// advancedLCG64 is an implementation of a 64-bit linear congruential generator (LCG).
// It takes the following parameters:
//   - state: The current state of the LCG.
//   - delta: The number of steps to advance the LCG.
//   - mul: The multiplier of the LCG.
//   - add: The increment of the LCG.
//
// The function advances the LCG by `delta` steps and returns the new state.
//
// The LCG algorithm is defined by the following recurrence relation:
//
//	state(n+1) = (state(n) * mul + add) mod 2^64
//
// The function uses an efficient algorithm to advance the LCG by `delta` steps
// without iterating `delta` times. It exploits the properties of the LCG and
// uses binary exponentiation to calculate the result in logarithmic time.
//
// The algorithm works as follows:
//  1. Initialize `accMul` to 1 and `accAdd` to 0.
//  2. Iterate while `delta` is greater than 0:
//     - If the least significant bit of `delta` is 1:
//     - Multiply `accMul` by `mul`.
//     - Set `accAdd` to `accAdd * mul + add`.
//     - Update `add` to `(mul + 1) * add`.
//     - Update `mul` to `mul * mul`.
//     - Right-shift `delta` by 1 (divide by 2).
//  3. Return `accMul * state + accAdd`.
//
// The time complexity of this function is O(log(delta)), as it iterates logarithmically
// with respect to `delta`. The space complexity is O(1), as it uses only a constant
// amount of additional memory.
func (p *PCG32) advancedLCG64(state, delta, mul, add uint64) uint64 {
	accMul := uint64(1)
	accAdd := uint64(0)

	for delta > 0 {
		if delta&1 != 0 {
			accMul *= mul
			accAdd = accAdd*mul + add
		}
		add = (mul + 1) * add
		mul *= mul
		delta /= 2
	}
	return accMul*state + accAdd
}

// Advance moves the PCG32 generator forward by `delta` steps.
// It updates the internal state of the generator using the `lcg64` function
// and returns the updated PCG32 instance.
func (p *PCG32) Advance(delta uint64) *PCG32 {
	p.state = p.advancedLCG64(p.state, delta, multiplier, incrementStep)
	return p
}

// Retreat moves the PCG32 generator backward by `delta` steps.
// It calculates the equivalent forward delta using the two's complement of `delta`
// and calls the `Advance` function with the calculated delta.
// It returns the updated PCG32 instance.
func (p *PCG32) Retreat(delta uint64) *PCG32 {
	safeDelta := ^delta + 1
	return p.Advance(safeDelta)
}

func (p *PCG32) Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to shuffle")
	}
	if n < 2 {
		return
	}
	// Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	for i := n - 1; i > 0; i-- {
		j := int(p.Uintn32(uint32(i + 1)))
		swap(i, j)
	}
}

// Perm returns a slice of n integers. The slice is a random permutation of the integers [0, n).
func (p *PCG32) Perm(n int) []int {
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = i
	}
	for i := 1; i < n; i++ {
		j := int(p.Uintn32(uint32(i + 1)))
		res[i], res[j] = res[j], res[i]
	}
	return res
}

// Read generates and fills the given byte slice with random bytes using the PCG32 random number generator.
//
// This function repeatedly generates 32-bit unsigned integers using the Uint32 function,
// and uses binary.LittleEndian.PutUint32 to split these integers into bytes and store them in the slice.
// This approach efficiently copies memory in accordance with the CPU's endian configuration, thereby enhancing performance.
//
// Parameters:
//   - buf: The byte slice to be filled with random bytes.
//
// Return values:
//   - n: The number of bytes generated and stored in the byte slice. It is always equal to len(buf).
//   - err: Always returns nil, indicating no error occurred.
func (p *PCG32) Read(buf []byte) (int, error) {
	n := len(buf)
	i := 0

	// loop unrolling: process 8 bytes in each iteration
	for ; i <= n-8; i += 8 {
		val1 := p.Uint32()
		val2 := p.Uint32()
		binary.LittleEndian.PutUint32(buf[i:], val1)
		binary.LittleEndian.PutUint32(buf[i+4:], val2)
	}

	// handle remaining bytes (less than 8 bytes)
	if i < n {
		remaining := buf[i:]
		for j := 0; j < len(remaining); j += 4 {
			if i+j < n {
				val := p.Uint32()
				// handle remaining bytes (less than real buffer size)
				for k := 0; k < 4 && (j+k) < len(remaining); k++ {
					remaining[j+k] = byte(val >> (8 * k))
				}
			}
		}
	}

	return n, nil
}
