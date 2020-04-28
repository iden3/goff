// Copyright 2019 ConsenSys AG
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

package cmd

import (
	"math/big"
)

type field struct {
	PackageName          string
        IfaceName            string
	ElementName          string
	Modulus              string
	NbWords              int
	NbBits               int
	NbWordsLastIndex     int
	NbWordsIndexesNoZero []int
	NbWordsIndexesFull   []int
	IdxFIPS              []int
	Q                    []uint64
	QInverse             []uint64
	ASM                  bool
	RSquare              []uint64
	One                  []uint64
	LegendreExponent     []uint64
	NoCarry              bool
	NoCarrySquare        bool // used if NoCarry is set, but some op may overflow in square optimization
	Benches              bool
	SqrtQ3Mod4           bool
	SqrtAtkin            bool
	SqrtTonelliShanks    bool
	SqrtE                uint64
	SqrtS                []uint64
	SqrtAtkinExponent    []uint64
	SqrtSMinusOneOver2   []uint64
	SqrtG                []uint64 // NonResidue ^  SqrtR (montgomery form)
	SqrtQ3Mod4Exponent   []uint64
	NonResidue           []uint64 // (montgomery form)
	Version              string
	NoCollidingNames     bool // if multiple elements are generated in the same package, triggers name collisions
}

// -------------------------------------------------------------------------------------------------
// Field data precompute functions
func newField(packageName, elementName, modulus string, benches bool, ifaceName string, noCollidingNames bool) (*field, error) {
	// parse modulus
	var bModulus big.Int
	if _, ok := bModulus.SetString(modulus, 10); !ok {
		return nil, errParseModulus
	}

	// field info
        if ifaceName == "" {
           ifaceName = elementName
        }
	F := &field{
		PackageName:      packageName,
		ElementName:      elementName,
		Modulus:          modulus,
		Benches:          benches,
                IfaceName:        ifaceName,
		NoCollidingNames: noCollidingNames,
	}
	F.Version = Version
	// pre compute field constants
	F.NbBits = bModulus.BitLen()
	F.NbWords = len(bModulus.Bits())
	if F.NbWords < 2 {
		return nil, errUnsupportedModulus
	}

	F.NbWordsLastIndex = F.NbWords - 1

	// set q from big int repr
	F.Q = toUint64Slice(&bModulus)

	//  setting qInverse
	_r := big.NewInt(1)
	_r.Lsh(_r, uint(F.NbWords)*64)
	_rInv := big.NewInt(1)
	_qInv := big.NewInt(0)
	extendedEuclideanAlgo(_r, &bModulus, _rInv, _qInv)
	_qInv.Mod(_qInv, _r)
	F.QInverse = toUint64Slice(_qInv)

	//  rsquare
	_rSquare := big.NewInt(2)
	exponent := big.NewInt(int64(F.NbWords) * 64 * 2)
	_rSquare.Exp(_rSquare, exponent, &bModulus)
	F.RSquare = toUint64Slice(_rSquare)

	var one big.Int
	one.SetUint64(1)
	one.Lsh(&one, uint(F.NbWords)*64).Mod(&one, &bModulus)
	F.One = toUint64Slice(&one)

	// indexes (template helpers)
	F.NbWordsIndexesFull = make([]int, F.NbWords)
	F.NbWordsIndexesNoZero = make([]int, F.NbWords-1)
	for i := 0; i < F.NbWords; i++ {
		F.NbWordsIndexesFull[i] = i
		if i > 0 {
			F.NbWordsIndexesNoZero[i-1] = i
		}
	}

	// See https://hackmd.io/@zkteam/modular_multiplication
	// if the last word of the modulus is smaller or equal to B,
	// we can simplify the montgomery multiplication
	const B = (^uint64(0) >> 1) - 1
	F.NoCarry = (F.Q[len(F.Q)-1] <= B) && F.NbWords <= 12
	const BSquare = (^uint64(0) >> 2)
	F.NoCarrySquare = F.Q[len(F.Q)-1] <= BSquare

	for i := F.NbWords; i <= 2*F.NbWords-2; i++ {
		F.IdxFIPS = append(F.IdxFIPS, i)
	}

	// Legendre exponent (p-1)/2
	var legendreExponent big.Int
	legendreExponent.SetUint64(1)
	legendreExponent.Sub(&bModulus, &legendreExponent)
	legendreExponent.Rsh(&legendreExponent, 1)
	F.LegendreExponent = toUint64Slice(&legendreExponent)

	// Sqrt pre computes
	var qMod big.Int
	qMod.SetUint64(4)
	if qMod.Mod(&bModulus, &qMod).Cmp(new(big.Int).SetUint64(3)) == 0 {
		// q ≡ 3 (mod 4)
		// using  z ≡ ± x^((p+1)/4) (mod q)
		F.SqrtQ3Mod4 = true
		var sqrtExponent big.Int
		sqrtExponent.SetUint64(1)
		sqrtExponent.Add(&bModulus, &sqrtExponent)
		sqrtExponent.Rsh(&sqrtExponent, 2)
		F.SqrtQ3Mod4Exponent = toUint64Slice(&sqrtExponent)
	} else {
		// q ≡ 1 (mod 4)
		qMod.SetUint64(8)
		if qMod.Mod(&bModulus, &qMod).Cmp(new(big.Int).SetUint64(5)) == 0 {
			// q ≡ 5 (mod 8)
			// use Atkin's algorithm
			// see modSqrt5Mod8Prime in math/big/int.go
			F.SqrtAtkin = true
			e := new(big.Int).Rsh(&bModulus, 3) // e = (q - 5) / 8
			F.SqrtAtkinExponent = toUint64Slice(e)
		} else {
			// use Tonelli-Shanks
			F.SqrtTonelliShanks = true

			// Write q-1 =2^e * s , s odd
			var s big.Int
			one.SetUint64(1)
			s.Sub(&bModulus, &one)

			e := s.TrailingZeroBits()
			s.Rsh(&s, e)
			F.SqrtE = uint64(e)
			F.SqrtS = toUint64Slice(&s)

			// find non residue
			var nonResidue big.Int
			nonResidue.SetInt64(2)
			one.SetUint64(1)
			for big.Jacobi(&nonResidue, &bModulus) != -1 {
				nonResidue.Add(&nonResidue, &one)
			}

			// g = nonresidue ^ s
			var g big.Int
			g.Exp(&nonResidue, &s, &bModulus)
			// store g in montgomery form
			g.Lsh(&g, uint(F.NbWords)*64).Mod(&g, &bModulus)
			F.SqrtG = toUint64Slice(&g)

			// store non residue in montgomery form
			nonResidue.Lsh(&nonResidue, uint(F.NbWords)*64).Mod(&nonResidue, &bModulus)
			F.NonResidue = toUint64Slice(&nonResidue)

			// (s+1) /2
			s.Sub(&s, &one).Rsh(&s, 1)
			F.SqrtSMinusOneOver2 = toUint64Slice(&s)
		}
	}

	// ASM
	F.ASM = F.NoCarry && F.NbWords <= 6 // max words without having to deal with spilling

	return F, nil
}

func toUint64Slice(b *big.Int) (s []uint64) {
	s = make([]uint64, len(b.Bits()))
	for i, v := range b.Bits() {
		s[i] = (uint64)(v)
	}
	return
}

// https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
// r > q, modifies rinv and qinv such that rinv.r - qinv.q = 1
func extendedEuclideanAlgo(r, q, rInv, qInv *big.Int) {
	var s1, s2, t1, t2, qi, tmpMuls, riPlusOne, tmpMult, a, b big.Int
	t1.SetUint64(1)
	rInv.Set(big.NewInt(1))
	qInv.Set(big.NewInt(0))
	a.Set(r)
	b.Set(q)

	// r_i+1 = r_i-1 - q_i.r_i
	// s_i+1 = s_i-1 - q_i.s_i
	// t_i+1 = t_i-1 - q_i.s_i
	for b.Sign() > 0 {
		qi.Div(&a, &b)
		riPlusOne.Mod(&a, &b)

		tmpMuls.Mul(&s1, &qi)
		tmpMult.Mul(&t1, &qi)

		s2.Set(&s1)
		t2.Set(&t1)

		s1.Sub(rInv, &tmpMuls)
		t1.Sub(qInv, &tmpMult)
		rInv.Set(&s2)
		qInv.Set(&t2)

		a.Set(&b)
		b.Set(&riPlusOne)
	}
	qInv.Neg(qInv)
}
