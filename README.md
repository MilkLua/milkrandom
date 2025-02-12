# MilkRandom

A prng library for MilkLua in pure Go.

## `./pcg32`

The PCG-32 (Permuted Congruential Generator) is a high-quality, fast random number generator designed by Melissa O'Neill. It offers excellent statistical properties and a good balance between speed and quality for non-cryptographic applications.

#### **Key Features**
- **Algorithm**: Combines a linear congruential generator with output permutation.
- **State Size**: 64 bits for the state, 64 bits for the increment.
- **Period**: $2^{64}$, providing a long sequence of unique outputs.
- **Performance**: Fast and efficient, suitable for various applications.
- **Output**: Generates 32-bit random numbers.

#### **Implementation Details**
1. **State Initialization**:
   - Uses two 64-bit integers: `state` and `inc`.
   - Can be seeded with any 64-bit unsigned integer.

2. **Random Number Generation**:
   - Applies a linear congruential step followed by a permutation operation.
   - Outputs are uniformly distributed 32-bit unsigned integers.

3. **Seeding Process**:
   - Initializes both `state` and `inc` based on the provided seed.
   - Performs a few steps to ensure proper initialization.

4. **Additional Functions**:
   - Provides methods for generating floating-point numbers and integers within specific ranges.

## `./pcg64`

The PCG-64 (Permuted Congruential Generator) is an extended version of PCG-32, offering a larger state space and output. It provides high-quality random numbers suitable for demanding applications requiring 64-bit outputs.

#### **Key Features**
- **Algorithm**: Uses a 128-bit linear congruential generator with output permutation.
- **State Size**: 128 bits for the state, 128 bits for the increment.
- **Period**: $2^{128}$, providing an extremely long sequence of unique outputs.
- **Performance**: Efficient for 64-bit systems while maintaining high statistical quality.
- **Output**: Generates 64-bit random numbers.

#### **Implementation Details**
1. **State Representation**:
   - Uses two 128-bit integers: `state` and `inc`, each represented by a pair of 64-bit values.
   - Implements a custom `uint128` type to handle 128-bit arithmetic.

2. **Random Number Generation**:
   - Applies a 128-bit linear congruential step followed by a permutation operation.
   - Outputs uniformly distributed 64-bit unsigned integers.

3. **Seeding Process**:
   - Initializes both `state` and `inc` based on the provided seed.
   - Performs multiple steps to ensure proper initialization of the 128-bit state.

4. **Concurrency Support**:
   - Provides a thread-safe version (`SafePCG64`) using mutexes for concurrent access.

5. **Additional Functions**:
   - Includes methods for generating various types of random numbers (float64, int64, etc.).
   - Implements custom 128-bit arithmetic operations (add128, mul128).

## `./splitmix64`

The `SplitMix64` random number generator is a simple and fast PRNG designed primarily for seeding other generators, such as `xoshiro256**`. It has excellent statistical properties and a long period, making it suitable for standalone use in non-cryptographic applications.

---

#### **Key Features**
- **Algorithm**: Uses simple arithmetic and bitwise operations to generate random numbers.
- **State Size**: 64 bits (a single unsigned 64-bit integer).
- **Performance**: Extremely fast and lightweight, ideal for initializing other PRNGs.
- **Applications**: Commonly used as a seed generator for more complex PRNGs like `xoshiro256**`.
- **Period**: $2^{64}$, ensuring a sufficiently long sequence of unique outputs.

---

#### **Implementation Details**
1. **State Initialization**:
   - The generator maintains a single 64-bit state.
   - It can be seeded with any 64-bit unsigned integer.

2. **Random Number Generation**:
   - The algorithm adds a constant to the state, applies bitwise shifts and XOR operations, and performs multiplications to produce a new random value.
   - Outputs are uniformly distributed over the range of 64-bit unsigned integers.

3. **Concurrency Support**:
   - A thread-safe version (`SafeSplitMix64`) uses mutexes to ensure safe concurrent access.

## `./xoshiro256starstar`

The `xoshiro256**` random number generator is a high-performance, general-purpose pseudorandom number generator (PRNG) with excellent statistical properties. It was designed by David Blackman and Sebastiano Vigna and is widely implemented in various programming languages. Below are the key features and implementation details:

---

#### **Key Features**
- **Algorithm**: Combines XOR, shift, and rotate operations for generating random numbers.
- **State Size**: 256 bits (four 64-bit unsigned integers).
- **Performance**: Optimized for sub-nanosecond speed, making it suitable for high-performance applications.
- **Applications**: Ideal for simulations, games, and other non-cryptographic use cases.
- **Parallelization**: Supports "jump" and "long jump" operations to create non-overlapping subsequences for parallel computations.

---

#### **Implementation Details**
1. **State Initialization**:
   - The state consists of four 64-bit integers.
   - It can be seeded using methods like `SplitMix64` or random byte arrays.

2. **Random Number Generation**:
   - The core function generates a 64-bit unsigned integer by performing bitwise operations on the internal state.
   - Variants exist to produce signed integers, 32-bit numbers, or floating-point values.

3. **Jump Functions**:
   - `Jump`: Advances the state by $2^{128}$ steps, creating independent streams of randomness.
   - `LongJump`: Advances the state by $2^{192}$ steps for even greater separation of streams.

4. **Concurrency Support**:
   - Thread-safe implementations use mutexes to ensure safe concurrent access to the generator.