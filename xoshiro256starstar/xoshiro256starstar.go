// Package xoshiro256starstar implements the Xoshiro256** random number generator.
package xoshiro256starstar

import (
	"encoding/binary"
	"errors"
	"math/bits"
	"sync"
	"time"
)

// Xoshiro256StarStar represents the state of a xoshiro256** random number generator.
type Xoshiro256StarStar struct {
	state [4]uint64
}

// SafeXoshiro256StarStar represents the state of a xoshiro256** random number generator with a mutex to make it safe for concurrent use.
type SafeXoshiro256StarStar struct {
	Xoshiro256StarStar
	mu sync.Mutex
}

// New creates a new xoshiro256StarStar instance seeded with the current time.
func New() *Xoshiro256StarStar {
	x := &Xoshiro256StarStar{}
	x.Seed(uint64(time.Now().UnixNano()))
	// "Warm up" the generator
	for i := 0; i < 10; i++ {
		x.Uint64()
	}
	return x
}

// NewSafe creates a new safe xoshiro256StarStar instance seeded with the current time.
func NewSafe() *SafeXoshiro256StarStar {
	x := &SafeXoshiro256StarStar{}
	x.Seed(uint64(time.Now().UnixNano()))
	// "Warm up" the generator
	for i := 0; i < 10; i++ {
		x.Uint64()
	}
	return x
}

// State returns the current state of the random number generator.
func (x *Xoshiro256StarStar) State() [4]uint64 {
	return x.state
}

// State returns the current state of the random number generator, which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) State() [4]uint64 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.Xoshiro256StarStar.State()
}

// Reset resets the state of the random number generator to the seed value.
func (x *Xoshiro256StarStar) Reset() {
	x.Seed(uint64(time.Now().UnixNano()))
}

// Reset resets the state of the random number generator to the seed value, which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Reset() {
	x.mu.Lock()
	defer x.mu.Unlock()
	x.Xoshiro256StarStar.Reset()
}

// MarshalBinary returns the binary encoding of the current state of the random number generator.
func (x *Xoshiro256StarStar) Marshal() ([]byte, error) {
	buf := make([]byte, 32)
	for i, v := range x.state {
		binary.LittleEndian.PutUint64(buf[i*8:], v)
	}
	return buf, nil
}

// MarshalBinary returns the binary encoding of the current state of the random number generator, which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Marshal() ([]byte, error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.Xoshiro256StarStar.Marshal()
}

// UnmarshalBinary sets the state of the random number generator to the state represented by the input data.
func (x *Xoshiro256StarStar) Unmarshal(data []byte) error {
	if len(data) != 32 {
		return errors.New("xoshiro256starstar: invalid state length")
	}
	for i := range x.state {
		x.state[i] = binary.LittleEndian.Uint64(data[i*8:])
	}
	return nil
}

// UnmarshalBinary sets the state of the random number generator to the state represented by the input data, which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Unmarshal(data []byte) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.Xoshiro256StarStar.Unmarshal(data)
}

// Seed initializes the state of the random number generator with the given seed value.
func (x *Xoshiro256StarStar) Seed(seed uint64) {
	if seed == 0 { // Seed with current time if seed is 0
		seed = uint64(time.Now().UnixNano())
	}
	splitmix64 := func(s *uint64) uint64 {
		*s += 0x9e3779b97f4a7c15
		z := *s
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		return z ^ (z >> 31)
	}
	s := seed
	x.state[0] = splitmix64(&s)
	x.state[1] = splitmix64(&s)
	x.state[2] = splitmix64(&s)
	x.state[3] = splitmix64(&s)
}

// Seed initializes the state of the random number generator with the given seed value, which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Seed(seed uint64) {
	x.mu.Lock()
	defer x.mu.Unlock()
	x.Xoshiro256StarStar.Seed(seed)
}

// Uint64 generates a random 64-bit unsigned integer.
func (x *Xoshiro256StarStar) Uint64() uint64 {
	result := bits.RotateLeft64(x.state[1]*5, 7) * 9
	t := x.state[1] << 17
	x.state[2] ^= x.state[0]
	x.state[3] ^= x.state[1]
	x.state[1] ^= x.state[2]
	x.state[0] ^= x.state[3]
	x.state[2] ^= t
	x.state[3] = bits.RotateLeft64(x.state[3], 45)
	return result
}

// Uint64 generates a random 64-bit unsigned integer, which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Uint64() uint64 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.Xoshiro256StarStar.Uint64()
}

// Jump advances the internal state by 2^128 calls to Next().
func (x *Xoshiro256StarStar) Jump() {
	jump := [4]uint64{0x180ec6d33cfd0aba, 0xd5a61266f0c9392c, 0xa9582618e03fc9aa, 0x39abdc4529b1661c}
	s0, s1, s2, s3 := uint64(0), uint64(0), uint64(0), uint64(0)
	for i := 0; i < len(jump); i++ {
		for b := uint64(0); b < 64; b++ {
			if (jump[i] & (1 << b)) != 0 {
				s0 ^= x.state[0]
				s1 ^= x.state[1]
				s2 ^= x.state[2]
				s3 ^= x.state[3]
			}
			x.Uint64()
		}
	}
	x.state[0], x.state[1], x.state[2], x.state[3] = s0, s1, s2, s3
}

// Jump advances the internal state by 2^128 calls to Next(), which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Jump() {
	x.mu.Lock()
	defer x.mu.Unlock()
	x.Xoshiro256StarStar.Jump()
}

// LongJump advances the internal state by 2^192 calls to Next().
func (x *Xoshiro256StarStar) LongJump() {
	jump := [4]uint64{0x76e15d3efefdcbbf, 0xc5004e441c522fb3, 0x77710069854ee241, 0x39109bb02acbe635}
	s0, s1, s2, s3 := uint64(0), uint64(0), uint64(0), uint64(0)
	for i := 0; i < len(jump); i++ {
		for b := uint64(0); b < 64; b++ {
			if (jump[i] & (1 << b)) != 0 {
				s0 ^= x.state[0]
				s1 ^= x.state[1]
				s2 ^= x.state[2]
				s3 ^= x.state[3]
			}
			x.Uint64()
		}
	}
	x.state[0], x.state[1], x.state[2], x.state[3] = s0, s1, s2, s3
}

// LongJump advances the internal state by 2^192 calls to Next(), which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) LongJump() {
	x.mu.Lock()
	defer x.mu.Unlock()
	x.Xoshiro256StarStar.LongJump()
}

// go:inline
// Int64 generates a random 64-bit signed integer.
func (x *Xoshiro256StarStar) Int64() int64 {
	return int64(x.Uint64() >> 1)
}

// go:inline
// Int64 generates a random 64-bit signed integer, which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Int64() int64 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return int64(x.Xoshiro256StarStar.Uint64() >> 1)
}

// go:inline
// Uint32 generates a random 32-bit unsigned integer.
func (x *Xoshiro256StarStar) Uint32() uint32 {
	return uint32(x.Uint64() >> 32)
}

// go:inline
// Uint32 generates a random 32-bit unsigned integer, which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Uint32() uint32 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return uint32(x.Xoshiro256StarStar.Uint64() >> 32)
}

// go:inline
// Int32 generates a random 32-bit signed integer.
func (x *Xoshiro256StarStar) Int32() int32 {
	return int32(x.Uint32() >> 1)
}

// go:inline
// Int32 generates a random 32-bit signed integer, which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Int32() int32 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return int32(x.Xoshiro256StarStar.Uint32() >> 1)
}

// Int generates a random integer in the range [0, n).
func (x *Xoshiro256StarStar) Int(n int) int {
	if n <= 0 {
		panic("xoshiro256starstar: argument to Int is <= 0")
	}
	if n&(n-1) == 0 { // n is 2^m, use mask
		return int(x.Uint64() & uint64(n-1))
	}
	max := uint64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := x.Uint64()
	for v > max {
		v = x.Uint64()
	}
	return int(v % uint64(n))
}

// Int generates a random integer in the range [0, n), which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Int(n int) int {
	if n <= 0 {
		panic("xoshiro256starstar: argument to Int is <= 0")
	}
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.Xoshiro256StarStar.Int(n)
}

// go:inline
// Float64 generates a random float64 in the range [0.0, 1.0).
func (x *Xoshiro256StarStar) Float64() float64 {
	return float64(x.Uint64()>>(64-53)) / (1 << 53)
}

// go:inline
// Float64 generates a random float64 in the range [0.0, 1.0), which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Float64() float64 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return float64(x.Xoshiro256StarStar.Uint64()>>(64-53)) / (1 << 53)
}

// go:inline
// Float32 generates a random float32 in the range [0.0, 1.0).
func (x *Xoshiro256StarStar) Float32() float32 {
	return float32(x.Uint32()>>(32-24)) / (1 << 24)
}

// go:inline
// Float32 generates a random float32 in the range [0.0, 1.0), which is safe for concurrent use.
func (x *SafeXoshiro256StarStar) Float32() float32 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return float32(x.Xoshiro256StarStar.Uint32()>>(32-24)) / (1 << 24)
}
