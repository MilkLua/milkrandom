// Package splitmix64 implements the SplitMix64 random number generator.
package splitmix64

import (
	"encoding/binary"
	"errors"
	"sync"
	"time"
)

// SplitMix64 represents the state of a SplitMix64 random number generator.
type SplitMix64 struct {
	state uint64
}

// SafeSplitMix64 represents the state of a SplitMix64 random number generator with a mutex to make it safe for concurrent use.
type SafeSplitMix64 struct {
	SplitMix64
	mu sync.Mutex
}

// New creates a new SplitMix64 instance seeded with the current time.
func New() *SplitMix64 {
	x := &SplitMix64{}
	x.Seed(uint64(time.Now().UnixNano()))
	return x
}

// NewSafe creates a new safe SplitMix64 instance seeded with the current time.
func NewSafe() *SafeSplitMix64 {
	x := &SafeSplitMix64{}
	x.Seed(uint64(time.Now().UnixNano()))
	return x
}

// State returns the current state of the random number generator.
func (x *SplitMix64) State() uint64 {
	return x.state
}

// State returns the current state of the random number generator, which is safe for concurrent use.
func (x *SafeSplitMix64) State() uint64 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.SplitMix64.State()
}

// Reset resets the state of the random number generator to the seed value.
func (x *SplitMix64) Reset() {
	x.Seed(uint64(time.Now().UnixNano()))
}

// Reset resets the state of the random number generator to the seed value, which is safe for concurrent use.
func (x *SafeSplitMix64) Reset() {
	x.mu.Lock()
	defer x.mu.Unlock()
	x.SplitMix64.Reset()
}

// Marshal returns the binary encoding of the current state of the random number generator.
func (x *SplitMix64) Marshal() ([]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, x.state)
	return buf, nil
}

// Marshal returns the binary encoding of the current state of the random number generator, which is safe for concurrent use.
func (x *SafeSplitMix64) Marshal() ([]byte, error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.SplitMix64.Marshal()
}

// Unmarshal sets the state of the random number generator to the state represented by the input data.
func (x *SplitMix64) Unmarshal(data []byte) error {
	if len(data) != 8 {
		return errors.New("splitmix64: invalid state length")
	}
	x.state = binary.LittleEndian.Uint64(data)
	return nil
}

// Unmarshal sets the state of the random number generator to the state represented by the input data, which is safe for concurrent use.
func (x *SafeSplitMix64) Unmarshal(data []byte) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.SplitMix64.Unmarshal(data)
}

// Seed initializes the state of the random number generator with the given seed value.
func (x *SplitMix64) Seed(seed uint64) {
	x.state = seed
}

// Seed initializes the state of the random number generator with the given seed value, which is safe for concurrent use.
func (x *SafeSplitMix64) Seed(seed uint64) {
	x.mu.Lock()
	defer x.mu.Unlock()
	x.SplitMix64.Seed(seed)
}

// Uint64 generates a random 64-bit unsigned integer.
func (x *SplitMix64) Uint64() uint64 {
	x.state += 0x9e3779b97f4a7c15
	z := x.state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

// Uint64 generates a random 64-bit unsigned integer, which is safe for concurrent use.
func (x *SafeSplitMix64) Uint64() uint64 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.SplitMix64.Uint64()
}

// Int64 generates a random 64-bit signed integer.
func (x *SplitMix64) Int64() int64 {
	return int64(x.Uint64() >> 1)
}

// Int64 generates a random 64-bit signed integer, which is safe for concurrent use.
func (x *SafeSplitMix64) Int64() int64 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.SplitMix64.Int64()
}

// Uint32 generates a random 32-bit unsigned integer.
func (x *SplitMix64) Uint32() uint32 {
	return uint32(x.Uint64() >> 32)
}

// Uint32 generates a random 32-bit unsigned integer, which is safe for concurrent use.
func (x *SafeSplitMix64) Uint32() uint32 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.SplitMix64.Uint32()
}

// Int32 generates a random 32-bit signed integer.
func (x *SplitMix64) Int32() int32 {
	return int32(x.Uint32() >> 1)
}

// Int32 generates a random 32-bit signed integer, which is safe for concurrent use.
func (x *SafeSplitMix64) Int32() int32 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.SplitMix64.Int32()
}

// Int generates a random integer in the range [0, n).
func (x *SplitMix64) Int(n int) int {
	if n <= 0 {
		panic("splitmix64: argument to Int is <= 0")
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
func (x *SafeSplitMix64) Int(n int) int {
	if n <= 0 {
		panic("splitmix64: argument to Int is <= 0")
	}
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.SplitMix64.Int(n)
}

// Float64 generates a random float64 in the range [0.0, 1.0).
func (x *SplitMix64) Float64() float64 {
	return float64(x.Uint64()>>(64-53)) / (1 << 53)
}

// Float64 generates a random float64 in the range [0.0, 1.0), which is safe for concurrent use.
func (x *SafeSplitMix64) Float64() float64 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.SplitMix64.Float64()
}

// Float32 generates a random float32 in the range [0.0, 1.0).
func (x *SplitMix64) Float32() float32 {
	return float32(x.Uint32()>>(32-24)) / (1 << 24)
}

// Float32 generates a random float32 in the range [0.0, 1.0), which is safe for concurrent use.
func (x *SafeSplitMix64) Float32() float32 {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.SplitMix64.Float32()
}
