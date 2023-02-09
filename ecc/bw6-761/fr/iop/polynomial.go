// Copyright 2020 ConsenSys Software Inc.
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

// Code generated by consensys/gnark-crypto DO NOT EDIT

package iop

import (
	"math/big"
	"math/bits"

	"github.com/consensys/gnark-crypto/ecc/bw6-761/fr"
	"github.com/consensys/gnark-crypto/ecc/bw6-761/fr/fft"
)

// Basis indicates the basis in which a polynomial is represented.
type Basis uint32

const (
	Canonical Basis = 1 << iota
	Lagrange
	LagrangeCoset
)

// Layout indicates if a polynomial has a BitReverse or a Regular layout
type Layout uint32

const (
	Regular Layout = 8 << iota
	BitReverse
)

// Form describes the form of a polynomial.
// TODO should be a regular enum?
type Form struct {
	Basis  Basis
	Layout Layout
}

// enum of the possible Form values for type-safe switches
// in this package
var (
	canonicalRegular        = Form{Canonical, Regular}
	canonicalBitReverse     = Form{Canonical, BitReverse}
	lagrangeRegular         = Form{Lagrange, Regular}
	lagrangeBitReverse      = Form{Lagrange, BitReverse}
	lagrangeCosetRegular    = Form{LagrangeCoset, Regular}
	lagrangeCosetBitReverse = Form{LagrangeCoset, BitReverse}
)

// Polynomial represents a polynomial, the vector of coefficients
// along with the basis and the layout.
type Polynomial struct {
	Coefficients []fr.Element
	Form
}

// NewPolynomial creates a new polynomial. The slice coeff NOT copied
// but directly assigned to the new polynomial.
func NewPolynomial(coeffs []fr.Element, form Form) *Polynomial {
	return &Polynomial{Coefficients: coeffs, Form: form}
}

// Clone returns a deep copy of the underlying data structure.
func (p *Polynomial) Clone() *Polynomial {
	r := &Polynomial{
		Coefficients: make([]fr.Element, len(p.Coefficients)),
		Form:         p.Form,
	}
	copy(r.Coefficients, p.Coefficients)
	return r
}

// Evaluate evaluates p at x.
// The code panics if the function is not in canonical form.
func (p *Polynomial) Evaluate(x fr.Element) fr.Element {

	var r fr.Element
	if p.Basis != Canonical {
		panic("p must be in canonical basis")
	}

	if p.Layout == Regular {
		for i := len(p.Coefficients) - 1; i >= 0; i-- {
			r.Mul(&r, &x).Add(&r, &p.Coefficients[i])
		}
	} else {
		nn := uint64(64 - bits.TrailingZeros(uint(len(p.Coefficients))))
		for i := len(p.Coefficients) - 1; i >= 0; i-- {
			iRev := bits.Reverse64(uint64(i)) >> nn
			r.Mul(&r, &x).Add(&r, &p.Coefficients[iRev])
		}
	}

	return r

}

// ToRegular changes the layout of p to Regular.
// Leaves p unchanged if p's layout was already Regular.
func (p *Polynomial) ToRegular() *Polynomial {
	if p.Layout == Regular {
		return p
	}
	fft.BitReverse(p.Coefficients)
	p.Layout = Regular
	return p
}

// ToBitReverse changes the layout of p to BitReverse.
// Leaves p unchanged if p's layout was already BitReverse.
func (p *Polynomial) ToBitReverse() *Polynomial {
	if p.Layout == BitReverse {
		return p
	}
	fft.BitReverse(p.Coefficients)
	p.Layout = BitReverse
	return p
}

// ToLagrange converts p to Lagrange form.
// Leaves p unchanged if p was already in Lagrange form.
func (p *Polynomial) ToLagrange(d *fft.Domain) *Polynomial {
	id := p.Form
	resize(p, d.Cardinality)
	switch id {
	case canonicalRegular:
		p.Layout = BitReverse
		d.FFT(p.Coefficients, fft.DIF)
	case canonicalBitReverse:
		p.Layout = Regular
		d.FFT(p.Coefficients, fft.DIT)
	case lagrangeRegular, lagrangeBitReverse:
		return p
	case lagrangeCosetRegular:
		p.Layout = Regular
		d.FFTInverse(p.Coefficients, fft.DIF, true)
		d.FFT(p.Coefficients, fft.DIT)
	case lagrangeCosetBitReverse:
		p.Layout = BitReverse
		d.FFTInverse(p.Coefficients, fft.DIT, true)
		d.FFT(p.Coefficients, fft.DIF)
	default:
		panic("unknown ID")
	}
	p.Basis = Lagrange
	return p
}

// ToCanonical converts p to canonical form.
// Leaves p unchanged if p was already in Canonical form.
func (p *Polynomial) ToCanonical(d *fft.Domain) *Polynomial {
	id := p.Form
	resize(p, d.Cardinality)
	switch id {
	case canonicalRegular, canonicalBitReverse:
		return p
	case lagrangeRegular:
		p.Layout = BitReverse
		d.FFTInverse(p.Coefficients, fft.DIF)
	case lagrangeBitReverse:
		p.Layout = Regular
		d.FFTInverse(p.Coefficients, fft.DIT)
	case lagrangeCosetRegular:
		p.Layout = BitReverse
		d.FFTInverse(p.Coefficients, fft.DIF, true)
	case lagrangeCosetBitReverse:
		p.Layout = Regular
		d.FFTInverse(p.Coefficients, fft.DIT, true)
	default:
		panic("unknown ID")
	}
	p.Basis = Canonical
	return p
}

func resize(p *Polynomial, newSize uint64) {
	z := make([]fr.Element, int(newSize)-len(p.Coefficients))
	p.Coefficients = append(p.Coefficients, z...)
}

// ToLagrangeCoset Sets p to q, in LagrangeCoset form and returns it.
func (p *Polynomial) ToLagrangeCoset(d *fft.Domain) *Polynomial {
	id := p.Form
	resize(p, d.Cardinality)
	switch id {
	case canonicalRegular:
		p.Layout = BitReverse
		d.FFT(p.Coefficients, fft.DIF, true)
	case canonicalBitReverse:
		p.Layout = Regular
		d.FFT(p.Coefficients, fft.DIT, true)
	case lagrangeRegular:
		p.Layout = Regular
		d.FFTInverse(p.Coefficients, fft.DIF)
		d.FFT(p.Coefficients, fft.DIT, true)
	case lagrangeBitReverse:
		p.Layout = BitReverse
		d.FFTInverse(p.Coefficients, fft.DIT)
		d.FFT(p.Coefficients, fft.DIF, true)
	case lagrangeCosetRegular, lagrangeCosetBitReverse:
		return p
	default:
		panic("unknown ID")
	}

	p.Basis = LagrangeCoset
	return p
}

// WrappedPolynomial wraps a polynomial so that it is
// interpreted as P'(X)=P(\omega^{s}X).
// Size is the real size of the polynomial (seen as a vector).
// For instance if len(P)=32 but P.Size=8, it means that P has been
// extended (e.g. it is evaluated on a larger set) but P is a polynomial
// of degree 7.
// blindedSize is the size of the polynomial when it is blinded. By
// default blindedSize=Size, until the polynomial is blinded.
type WrappedPolynomial struct {
	*Polynomial
	shift       int
	size        int
	blindedSize int
}

// NewWrappedPolynomial returned a WrappedPolynomial from p.
// ! Warning this does not do a deep copy of p, and modifications on the wrapped
// polynomial will modify the underlying coefficients of p.
func NewWrappedPolynomial(p *Polynomial) *WrappedPolynomial {
	return &WrappedPolynomial{
		Polynomial:  p,
		size:        len(p.Coefficients),
		blindedSize: len(p.Coefficients),
	}
}

// Shift the wrapped polynomial; it doesn't modify the underlying data structure,
// but flag the WrappedPolynomial such that it will be interpreted as p(\omega^shift X)
func (wp *WrappedPolynomial) Shift(shift int) *WrappedPolynomial {
	wp.shift = shift
	return wp
}

// BlindedSize returns the the size of the polynomial when it is blinded. By
// default blindedSize=Size, until the polynomial is blinded.
func (wp *WrappedPolynomial) BlindedSize() int {
	return wp.blindedSize
}

// Blind blinds a polynomial q by adding Q(X)*(X^{n}-1),
// where deg Q = blindingOrder and Q is random, and n is the
// size of q. Sets the result to p and returns it.
//
// blindingOrder is the degree of Q, where the blinding is Q(X)*(X^{n}-1)
// where n is the size of wp. The size of wp is modified since the underlying
// polynomial is of bigger degree now. The new size is wp.Size+1+blindingOrder.
//
// /!\ The code panics if wq is not in canonical, regular layout
func (wp *WrappedPolynomial) Blind(blindingOrder int) *WrappedPolynomial {
	// check that wp is in canonical basis
	if wp.Form != canonicalRegular {
		panic("the input must be in canonical basis, regular layout")
	}

	// we add Q*(x^{n}-1) so the new size is deg(Q)+n+1
	// where n is the size of wq.
	newSize := wp.size + blindingOrder + 1

	// Resize wp. The size of wq might has already been increased
	// (e.g. when the polynomial is evaluated on a larger domain),
	// if that's the case we don't resize the polynomial.
	offset := newSize - len(wp.Coefficients)
	if offset > 0 {
		z := make([]fr.Element, offset)
		wp.Coefficients = append(wp.Coefficients, z...)
	}

	// blinding: we add Q(X)(X^{n}-1) to P, where deg(Q)=blindingOrder
	var r fr.Element

	for i := 0; i <= blindingOrder; i++ {
		r.SetRandom()
		wp.Coefficients[i].Sub(&wp.Coefficients[i], &r)
		wp.Coefficients[i+wp.size].Add(&wp.Coefficients[i+wp.size], &r)
	}
	wp.blindedSize = newSize

	return wp
}

// Evaluate evaluates p at x.
// The code panics if the function is not in canonical form.
func (wp *WrappedPolynomial) Evaluate(x fr.Element) fr.Element {

	if wp.shift == 0 {
		return wp.Polynomial.Evaluate(x)
	}

	// TODO find a way to retrieve the root properly instead of re generating the fft domain
	d := fft.NewDomain(uint64(wp.size))
	var g fr.Element
	if wp.shift <= 5 {
		g = smallExp(d.Generator, wp.shift)
		x.Mul(&x, &g)
		return wp.Polynomial.Evaluate(x)
	}

	bs := big.NewInt(int64(wp.shift))
	g = *g.Exp(g, bs)
	x.Mul(&x, &g)
	return wp.Polynomial.Evaluate(x)
}

// Clone returns a deep copy of wp. The underlying polynomial is cloned;
// see also ShallowClone to perform a ShallowClone on the underlying polynomial.
func (wp *WrappedPolynomial) Clone() *WrappedPolynomial {
	res := wp.ShallowClone()
	res.Polynomial = wp.Polynomial.Clone()
	return res
}

// ShallowClone returns a shallow copy of wp. The underlying polynomial coefficient
// is NOT cloned and both objects will point to the same coefficient vector.
func (wp *WrappedPolynomial) ShallowClone() *WrappedPolynomial {
	res := *wp
	return &res
}

// GetCoeff returns the i-th entry of wp, taking the layout in account.
func (wp *WrappedPolynomial) GetCoeff(i int) fr.Element {

	n := len(wp.Coefficients)
	rho := n / wp.size
	if wp.Polynomial.Form.Layout == Regular {
		return wp.Coefficients[(i+rho*wp.shift)%n]
	} else {
		nn := uint64(64 - bits.TrailingZeros(uint(n)))
		iRev := bits.Reverse64(uint64((i+rho*wp.shift)%n)) >> nn
		return wp.Coefficients[iRev]
	}

}
