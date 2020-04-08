// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by goff (v0.2.0) DO NOT EDIT

// Package bn256 contains field arithmetic operations
package bn256

import (
	"math/bits"

	"golang.org/x/sys/cpu"
)

var supportAdx = cpu.X86.HasADX && cpu.X86.HasBMI2

func madd(a, b, t, u, v uint64) (uint64, uint64, uint64) {
	var carry uint64
	hi, lo := bits.Mul64(a, b)
	v, carry = bits.Add64(lo, v, 0)
	u, carry = bits.Add64(hi, u, carry)
	t, _ = bits.Add64(t, 0, carry)
	return t, u, v
}

// madd0 hi = a*b + c (discards lo bits)
func madd0(a, b, c uint64) (hi uint64) {
	var carry, lo uint64
	hi, lo = bits.Mul64(a, b)
	_, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}

// madd1 hi, lo = a*b + c
func madd1(a, b, c uint64) (hi uint64, lo uint64) {
	var carry uint64
	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}

// madd2 hi, lo = a*b + c + d
func madd2(a, b, c, d uint64) (hi uint64, lo uint64) {
	var carry uint64
	hi, lo = bits.Mul64(a, b)
	c, carry = bits.Add64(c, d, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}

// madd2s superhi, hi, lo = 2*a*b + c + d + e
func madd2s(a, b, c, d, e uint64) (superhi, hi, lo uint64) {
	var carry, sum uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)

	sum, carry = bits.Add64(c, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, sum, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	hi, _ = bits.Add64(hi, 0, d)
	return
}

func madd1s(a, b, d, e uint64) (superhi, hi, lo uint64) {
	var carry uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)
	lo, carry = bits.Add64(lo, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	hi, _ = bits.Add64(hi, 0, d)
	return
}

func madd2sb(a, b, c, e uint64) (superhi, hi, lo uint64) {
	var carry, sum uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)

	sum, carry = bits.Add64(c, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, sum, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}

func madd1sb(a, b, e uint64) (superhi, hi, lo uint64) {
	var carry uint64

	hi, lo = bits.Mul64(a, b)
	lo, carry = bits.Add64(lo, lo, 0)
	hi, superhi = bits.Add64(hi, hi, carry)
	lo, carry = bits.Add64(lo, e, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	return
}

func madd3(a, b, c, d, e uint64) (hi uint64, lo uint64) {
	var carry uint64
	hi, lo = bits.Mul64(a, b)
	c, carry = bits.Add64(c, d, 0)
	hi, _ = bits.Add64(hi, 0, carry)
	lo, carry = bits.Add64(lo, c, 0)
	hi, _ = bits.Add64(hi, e, carry)
	return
}
