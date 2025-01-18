package pcg32

import (
	"encoding/binary"
	"errors"
	"math/bits"
	"time"
)

// PCG32 represents the state of a PCG-32 random number generator.
type PCG32 struct {
	state uint64
	inc   uint64
}

// New creates a new PCG32 instance seeded with the current time.
func New() *PCG32 {
	p := &PCG32{}
	p.Seed(uint64(time.Now().UnixNano()))
	return p
}

// State returns the current state of the random number generator.
func (p *PCG32) State() (uint64, uint64) {
	return p.state, p.inc
}

// Reset resets the state of the random number generator to the seed value.
func (p *PCG32) Reset() {
	p.Seed(uint64(time.Now().UnixNano()))
}

// Marshal returns the binary encoding of the current state of the random number generator.
func (p *PCG32) Marshal() ([]byte, error) {
	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf[0:], p.state)
	binary.LittleEndian.PutUint64(buf[8:], p.inc)
	return buf, nil
}

// Unmarshal sets the state of the random number generator to the state represented by the input data.
func (p *PCG32) Unmarshal(data []byte) error {
	if len(data) != 16 {
		return errors.New("pcg32: invalid state length")
	}
	p.state = binary.LittleEndian.Uint64(data[0:])
	p.inc = binary.LittleEndian.Uint64(data[8:])
	return nil
}

// Seed initializes the state of the random number generator with the given seed value.
func (p *PCG32) Seed(seed uint64) {
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	p.state = 0
	p.inc = (seed << 1) | 1
	p.Next()
	p.state += seed
	p.Next()
}

// Next generates a random 32-bit unsigned integer.
func (p *PCG32) Next() uint32 {
	oldState := p.state
	p.state = oldState*6364136223846793005 + p.inc
	xorshifted := uint32(((oldState >> 18) ^ oldState) >> 27)
	rot := uint32(oldState >> 59)
	return bits.RotateLeft32(xorshifted, -int(rot))
}

// Float64 generates a random float64 in the range [0.0, 1.0).
func (p *PCG32) Float64() float64 {
	return float64(p.Next()) / (1 << 32)
}

// Float32 generates a random float32 in the range [0.0, 1.0).
func (p *PCG32) Float32() float32 {
	return float32(p.Next()) / (1 << 32)
}

// Int31 generates a random 31-bit signed integer.
func (p *PCG32) Int32() int32 {
	return int32(p.Next() >> 1)
}

// Int generates a random integer in the range [0, n).
func (p *PCG32) Int(n int) int {
	if n <= 0 {
		panic("pcg32: argument to Int is <= 0")
	}
	if n&(n-1) == 0 { // n is 2^m, use mask
		return int(p.Next() & uint32(n-1))
	}
	max := uint32((1 << 31) - 1 - (1<<31)%uint32(n))
	v := p.Next()
	for v > max {
		v = p.Next()
	}
	return int(v % uint32(n))
}
