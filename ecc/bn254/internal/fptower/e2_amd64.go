// Copyright 2020 Consensys Software Inc.
// Licensed under the Apache License, Version 2.0. See the LICENSE file for details.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package fptower

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
)

// q + r'.r = 1, i.e., qInvNeg = - q⁻¹ mod r
// used for Montgomery reduction
const qInvNeg uint64 = 9786893198990664585

// Field modulus q (Fp)
const (
	q0 uint64 = 4332616871279656263
	q1 uint64 = 10917124144477883021
	q2 uint64 = 13281191951274694749
	q3 uint64 = 3486998266802970665
)

var qElement = fp.Element{
	q0,
	q1,
	q2,
	q3,
}

//go:noescape
func addE2(res, x, y *E2)

//go:noescape
func subE2(res, x, y *E2)

//go:noescape
func doubleE2(res, x *E2)

//go:noescape
func negE2(res, x *E2)

//go:noescape
func mulNonResE2(res, x *E2)

//go:noescape
func squareAdxE2(res, x *E2)

//go:noescape
func mulAdxE2(res, x, y *E2)

// MulByNonResidue multiplies a E2 by (9,1)
func (z *E2) MulByNonResidue(x *E2) *E2 {
	mulNonResE2(z, x)
	return z
}

// Mul sets z to the E2-product of x,y, returns z
func (z *E2) Mul(x, y *E2) *E2 {
	mulAdxE2(z, x, y)
	return z
}

// Square sets z to the E2-product of x,x, returns z
func (z *E2) Square(x *E2) *E2 {
	squareAdxE2(z, x)
	return z
}
