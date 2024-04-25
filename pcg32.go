package pcg

// ref: https://gist.github.com/ivan-pi/060e38d5f9a86c57923a61fbf18d095c
const (
	defaultState  = 0x853c49e6748fea9b //  9600629759793949339
	multiplier    = 0x5851f42d4c957f2d //  6364136223846793005
	incrementStep = 0xda3e39cb94b95bdb // 15726070495360670683
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

// Random generates a pseudorandom 32-bit unsigned integer using the PCG32 algorithm.
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
func (p *PCG32) Random() uint32 {
	old := p.state
	p.state = old*multiplier + p.increment

	xorshifted := uint32(((old >> 18) ^ old) >> 27)
	rot := uint32(old >> 59)

	return (xorshifted >> rot) | (xorshifted << (neg_mask - rot))
}

// Bounded generates a pseudorandom number in the range [0, bound) using the PCG32 algorithm.
func (p *PCG32) Bounded(bound uint32) uint32 {
	if bound == 0 {
		return 0
	}

	threshold := -bound % bound
	for {
		r := p.Random()
		if r >= threshold {
			return r % bound
		}
	}
}

// lcg64 is an implementation of a 64-bit linear congruential generator (LCG).
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
func (p *PCG32) lcg64(state, delta, mul, add uint64) uint64 {
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
	p.state = p.lcg64(p.state, delta, multiplier, incrementStep)
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
