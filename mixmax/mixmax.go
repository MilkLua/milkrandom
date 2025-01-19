// Package mixmax implements the MIXMAX random number generator.
package mixmax

import (
	"encoding/binary"
	"errors"
	"sync"
	"time"
)

const (
	N    = 240                // Default matrix size
	MOD  = 0xfffffffffffffffe // 2^61 - 2
	TWOM = 1.0 / (1 << 61)
)

// MIXMAX represents the state of a MIXMAX random number generator.
type MIXMAX struct {
	state   [N]uint64
	counter int
}

// SafeMIXMAX represents the state of a MIXMAX random number generator with a mutex to make it safe for concurrent use.
type SafeMIXMAX struct {
	MIXMAX
	mu sync.Mutex
}

// New creates a new MIXMAX instance seeded with the current time.
func New() *MIXMAX {
	m := &MIXMAX{}
	m.Seed(uint64(time.Now().UnixNano()))
	// "Warm up" the generator
	for i := 0; i < 10; i++ {
		m.Uint64()
	}
	return m
}

// NewSafe creates a new safe MIXMAX instance seeded with the current time.
func NewSafe() *SafeMIXMAX {
	m := &SafeMIXMAX{}
	m.Seed(uint64(time.Now().UnixNano()))
	// "Warm up" the generator
	for i := 0; i < 10; i++ {
		m.Uint64()
	}
	return m
}

// go:inline
// State returns the current state of the random number generator.
func (m *MIXMAX) State() [N]uint64 {
	return m.state
}

// State returns the current state of the random number generator, which is safe for concurrent use.
func (m *SafeMIXMAX) State() [N]uint64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.MIXMAX.State()
}

// Reset resets the state of the random number generator to the seed value.
//
//go:inline
func (m *MIXMAX) Reset() {
	m.Seed(uint64(time.Now().UnixNano()))
}

// go:inline
// Reset resets the state of the random number generator to the seed value, which is safe for concurrent use.
func (m *SafeMIXMAX) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MIXMAX.Reset()
}

// Marshal returns the binary encoding of the current state of the random number generator.
func (m *MIXMAX) Marshal() ([]byte, error) {
	buf := make([]byte, N*8+8)
	for i, v := range m.state {
		binary.LittleEndian.PutUint64(buf[i*8:], v)
	}
	binary.LittleEndian.PutUint64(buf[N*8:], uint64(m.counter))
	return buf, nil
}

// go:inline
// Marshal returns the binary encoding of the current state of the random number generator, which is safe for concurrent use.
func (m *SafeMIXMAX) Marshal() ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.MIXMAX.Marshal()
}

// Unmarshal sets the state of the random number generator to the state represented by the input data.
func (m *MIXMAX) Unmarshal(data []byte) error {
	if len(data) != N*8+8 {
		return errors.New("mixmax: invalid state length")
	}
	for i := range m.state {
		m.state[i] = binary.LittleEndian.Uint64(data[i*8:])
	}
	m.counter = int(binary.LittleEndian.Uint64(data[N*8:]))
	return nil
}

// Unmarshal sets the state of the random number generator to the state represented by the input data, which is safe for concurrent use.
func (m *SafeMIXMAX) Unmarshal(data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.MIXMAX.Unmarshal(data)
}

// Seed initializes the state of the random number generator with the given seed value.
func (m *MIXMAX) Seed(seed uint64) {
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	m.state[0] = seed
	for i := 1; i < N; i++ {
		m.state[i] = (m.state[i-1]*6364136223846793005 + 1) & MOD
	}
	m.counter = 0
}

// go:inline
// Seed initializes the state of the random number generator with the given seed value, which is safe for concurrent use.
func (m *SafeMIXMAX) Seed(seed uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.MIXMAX.Seed(seed)
}

// Uint64 generates a random 64-bit unsigned integer.
func (m *MIXMAX) Uint64() uint64 {
	if m.counter == 0 {
		m.iterate()
	}
	m.counter--
	return m.state[m.counter]
}

// go:inline
// Uint64 generates a random 64-bit unsigned integer, which is safe for concurrent use.
func (m *SafeMIXMAX) Uint64() uint64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.MIXMAX.Uint64()
}

// iterate advances the internal state of the generator.
func (m *MIXMAX) iterate() {
	var t uint64
	for i := 0; i < N; i++ {
		t = m.state[(i+1)%N]
		m.state[i] = (m.state[i] + t) & MOD
		if m.state[i] < t {
			m.state[i]++
		}
	}
	m.counter = N
}

// Float64 generates a random float64 in the range [0.0, 1.0).
//
//go:inline
func (m *MIXMAX) Float64() float64 {
	return float64(m.Uint64()) * TWOM
}

// Float64 generates a random float64 in the range [0.0, 1.0), which is safe for concurrent use.
func (m *SafeMIXMAX) Float64() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.MIXMAX.Float64()
}

// Int64 generates a random 63-bit signed integer.
//
//go:inline
func (m *MIXMAX) Int64() int64 {
	return int64(m.Uint64() >> 1)
}

// Int64 generates a random 63-bit signed integer, which is safe for concurrent use.
//
//go:inline
func (m *SafeMIXMAX) Int64() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.MIXMAX.Int64()
}
