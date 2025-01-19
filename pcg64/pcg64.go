// Package pcg64 implements the PCG-64 random number generator.
package pcg64

import (
	"encoding/binary"
	"errors"
	"math/bits"
	"sync"
	"time"
)

// PCG64 represents the state of a PCG-64 random number generator.
type PCG64 struct {
	state uint128
	inc   uint128
}

// SafePCG64 represents the state of a PCG-64 random number generator with a mutex to make it safe for concurrent use.
type SafePCG64 struct {
	PCG64
	mu sync.Mutex
}

// uint128 is a simple representation of a 128-bit unsigned integer
type uint128 struct {
	low  uint64
	high uint64
}

// New creates a new PCG64 instance seeded with the current time.
func New() *PCG64 {
	p := &PCG64{}
	p.Seed(uint64(time.Now().UnixNano()))
	return p
}

// NewSafe creates a new safe PCG64 instance seeded with the current time.
func NewSafe() *SafePCG64 {
	p := &SafePCG64{}
	p.Seed(uint64(time.Now().UnixNano()))
	return p
}

// State returns the current state of the random number generator.
func (p *PCG64) State() (uint64, uint64, uint64, uint64) {
	return p.state.low, p.state.high, p.inc.low, p.inc.high
}

// State returns the current state of the random number generator, which is safe for concurrent use.
func (p *SafePCG64) State() (uint64, uint64, uint64, uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.PCG64.State()
}

// Reset resets the state of the random number generator to the seed value.
func (p *PCG64) Reset() {
	p.Seed(uint64(time.Now().UnixNano()))
}

// Reset resets the state of the random number generator to the seed value, which is safe for concurrent use.
func (p *SafePCG64) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.PCG64.Reset()
}

// Marshal returns the binary encoding of the current state of the random number generator.
func (p *PCG64) Marshal() ([]byte, error) {
	buf := make([]byte, 32)
	binary.LittleEndian.PutUint64(buf[0:], p.state.low)
	binary.LittleEndian.PutUint64(buf[8:], p.state.high)
	binary.LittleEndian.PutUint64(buf[16:], p.inc.low)
	binary.LittleEndian.PutUint64(buf[24:], p.inc.high)
	return buf, nil
}

// Unmarshal sets the state of the random number generator to the state represented by the input data.
func (p *PCG64) Unmarshal(data []byte) error {
	if len(data) != 32 {
		return errors.New("pcg64: invalid state length")
	}
	p.state.low = binary.LittleEndian.Uint64(data[0:])
	p.state.high = binary.LittleEndian.Uint64(data[8:])
	p.inc.low = binary.LittleEndian.Uint64(data[16:])
	p.inc.high = binary.LittleEndian.Uint64(data[24:])
	return nil
}

// Seed initializes the state of the random number generator with the given seed value.
func (p *PCG64) Seed(seed uint64) {
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	p.state = uint128{low: 0, high: 0}
	p.inc = uint128{low: seed | 1, high: 0}
	p.Next()
	p.state = add128(p.state, p.inc)
	p.Next()
}

// Next generates a random 64-bit unsigned integer.
func (p *PCG64) Next() uint64 {
	oldState := p.state
	p.state = add128(mul128(p.state, uint128{low: 0x5851f42d4c957f2d, high: 0x14057b7ef767814f}), p.inc)
	xorshifted := uint64(((oldState.high ^ oldState.low) >> 29) | ((oldState.high ^ oldState.low) << 35))
	rot := uint64(oldState.high >> 58)
	return bits.RotateLeft64(xorshifted, -int(rot))
}

// Float64 generates a random float64 in the range [0.0, 1.0).
func (p *PCG64) Float64() float64 {
	return float64(p.Next()>>(64-53)) / (1 << 53)
}

// Seed initializes the state of the random number generator with the given seed value, which is safe for concurrent use.
func (p *SafePCG64) Seed(seed uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.PCG64.Seed(seed)
}

// Next generates a random 64-bit unsigned integer, which is safe for concurrent use.
func (p *SafePCG64) Next() uint64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.PCG64.Next()
}

// Float64 generates a random float64 in the range [0.0, 1.0), which is safe for concurrent use.
func (p *SafePCG64) Float64() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.PCG64.Float64()
}

// Int64 generates a random 64-bit signed integer.
func (p *PCG64) Int64() int64 {
	return int64(p.Next() >> 1)
}

// Uint32 generates a random 32-bit unsigned integer.
func (p *PCG64) Uint32() uint32 {
	return uint32(p.Next() >> 32)
}

// Int32 generates a random 32-bit signed integer.
func (p *PCG64) Int32() int32 {
	return int32(p.Uint32() >> 1)
}

// Float32 generates a random float32 in the range [0.0, 1.0).
func (p *PCG64) Float32() float32 {
	return float32(p.Uint32()>>(32-24)) / (1 << 24)
}

// Int generates a random integer in the range [0, n).
func (p *PCG64) Int(n int) int {
	if n <= 0 {
		panic("pcg64: argument to Int is <= 0")
	}
	if n&(n-1) == 0 { // n is 2^m, use mask
		return int(p.Next() & uint64(n-1))
	}
	max := uint64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := p.Next()
	for v > max {
		v = p.Next()
	}
	return int(v % uint64(n))
}

// Helper functions for 128-bit arithmetic

func add128(a, b uint128) uint128 {
	low := a.low + b.low
	high := a.high + b.high
	if low < a.low {
		high++
	}
	return uint128{low: low, high: high}
}

func mul128(a, b uint128) uint128 {
	// if the number is smaller then uint64,then use normal uint64 * uint64
	if a.high == 0 && b.high == 0 {
		return uint128{low: a.low * b.low, high: 0}
	}

	// split into 2 parts
	a1, a0 := a.high, a.low
	b1, b0 := b.high, b.low

	// Karatsuba
	z0 := a0 * b0
	z2 := a1 * b1
	z1 := (a0+a1)*(b0+b1) - z0 - z2

	low := z0 + (z1 << 64)
	high := z2 + (z1 >> 64)

	if low < z0 {
		high++
	}

	return uint128{low: low, high: high}
}
