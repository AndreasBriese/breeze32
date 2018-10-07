// breeze32.go

// The MIT License (MIT)

// Copyright (c) 2014 Andreas Briese, eduToolbox@Bri-C GmbH, Sarstedt GERMANY

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package breeze32

import (
	"math"
	"sync"
	"time"
)

// Breeze32
//
// implements a cb-prng with four LM
// seeds with 64bit (uint64)
// 32 Byte outputState
type Breeze32 struct {
	State  [4]uint64
	State1 float64
	State2 float64
	State3 float64
	State4 float64
	idx    uint8
	mutex  sync.Mutex
}

var (
	maskUint8 = [8]uint64{
		0xff, 0xff00,
		0xff0000, 0xff000000,
		0xff00000000, 0xff0000000000,
		0xff000000000000, 0xff00000000000000,
	}
)

const (
	f64DIV = float64(1 << 64)
	f53DIV = float64(1 << 53)
	f32DIV = float64(1 << 32)
	f23DIV = float32(1 << 23)
)

// Reset resets to the initial (empty) state
// before initializing.
func (l *Breeze32) Reset() {
	*l = Breeze32{}
}

// Init initializes from user input by calling initr() to process the input to become seeds (seedr(seed)) for the LMs.
// Init reseeds the LMs but it does NOT reset the prng:
//    it seeds based on the previous output states, internal bitshift and idx values
func (l *Breeze32) Init(s uint64) *Breeze32 {
	l.seedr([2]uint64{s, s ^ 0xffffffffffffffff})
	return l
}

// seedr calculates the startvalues of the LMs and
// calls for the initial 'startrounds' roundtrips to shift circle
// once or more times over the output states
func (l *Breeze32) seedr(seed [2]uint64) {
	splittr := func(seed uint64) (s1, s2, s3 uint64) {
		s1 = (1 << 22) - seed>>43
		s2 = (1 << 23) - seed<<21>>42
		s3 = (1 << 22) - seed<<43>>43
		return s1, s2, s3
	}

	s1, s2, s3 := splittr(seed[0])
	s1, s2, s3 = splittr(s1 ^ s2 ^ s3 | seed[1])

	l.State1 = 1.0 / float64(s1)
	l.State2 = 1.0 - 1.0/float64(s2)
	l.State3 = 1.0 / float64(s3)
	l.State4 = 1.0 / float64(s2)

	startrounds := 9
	for startrounds > 0 {
		l.roundTrip()
		startrounds--
	}

}

// checks if LM's are seeded
func (l *Breeze32) isSeeded() bool {
	for _, v := range l.State {
		if v > 0 {
			return true
		}
	}
	return false
}

// roundTrip calculates the next LMs states
// uncomment for tests the states to be != 0 meening exhaustion (else reseeds from previous states)
// interchanges the states between LMs after 'mirroring them at 1' <- Lorenz system in IEEE-754 FP returns highest density near zero
// processes the output states from LMs states
// mixin (xoring) output states
func (l *Breeze32) roundTrip() {
	newState1 := (1.0 - l.State1)
	newState1 *= 4.0 * l.State1
	newState2 := (1.0 - l.State2)
	newState2 *= 3.99 * l.State2
	newState3 := (1.0 - l.State3)
	newState3 *= 3.98 * l.State3
	newState4 := (1.0 - l.State4)
	newState4 *= 3.995 * l.State4
	switch newState1 * newState2 * newState3 * newState4 {
	case 0:
		s1 := math.Float64bits(l.State1) ^ math.Float64bits(l.State2) ^ math.Float64bits(l.State3) ^ math.Float64bits(l.State4)
		l.seedr([2]uint64{s1, s1 ^ 0xffffffffffffffff})
		// panic("LM is gone")
	default:
		l.State1 = 1.0 - newState2
		l.State2 = 1.0 - newState3
		l.State3 = 1.0 - newState4
		l.State4 = 1.0 - newState1
	}

	// l.mutex.Lock()
	// defer l.mutex.Unlock()
	l.idx = 0

	l.State[0] ^= ((math.Float64bits(l.State1) >> 5) << 18)
	l.State[0] ^= ((math.Float64bits(l.State3) >> 5) << 22)
	l.State[0] ^= ((math.Float64bits(l.State2) >> 5) & 0x2fffff)

	l.State[1] ^= ((math.Float64bits(l.State3) >> 5) << 18)
	l.State[1] ^= ((math.Float64bits(l.State2) >> 5) << 22)
	l.State[1] ^= ((math.Float64bits(l.State4) >> 5) & 0x2fffff)

	l.State[2] ^= ((math.Float64bits(l.State2) >> 5) << 18)
	l.State[2] ^= ((math.Float64bits(l.State4) >> 5) << 22)
	l.State[2] ^= ((math.Float64bits(l.State1) >> 5) & 0x2fffff)

	l.State[3] ^= ((math.Float64bits(l.State4) >> 5) << 18)
	l.State[3] ^= ((math.Float64bits(l.State1) >> 5) << 22)
	l.State[3] ^= ((math.Float64bits(l.State3) >> 5) & 0x2fffff)

}

// RandIntN(n int) int
// takes an upper limit n and returns random values <n [0,n)
// conditional autoseeding with time.Now().UnixNano()
func (l *Breeze32) RandIntN(n int) int {
	if !l.isSeeded() {
		l.Init(uint64(time.Now().UnixNano()))
	}

	switch {
	case n <= 1<<32:
		return int(l.uint32N(n + 1))
	}

	return int(l.uint64N(n + 1))
}

// uint64N(n int) (r uint64)
// private func
// takes an upper limit n and returns random values <n [0,n)
func (l *Breeze32) uint64N(n int) (r uint64) {
mining:
	mf := float64(l.State[l.idx>>3]) / f64DIV
	r = uint64(float64(n) * mf)
	l.idx += 8
	if l.idx > 31 {
		l.roundTrip()
	}
	if r == 0 {
		goto mining
	}
	return r - 1
}

// Uint64(r *uint64) uint64
// places a random uint64 in r
// and returns r (for comodity)
// blazing fast <- no further floating point arithmetics
func (l *Breeze32) Uint64(r *uint64) uint64 {
	*r = l.State[l.idx>>3]
	l.idx += 8
	if l.idx > 31 {
		l.roundTrip()
	}
	return *r
}

// uint32N(n int) (r uint32)
// private func
// takes an upper limit n and returns random values <n [0,n)
func (l *Breeze32) uint32N(n int) (r uint32) {
mining:
	// mf := float64((l.State[l.idx>>3]&maskUint32[(l.idx>>2)&1])>>(((l.idx>>2)&1)<<5)) / f32DIV
	mf := float64(l.State[l.idx>>3]>>(((l.idx>>2)&1)<<5)) / f32DIV
	r = uint32(float64(n) * mf)
	l.idx += 4
	if l.idx > 31 {
		l.roundTrip()
	}
	if r == 0 {
		goto mining
	}
	return r - 1
}

// Uint32(r *uint32) uint32
// places a random uint32 in r
// and returns r (for comodity)
// blazing fast <- no further floating point arithmetics
func (l *Breeze32) Uint32(r *uint32) uint32 {
	// *r = uint32(l.State[l.idx>>3]&maskUint32[(l.idx>>2)&1]) >> (((l.idx >> 2) & 1) << 5)
	*r = uint32(l.State[l.idx>>3] >> (((l.idx >> 2) & 1) << 5))
	l.idx += 4
	if l.idx > 31 {
		l.roundTrip()
	}
	return *r
}

// unused -> UintN32 is as fast
// func (l *Breeze32) Uint8N(n int) (r uint8) {
// mining:
// 	mf := float64(l.State[l.idx>>3] & maskUint8[l.idx&7] >> ((l.idx & 7) << 3))
// 	mf /= 256.
// 	r = uint8(float64(n) * mf)
// 	l.idx++
// 	if l.idx > 31 {
// 		l.roundTrip()
// 	}
// 	if r == 0 {
// 		goto mining
// 	}
// 	return r - 1
// }

func (l *Breeze32) Uint8(r *uint8) uint8 {
	*r = uint8(l.State[l.idx>>3] & maskUint8[l.idx&7] >> ((l.idx & 7) << 3))
	l.idx++
	if l.idx > 31 {
		l.roundTrip()
	}
	return *r
}

// Byte(*uint8) returns one Byte  --  not threadsafe
// for compatibility in my code
func (l *Breeze32) Byte(r *uint8) uint8 {
	return l.Uint8(r)
}

// ByteMP(*uint8) returns threadsafe one Byte
func (l *Breeze32) ByteMP(r *uint8) uint8 {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.idx++
	if l.idx > 31 {
		l.roundTrip()
	}
	*r = uint8(l.State[l.idx>>3] & maskUint8[l.idx&7] >> ((l.idx & 7) << 3))
	return *r
}

// Float64() returns a float64 (double precision)
func (l *Breeze32) Float64() float64 {
	r := uint64(1)
	l.Uint64(&r)
	return float64(r>>11) / f53DIV
}

// Float32() returns a float32 (single precision)
func (l *Breeze32) Float32() float32 {
	r := uint32(1)
	l.Uint32(&r)
	return float32(r>>9) / f23DIV
}
