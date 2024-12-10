// Copyright 2020 Consensys Software Inc.
// Licensed under the Apache License, Version 2.0. See the LICENSE file for details.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package bls24317

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bls24-317/fp"
	"github.com/consensys/gnark-crypto/ecc/bls24-317/fr"
	"github.com/consensys/gnark-crypto/internal/parallel"
	"math/big"
	"runtime"
)

// G1Affine is a point in affine coordinates (x,y)
type G1Affine struct {
	X, Y fp.Element
}

// G1Jac is a point in Jacobian coordinates (x=X/Z², y=Y/Z³)
type G1Jac struct {
	X, Y, Z fp.Element
}

// g1JacExtended is a point in extended Jacobian coordinates (x=X/ZZ, y=Y/ZZZ, ZZ³=ZZZ²)
type g1JacExtended struct {
	X, Y, ZZ, ZZZ fp.Element
}

// -------------------------------------------------------------------------------------------------
// Affine coordinates

// Set sets p to a in affine coordinates.
func (p *G1Affine) Set(a *G1Affine) *G1Affine {
	p.X, p.Y = a.X, a.Y
	return p
}

// SetInfinity sets p to the infinity point, which is encoded as (0,0).
// N.B.: (0,0) is never on the curve for j=0 curves (Y²=X³+B).
func (p *G1Affine) SetInfinity() *G1Affine {
	p.X.SetZero()
	p.Y.SetZero()
	return p
}

// ScalarMultiplication computes and returns p = [s]a
// where p and a are affine points.
func (p *G1Affine) ScalarMultiplication(a *G1Affine, s *big.Int) *G1Affine {
	var _p G1Jac
	_p.FromAffine(a)
	_p.mulGLV(&_p, s)
	p.FromJacobian(&_p)
	return p
}

// ScalarMultiplicationBase computes and returns p = [s]g
// where g is the affine point generating the prime subgroup.
func (p *G1Affine) ScalarMultiplicationBase(s *big.Int) *G1Affine {
	var _p G1Jac
	_p.mulGLV(&g1Gen, s)
	p.FromJacobian(&_p)
	return p
}

// Add adds two points in affine coordinates.
// It uses the Jacobian addition with a.Z=b.Z=1 and converts the result to affine coordinates.
//
// https://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-mmadd-2007-bl
func (p *G1Affine) Add(a, b *G1Affine) *G1Affine {
	var q G1Jac
	// a is infinity, return b
	if a.IsInfinity() {
		p.Set(b)
		return p
	}
	// b is infinity, return a
	if b.IsInfinity() {
		p.Set(a)
		return p
	}
	if a.X.Equal(&b.X) {
		// if b == a, we double instead
		if a.Y.Equal(&b.Y) {
			q.DoubleMixed(a)
			return p.FromJacobian(&q)
		} else {
			// if b == -a, we return 0
			return p.SetInfinity()
		}
	}
	var H, HH, I, J, r, V fp.Element
	H.Sub(&b.X, &a.X)
	HH.Square(&H)
	I.Double(&HH).Double(&I)
	J.Mul(&H, &I)
	r.Sub(&b.Y, &a.Y)
	r.Double(&r)
	V.Mul(&a.X, &I)
	q.X.Square(&r).
		Sub(&q.X, &J).
		Sub(&q.X, &V).
		Sub(&q.X, &V)
	q.Y.Sub(&V, &q.X).
		Mul(&q.Y, &r)
	J.Mul(&a.Y, &J).Double(&J)
	q.Y.Sub(&q.Y, &J)
	q.Z.Double(&H)

	return p.FromJacobian(&q)
}

// Double doubles a point in affine coordinates.
// It converts the point to Jacobian coordinates, doubles it using Jacobian
// addition with a.Z=1, and converts it back to affine coordinates.
//
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-mdbl-2007-bl
func (p *G1Affine) Double(a *G1Affine) *G1Affine {
	var q G1Jac
	q.FromAffine(a)
	q.DoubleMixed(a)
	p.FromJacobian(&q)
	return p
}

// Sub subtracts two points in affine coordinates.
// It uses a similar approach to Add, but negates the second point before adding.
func (p *G1Affine) Sub(a, b *G1Affine) *G1Affine {
	var bneg G1Affine
	bneg.Neg(b)
	p.Add(a, &bneg)
	return p
}

// Equal tests if two points in affine coordinates are equal.
func (p *G1Affine) Equal(a *G1Affine) bool {
	return p.X.Equal(&a.X) && p.Y.Equal(&a.Y)
}

// Neg sets p to the affine negative point -a = (a.X, -a.Y).
func (p *G1Affine) Neg(a *G1Affine) *G1Affine {
	p.X = a.X
	p.Y.Neg(&a.Y)
	return p
}

// FromJacobian converts a point p1 from Jacobian to affine coordinates.
func (p *G1Affine) FromJacobian(p1 *G1Jac) *G1Affine {

	var a, b fp.Element

	if p1.Z.IsZero() {
		p.X.SetZero()
		p.Y.SetZero()
		return p
	}

	a.Inverse(&p1.Z)
	b.Square(&a)
	p.X.Mul(&p1.X, &b)
	p.Y.Mul(&p1.Y, &b).Mul(&p.Y, &a)

	return p
}

// String returns the string representation E(x,y) of the affine point p or "O" if it is infinity.
func (p *G1Affine) String() string {
	if p.IsInfinity() {
		return "O"
	}
	return "E([" + p.X.String() + "," + p.Y.String() + "])"
}

// IsInfinity checks if the affine point p is infinity, which is encoded as (0,0).
// N.B.: (0,0) is never on the curve for j=0 curves (Y²=X³+B).
func (p *G1Affine) IsInfinity() bool {
	return p.X.IsZero() && p.Y.IsZero()
}

// IsOnCurve returns true if the affine point p in on the curve.
func (p *G1Affine) IsOnCurve() bool {
	var point G1Jac
	point.FromAffine(p)
	return point.IsOnCurve() // call this function to handle infinity point
}

// IsInSubGroup returns true if the affine point p is in the correct subgroup, false otherwise.
func (p *G1Affine) IsInSubGroup() bool {
	var _p G1Jac
	_p.FromAffine(p)
	return _p.IsInSubGroup()
}

// -------------------------------------------------------------------------------------------------
// Jacobian coordinates

// Set sets p to a in Jacobian coordinates.
func (p *G1Jac) Set(q *G1Jac) *G1Jac {
	p.X, p.Y, p.Z = q.X, q.Y, q.Z
	return p
}

// Equal tests if two points in Jacobian coordinates are equal.
func (p *G1Jac) Equal(q *G1Jac) bool {
	// If one point is infinity, the other must also be infinity.
	if p.Z.IsZero() {
		return q.Z.IsZero()
	}
	// If the other point is infinity, return false since we can't
	// the following checks would be incorrect.
	if q.Z.IsZero() {
		return false
	}

	var pZSquare, aZSquare fp.Element
	pZSquare.Square(&p.Z)
	aZSquare.Square(&q.Z)

	var lhs, rhs fp.Element
	lhs.Mul(&p.X, &aZSquare)
	rhs.Mul(&q.X, &pZSquare)
	if !lhs.Equal(&rhs) {
		return false
	}
	lhs.Mul(&p.Y, &aZSquare).Mul(&lhs, &q.Z)
	rhs.Mul(&q.Y, &pZSquare).Mul(&rhs, &p.Z)

	return lhs.Equal(&rhs)
}

// Neg sets p to the Jacobian negative point -q = (q.X, -q.Y, q.Z).
func (p *G1Jac) Neg(q *G1Jac) *G1Jac {
	*p = *q
	p.Y.Neg(&q.Y)
	return p
}

// AddAssign sets p to p+a in Jacobian coordinates.
//
// https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#addition-add-2007-bl
func (p *G1Jac) AddAssign(q *G1Jac) *G1Jac {

	// p is infinity, return q
	if p.Z.IsZero() {
		p.Set(q)
		return p
	}

	// q is infinity, return p
	if q.Z.IsZero() {
		return p
	}

	var Z1Z1, Z2Z2, U1, U2, S1, S2, H, I, J, r, V fp.Element
	Z1Z1.Square(&q.Z)
	Z2Z2.Square(&p.Z)
	U1.Mul(&q.X, &Z2Z2)
	U2.Mul(&p.X, &Z1Z1)
	S1.Mul(&q.Y, &p.Z).
		Mul(&S1, &Z2Z2)
	S2.Mul(&p.Y, &q.Z).
		Mul(&S2, &Z1Z1)

	// if p == q, we double instead
	if U1.Equal(&U2) && S1.Equal(&S2) {
		return p.DoubleAssign()
	}

	H.Sub(&U2, &U1)
	I.Double(&H).
		Square(&I)
	J.Mul(&H, &I)
	r.Sub(&S2, &S1).Double(&r)
	V.Mul(&U1, &I)
	p.X.Square(&r).
		Sub(&p.X, &J).
		Sub(&p.X, &V).
		Sub(&p.X, &V)
	p.Y.Sub(&V, &p.X).
		Mul(&p.Y, &r)
	S1.Mul(&S1, &J).Double(&S1)
	p.Y.Sub(&p.Y, &S1)
	p.Z.Add(&p.Z, &q.Z)
	p.Z.Square(&p.Z).
		Sub(&p.Z, &Z1Z1).
		Sub(&p.Z, &Z2Z2).
		Mul(&p.Z, &H)

	return p
}

// SubAssign sets p to p-a in Jacobian coordinates.
// It uses a similar approach to AddAssign, but negates the point a before adding.
func (p *G1Jac) SubAssign(q *G1Jac) *G1Jac {
	var tmp G1Jac
	tmp.Set(q)
	tmp.Y.Neg(&tmp.Y)
	p.AddAssign(&tmp)
	return p
}

// Double sets p to [2]q in Jacobian coordinates.
//
// https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#doubling-dbl-2007-bl
func (p *G1Jac) DoubleMixed(a *G1Affine) *G1Jac {
	var XX, YY, YYYY, S, M, T fp.Element
	XX.Square(&a.X)
	YY.Square(&a.Y)
	YYYY.Square(&YY)
	S.Add(&a.X, &YY).
		Square(&S).
		Sub(&S, &XX).
		Sub(&S, &YYYY).
		Double(&S)
	M.Double(&XX).
		Add(&M, &XX) // -> + A, but A=0 here
	T.Square(&M).
		Sub(&T, &S).
		Sub(&T, &S)
	p.X.Set(&T)
	p.Y.Sub(&S, &T).
		Mul(&p.Y, &M)
	YYYY.Double(&YYYY).
		Double(&YYYY).
		Double(&YYYY)
	p.Y.Sub(&p.Y, &YYYY)
	p.Z.Double(&a.Y)

	return p
}

// AddMixed sets p to p+a in Jacobian coordinates, where a.Z = 1.
//
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-madd-2007-bl
func (p *G1Jac) AddMixed(a *G1Affine) *G1Jac {

	//if a is infinity return p
	if a.IsInfinity() {
		return p
	}
	// p is infinity, return a
	if p.Z.IsZero() {
		p.X = a.X
		p.Y = a.Y
		p.Z.SetOne()
		return p
	}

	var Z1Z1, U2, S2, H, HH, I, J, r, V fp.Element
	Z1Z1.Square(&p.Z)
	U2.Mul(&a.X, &Z1Z1)
	S2.Mul(&a.Y, &p.Z).
		Mul(&S2, &Z1Z1)

	// if p == a, we double instead
	if U2.Equal(&p.X) && S2.Equal(&p.Y) {
		return p.DoubleMixed(a)
	}

	H.Sub(&U2, &p.X)
	HH.Square(&H)
	I.Double(&HH).Double(&I)
	J.Mul(&H, &I)
	r.Sub(&S2, &p.Y).Double(&r)
	V.Mul(&p.X, &I)
	p.X.Square(&r).
		Sub(&p.X, &J).
		Sub(&p.X, &V).
		Sub(&p.X, &V)
	J.Mul(&J, &p.Y).Double(&J)
	p.Y.Sub(&V, &p.X).
		Mul(&p.Y, &r)
	p.Y.Sub(&p.Y, &J)
	p.Z.Add(&p.Z, &H)
	p.Z.Square(&p.Z).
		Sub(&p.Z, &Z1Z1).
		Sub(&p.Z, &HH)

	return p
}

// Double sets p to [2]q in Jacobian coordinates.
//
// https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#doubling-dbl-2007-bl
func (p *G1Jac) Double(q *G1Jac) *G1Jac {
	p.Set(q)
	p.DoubleAssign()
	return p
}

// DoubleAssign doubles p in Jacobian coordinates.
//
// https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#doubling-dbl-2007-bl
func (p *G1Jac) DoubleAssign() *G1Jac {

	var XX, YY, YYYY, ZZ, S, M, T fp.Element

	XX.Square(&p.X)
	YY.Square(&p.Y)
	YYYY.Square(&YY)
	ZZ.Square(&p.Z)
	S.Add(&p.X, &YY)
	S.Square(&S).
		Sub(&S, &XX).
		Sub(&S, &YYYY).
		Double(&S)
	M.Double(&XX).Add(&M, &XX)
	p.Z.Add(&p.Z, &p.Y).
		Square(&p.Z).
		Sub(&p.Z, &YY).
		Sub(&p.Z, &ZZ)
	T.Square(&M)
	p.X = T
	T.Double(&S)
	p.X.Sub(&p.X, &T)
	p.Y.Sub(&S, &p.X).
		Mul(&p.Y, &M)
	YYYY.Double(&YYYY).Double(&YYYY).Double(&YYYY)
	p.Y.Sub(&p.Y, &YYYY)

	return p
}

// ScalarMultiplication computes and returns p = [s]a
// where p and a are Jacobian points.
// using the GLV technique.
// see https://www.iacr.org/archive/crypto2001/21390189.pdf
func (p *G1Jac) ScalarMultiplication(q *G1Jac, s *big.Int) *G1Jac {
	return p.mulGLV(q, s)
}

// ScalarMultiplicationBase computes and returns p = [s]g
// where g is the prime subgroup generator.
func (p *G1Jac) ScalarMultiplicationBase(s *big.Int) *G1Jac {
	return p.mulGLV(&g1Gen, s)

}

// String converts p to affine coordinates and returns its string representation E(x,y) or "O" if it is infinity.
func (p *G1Jac) String() string {
	_p := G1Affine{}
	_p.FromJacobian(p)
	return _p.String()
}

// FromAffine converts a point a from affine to Jacobian coordinates.
func (p *G1Jac) FromAffine(a *G1Affine) *G1Jac {
	if a.IsInfinity() {
		p.Z.SetZero()
		p.X.SetOne()
		p.Y.SetOne()
		return p
	}
	p.Z.SetOne()
	p.X.Set(&a.X)
	p.Y.Set(&a.Y)
	return p
}

// IsOnCurve returns true if the Jacobian point p in on the curve.
func (p *G1Jac) IsOnCurve() bool {
	var left, right, tmp, ZZ fp.Element
	left.Square(&p.Y)
	right.Square(&p.X).Mul(&right, &p.X)
	ZZ.Square(&p.Z)
	tmp.Square(&ZZ).Mul(&tmp, &ZZ)
	// Mul tmp by bCurveCoeff=4
	tmp.Double(&tmp).Double(&tmp)
	right.Add(&right, &tmp)
	return left.Equal(&right)
}

// IsInSubGroup returns true if p is on the r-torsion, false otherwise.

// Z[r,0]+Z[-lambdaG1Affine, 1] is the kernel
// of (u,v)->u+lambdaG1Affinev mod r. Expressing r, lambdaG1Affine as
// polynomials in x, a short vector of this Zmodule is
// 1, x⁴. So we check that p+x⁴ϕ(p)
// is the infinity.
func (p *G1Jac) IsInSubGroup() bool {

	var res G1Jac
	res.phi(p).
		mulWindowed(&res, &xGen).
		mulWindowed(&res, &xGen).
		mulWindowed(&res, &xGen).
		mulWindowed(&res, &xGen).
		AddAssign(p)

	return res.IsOnCurve() && res.Z.IsZero()

}

// mulWindowed computes the 2-bits windowed double-and-add scalar
// multiplication p=[s]q in Jacobian coordinates.
func (p *G1Jac) mulWindowed(q *G1Jac, s *big.Int) *G1Jac {

	var res G1Jac
	var ops [3]G1Jac

	ops[0].Set(q)
	if s.Sign() == -1 {
		ops[0].Neg(&ops[0])
	}
	res.Set(&g1Infinity)
	ops[1].Double(&ops[0])
	ops[2].Set(&ops[0]).AddAssign(&ops[1])

	b := s.Bytes()
	for i := range b {
		w := b[i]
		mask := byte(0xc0)
		for j := 0; j < 4; j++ {
			res.DoubleAssign().DoubleAssign()
			c := (w & mask) >> (6 - 2*j)
			if c != 0 {
				res.AddAssign(&ops[c-1])
			}
			mask = mask >> 2
		}
	}
	p.Set(&res)

	return p

}

// phi sets p to ϕ(a) where ϕ: (x,y) → (w x,y),
// where w is a third root of unity.
func (p *G1Jac) phi(q *G1Jac) *G1Jac {
	p.Set(q)
	p.X.Mul(&p.X, &thirdRootOneG1)
	return p
}

// mulGLV computes the scalar multiplication using a windowed-GLV method
//
// see https://www.iacr.org/archive/crypto2001/21390189.pdf
func (p *G1Jac) mulGLV(q *G1Jac, s *big.Int) *G1Jac {

	var table [15]G1Jac
	var res G1Jac
	var k1, k2 fr.Element

	res.Set(&g1Infinity)

	// table[b3b2b1b0-1] = b3b2 ⋅ ϕ(q) + b1b0*q
	table[0].Set(q)
	table[3].phi(q)

	// split the scalar, modifies ±q, ϕ(q) accordingly
	k := ecc.SplitScalar(s, &glvBasis)

	if k[0].Sign() == -1 {
		k[0].Neg(&k[0])
		table[0].Neg(&table[0])
	}
	if k[1].Sign() == -1 {
		k[1].Neg(&k[1])
		table[3].Neg(&table[3])
	}

	// precompute table (2 bits sliding window)
	// table[b3b2b1b0-1] = b3b2 ⋅ ϕ(q) + b1b0 ⋅ q if b3b2b1b0 != 0
	table[1].Double(&table[0])
	table[2].Set(&table[1]).AddAssign(&table[0])
	table[4].Set(&table[3]).AddAssign(&table[0])
	table[5].Set(&table[3]).AddAssign(&table[1])
	table[6].Set(&table[3]).AddAssign(&table[2])
	table[7].Double(&table[3])
	table[8].Set(&table[7]).AddAssign(&table[0])
	table[9].Set(&table[7]).AddAssign(&table[1])
	table[10].Set(&table[7]).AddAssign(&table[2])
	table[11].Set(&table[7]).AddAssign(&table[3])
	table[12].Set(&table[11]).AddAssign(&table[0])
	table[13].Set(&table[11]).AddAssign(&table[1])
	table[14].Set(&table[11]).AddAssign(&table[2])

	// bounds on the lattice base vectors guarantee that k1, k2 are len(r)/2 or len(r)/2+1 bits long max
	// this is because we use a probabilistic scalar decomposition that replaces a division by a right-shift
	k1 = k1.SetBigInt(&k[0]).Bits()
	k2 = k2.SetBigInt(&k[1]).Bits()

	// we don't target constant-timeness so we check first if we increase the bounds or not
	maxBit := k1.BitLen()
	if k2.BitLen() > maxBit {
		maxBit = k2.BitLen()
	}
	hiWordIndex := (maxBit - 1) / 64

	// loop starts from len(k1)/2 or len(k1)/2+1 due to the bounds
	for i := hiWordIndex; i >= 0; i-- {
		mask := uint64(3) << 62
		for j := 0; j < 32; j++ {
			res.Double(&res).Double(&res)
			b1 := (k1[i] & mask) >> (62 - 2*j)
			b2 := (k2[i] & mask) >> (62 - 2*j)
			if b1|b2 != 0 {
				s := (b2<<2 | b1)
				res.AddAssign(&table[s-1])
			}
			mask = mask >> 2
		}
	}

	p.Set(&res)
	return p
}

// ClearCofactor maps a point in curve to r-torsion
func (p *G1Affine) ClearCofactor(a *G1Affine) *G1Affine {
	var _p G1Jac
	_p.FromAffine(a)
	_p.ClearCofactor(&_p)
	p.FromJacobian(&_p)
	return p
}

// ClearCofactor maps a point in E(Fp) to E(Fp)[r]
func (p *G1Jac) ClearCofactor(q *G1Jac) *G1Jac {
	// cf https://eprint.iacr.org/2019/403.pdf, 5
	var res G1Jac
	res.mulWindowed(q, &xGen).Neg(&res).AddAssign(q)
	p.Set(&res)
	return p

}

// JointScalarMultiplication computes [s1]a1+[s2]a2 using Strauss-Shamir technique
// where a1 and a2 are affine points.
func (p *G1Jac) JointScalarMultiplication(a1, a2 *G1Affine, s1, s2 *big.Int) *G1Jac {

	var res, p1, p2 G1Jac
	res.Set(&g1Infinity)
	p1.FromAffine(a1)
	p2.FromAffine(a2)

	var table [15]G1Jac

	var k1, k2 big.Int
	if s1.Sign() == -1 {
		k1.Neg(s1)
		table[0].Neg(&p1)
	} else {
		k1.Set(s1)
		table[0].Set(&p1)
	}
	if s2.Sign() == -1 {
		k2.Neg(s2)
		table[3].Neg(&p2)
	} else {
		k2.Set(s2)
		table[3].Set(&p2)
	}

	// precompute table (2 bits sliding window)
	table[1].Double(&table[0])
	table[2].Set(&table[1]).AddAssign(&table[0])
	table[4].Set(&table[3]).AddAssign(&table[0])
	table[5].Set(&table[3]).AddAssign(&table[1])
	table[6].Set(&table[3]).AddAssign(&table[2])
	table[7].Double(&table[3])
	table[8].Set(&table[7]).AddAssign(&table[0])
	table[9].Set(&table[7]).AddAssign(&table[1])
	table[10].Set(&table[7]).AddAssign(&table[2])
	table[11].Set(&table[7]).AddAssign(&table[3])
	table[12].Set(&table[11]).AddAssign(&table[0])
	table[13].Set(&table[11]).AddAssign(&table[1])
	table[14].Set(&table[11]).AddAssign(&table[2])

	var s [2]fr.Element
	s[0] = s[0].SetBigInt(&k1).Bits()
	s[1] = s[1].SetBigInt(&k2).Bits()

	maxBit := k1.BitLen()
	if k2.BitLen() > maxBit {
		maxBit = k2.BitLen()
	}
	hiWordIndex := (maxBit - 1) / 64

	for i := hiWordIndex; i >= 0; i-- {
		mask := uint64(3) << 62
		for j := 0; j < 32; j++ {
			res.Double(&res).Double(&res)
			b1 := (s[0][i] & mask) >> (62 - 2*j)
			b2 := (s[1][i] & mask) >> (62 - 2*j)
			if b1|b2 != 0 {
				s := (b2<<2 | b1)
				res.AddAssign(&table[s-1])
			}
			mask = mask >> 2
		}
	}

	p.Set(&res)
	return p

}

// JointScalarMultiplicationBase computes [s1]g+[s2]a using Straus-Shamir technique
// where g is the prime subgroup generator.
func (p *G1Jac) JointScalarMultiplicationBase(a *G1Affine, s1, s2 *big.Int) *G1Jac {
	return p.JointScalarMultiplication(&g1GenAff, a, s1, s2)

}

// -------------------------------------------------------------------------------------------------
// extended Jacobian coordinates

// Set sets p to a in extended Jacobian coordinates.
func (p *g1JacExtended) Set(q *g1JacExtended) *g1JacExtended {
	p.X, p.Y, p.ZZ, p.ZZZ = q.X, q.Y, q.ZZ, q.ZZZ
	return p
}

// SetInfinity sets p to the infinity point (1,1,0,0).
func (p *g1JacExtended) SetInfinity() *g1JacExtended {
	p.X.SetOne()
	p.Y.SetOne()
	p.ZZ = fp.Element{}
	p.ZZZ = fp.Element{}
	return p
}

// IsInfinity checks if the p is infinity, i.e. p.ZZ=0.
func (p *g1JacExtended) IsInfinity() bool {
	return p.ZZ.IsZero()
}

// fromJacExtended converts an extended Jacobian point to an affine point.
func (p *G1Affine) fromJacExtended(q *g1JacExtended) *G1Affine {
	if q.ZZ.IsZero() {
		p.X = fp.Element{}
		p.Y = fp.Element{}
		return p
	}
	p.X.Inverse(&q.ZZ).Mul(&p.X, &q.X)
	p.Y.Inverse(&q.ZZZ).Mul(&p.Y, &q.Y)
	return p
}

// fromJacExtended converts an extended Jacobian point to a Jacobian point.
func (p *G1Jac) fromJacExtended(q *g1JacExtended) *G1Jac {
	if q.ZZ.IsZero() {
		p.Set(&g1Infinity)
		return p
	}
	p.X.Mul(&q.ZZ, &q.X).Mul(&p.X, &q.ZZ)
	p.Y.Mul(&q.ZZZ, &q.Y).Mul(&p.Y, &q.ZZZ)
	p.Z.Set(&q.ZZZ)
	return p
}

// unsafeFromJacExtended converts an extended Jacobian point, distinct from Infinity, to a Jacobian point.
func (p *G1Jac) unsafeFromJacExtended(q *g1JacExtended) *G1Jac {
	p.X.Square(&q.ZZ).Mul(&p.X, &q.X)
	p.Y.Square(&q.ZZZ).Mul(&p.Y, &q.Y)
	p.Z = q.ZZZ
	return p
}

// add sets p to p+q in extended Jacobian coordinates.
//
// https://www.hyperelliptic.org/EFD/g1p/auto-shortw-xyzz.html#addition-add-2008-s
func (p *g1JacExtended) add(q *g1JacExtended) *g1JacExtended {
	//if q is infinity return p
	if q.ZZ.IsZero() {
		return p
	}
	// p is infinity, return q
	if p.ZZ.IsZero() {
		p.Set(q)
		return p
	}

	var A, B, U1, U2, S1, S2 fp.Element

	// p2: q, p1: p
	U2.Mul(&q.X, &p.ZZ)
	U1.Mul(&p.X, &q.ZZ)
	A.Sub(&U2, &U1)
	S2.Mul(&q.Y, &p.ZZZ)
	S1.Mul(&p.Y, &q.ZZZ)
	B.Sub(&S2, &S1)

	if A.IsZero() {
		if B.IsZero() {
			return p.double(q)

		}
		p.ZZ = fp.Element{}
		p.ZZZ = fp.Element{}
		return p
	}

	var P, R, PP, PPP, Q, V fp.Element
	P.Sub(&U2, &U1)
	R.Sub(&S2, &S1)
	PP.Square(&P)
	PPP.Mul(&P, &PP)
	Q.Mul(&U1, &PP)
	V.Mul(&S1, &PPP)

	p.X.Square(&R).
		Sub(&p.X, &PPP).
		Sub(&p.X, &Q).
		Sub(&p.X, &Q)
	p.Y.Sub(&Q, &p.X).
		Mul(&p.Y, &R).
		Sub(&p.Y, &V)
	p.ZZ.Mul(&p.ZZ, &q.ZZ).
		Mul(&p.ZZ, &PP)
	p.ZZZ.Mul(&p.ZZZ, &q.ZZZ).
		Mul(&p.ZZZ, &PPP)

	return p
}

// double sets p to [2]q in Jacobian extended coordinates.
//
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-xyzz.html#doubling-dbl-2008-s-1
// N.B.: since we consider any point on Z=0 as the point at infinity
// this doubling formula works for infinity points as well.
func (p *g1JacExtended) double(q *g1JacExtended) *g1JacExtended {
	var U, V, W, S, XX, M fp.Element

	U.Double(&q.Y)
	V.Square(&U)
	W.Mul(&U, &V)
	S.Mul(&q.X, &V)
	XX.Square(&q.X)
	M.Double(&XX).
		Add(&M, &XX) // -> + A, but A=0 here
	U.Mul(&W, &q.Y)

	p.X.Square(&M).
		Sub(&p.X, &S).
		Sub(&p.X, &S)
	p.Y.Sub(&S, &p.X).
		Mul(&p.Y, &M).
		Sub(&p.Y, &U)
	p.ZZ.Mul(&V, &q.ZZ)
	p.ZZZ.Mul(&W, &q.ZZZ)

	return p
}

// addMixed sets p to p+q in extended Jacobian coordinates, where a.ZZ=1.
//
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-xyzz.html#addition-madd-2008-s
func (p *g1JacExtended) addMixed(a *G1Affine) *g1JacExtended {

	//if a is infinity return p
	if a.IsInfinity() {
		return p
	}
	// p is infinity, return a
	if p.ZZ.IsZero() {
		p.X = a.X
		p.Y = a.Y
		p.ZZ.SetOne()
		p.ZZZ.SetOne()
		return p
	}

	var P, R fp.Element

	// p2: a, p1: p
	P.Mul(&a.X, &p.ZZ)
	P.Sub(&P, &p.X)

	R.Mul(&a.Y, &p.ZZZ)
	R.Sub(&R, &p.Y)

	if P.IsZero() {
		if R.IsZero() {
			return p.doubleMixed(a)

		}
		p.ZZ = fp.Element{}
		p.ZZZ = fp.Element{}
		return p
	}

	var PP, PPP, Q, Q2, RR, X3, Y3 fp.Element

	PP.Square(&P)
	PPP.Mul(&P, &PP)
	Q.Mul(&p.X, &PP)
	RR.Square(&R)
	X3.Sub(&RR, &PPP)
	Q2.Double(&Q)
	p.X.Sub(&X3, &Q2)
	Y3.Sub(&Q, &p.X).Mul(&Y3, &R)
	R.Mul(&p.Y, &PPP)
	p.Y.Sub(&Y3, &R)
	p.ZZ.Mul(&p.ZZ, &PP)
	p.ZZZ.Mul(&p.ZZZ, &PPP)

	return p

}

// subMixed works the same as addMixed, but negates a.Y.
//
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-xyzz.html#addition-madd-2008-s
func (p *g1JacExtended) subMixed(a *G1Affine) *g1JacExtended {

	//if a is infinity return p
	if a.IsInfinity() {
		return p
	}
	// p is infinity, return a
	if p.ZZ.IsZero() {
		p.X = a.X
		p.Y.Neg(&a.Y)
		p.ZZ.SetOne()
		p.ZZZ.SetOne()
		return p
	}

	var P, R fp.Element

	// p2: a, p1: p
	P.Mul(&a.X, &p.ZZ)
	P.Sub(&P, &p.X)

	R.Mul(&a.Y, &p.ZZZ)
	R.Neg(&R)
	R.Sub(&R, &p.Y)

	if P.IsZero() {
		if R.IsZero() {
			return p.doubleNegMixed(a)

		}
		p.ZZ = fp.Element{}
		p.ZZZ = fp.Element{}
		return p
	}

	var PP, PPP, Q, Q2, RR, X3, Y3 fp.Element

	PP.Square(&P)
	PPP.Mul(&P, &PP)
	Q.Mul(&p.X, &PP)
	RR.Square(&R)
	X3.Sub(&RR, &PPP)
	Q2.Double(&Q)
	p.X.Sub(&X3, &Q2)
	Y3.Sub(&Q, &p.X).Mul(&Y3, &R)
	R.Mul(&p.Y, &PPP)
	p.Y.Sub(&Y3, &R)
	p.ZZ.Mul(&p.ZZ, &PP)
	p.ZZZ.Mul(&p.ZZZ, &PPP)

	return p

}

// doubleNegMixed works the same as double, but negates q.Y.
func (p *g1JacExtended) doubleNegMixed(a *G1Affine) *g1JacExtended {

	var U, V, W, S, XX, M, S2, L fp.Element

	U.Double(&a.Y)
	U.Neg(&U)
	V.Square(&U)
	W.Mul(&U, &V)
	S.Mul(&a.X, &V)
	XX.Square(&a.X)
	M.Double(&XX).
		Add(&M, &XX) // -> + A, but A=0 here
	S2.Double(&S)
	L.Mul(&W, &a.Y)

	p.X.Square(&M).
		Sub(&p.X, &S2)
	p.Y.Sub(&S, &p.X).
		Mul(&p.Y, &M).
		Add(&p.Y, &L)
	p.ZZ.Set(&V)
	p.ZZZ.Set(&W)

	return p
}

// doubleMixed sets p to [2]a in Jacobian extended coordinates, where a.ZZ=1.
//
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-xyzz.html#doubling-dbl-2008-s-1
func (p *g1JacExtended) doubleMixed(a *G1Affine) *g1JacExtended {

	var U, V, W, S, XX, M, S2, L fp.Element

	U.Double(&a.Y)
	V.Square(&U)
	W.Mul(&U, &V)
	S.Mul(&a.X, &V)
	XX.Square(&a.X)
	M.Double(&XX).
		Add(&M, &XX) // -> + A, but A=0 here
	S2.Double(&S)
	L.Mul(&W, &a.Y)

	p.X.Square(&M).
		Sub(&p.X, &S2)
	p.Y.Sub(&S, &p.X).
		Mul(&p.Y, &M).
		Sub(&p.Y, &L)
	p.ZZ.Set(&V)
	p.ZZZ.Set(&W)

	return p
}

// BatchJacobianToAffineG1 converts points in Jacobian coordinates to Affine coordinates
// performing a single field inversion using the Montgomery batch inversion trick.
func BatchJacobianToAffineG1(points []G1Jac) []G1Affine {
	result := make([]G1Affine, len(points))
	zeroes := make([]bool, len(points))
	accumulator := fp.One()

	// batch invert all points[].Z coordinates with Montgomery batch inversion trick
	// (stores points[].Z^-1 in result[i].X to avoid allocating a slice of fr.Elements)
	for i := 0; i < len(points); i++ {
		if points[i].Z.IsZero() {
			zeroes[i] = true
			continue
		}
		result[i].X = accumulator
		accumulator.Mul(&accumulator, &points[i].Z)
	}

	var accInverse fp.Element
	accInverse.Inverse(&accumulator)

	for i := len(points) - 1; i >= 0; i-- {
		if zeroes[i] {
			// do nothing, (X=0, Y=0) is infinity point in affine
			continue
		}
		result[i].X.Mul(&result[i].X, &accInverse)
		accInverse.Mul(&accInverse, &points[i].Z)
	}

	// batch convert to affine.
	parallel.Execute(len(points), func(start, end int) {
		for i := start; i < end; i++ {
			if zeroes[i] {
				// do nothing, (X=0, Y=0) is infinity point in affine
				continue
			}
			var a, b fp.Element
			a = result[i].X
			b.Square(&a)
			result[i].X.Mul(&points[i].X, &b)
			result[i].Y.Mul(&points[i].Y, &b).
				Mul(&result[i].Y, &a)
		}
	})

	return result
}

// BatchScalarMultiplicationG1 multiplies the same base by all scalars
// and return resulting points in affine coordinates
// uses a simple windowed-NAF-like multiplication algorithm.
func BatchScalarMultiplicationG1(base *G1Affine, scalars []fr.Element) []G1Affine {
	// approximate cost in group ops is
	// cost = 2^{c-1} + n(scalar.nbBits+nbChunks)

	nbPoints := uint64(len(scalars))
	min := ^uint64(0)
	bestC := 0
	for c := 2; c <= 16; c++ {
		cost := uint64(1 << (c - 1)) // pre compute the table
		nbChunks := computeNbChunks(uint64(c))
		cost += nbPoints * (uint64(c) + 1) * nbChunks // doublings + point add
		if cost < min {
			min = cost
			bestC = c
		}
	}
	c := uint64(bestC) // window size
	nbChunks := int(computeNbChunks(c))

	// last window may be slightly larger than c; in which case we need to compute one
	// extra element in the baseTable
	maxC := lastC(c)
	if c > maxC {
		maxC = c
	}

	// precompute all powers of base for our window
	// note here that if performance is critical, we can implement as in the msmX methods
	// this allocation to be on the stack
	baseTable := make([]G1Jac, (1 << (maxC - 1)))
	baseTable[0].FromAffine(base)
	for i := 1; i < len(baseTable); i++ {
		baseTable[i] = baseTable[i-1]
		baseTable[i].AddMixed(base)
	}
	// convert our base exp table into affine to use AddMixed
	baseTableAff := BatchJacobianToAffineG1(baseTable)
	toReturn := make([]G1Jac, len(scalars))

	// partition the scalars into digits
	digits, _ := partitionScalars(scalars, c, runtime.NumCPU())

	// for each digit, take value in the base table, double it c time, voilà.
	parallel.Execute(len(scalars), func(start, end int) {
		var p G1Jac
		for i := start; i < end; i++ {
			p.Set(&g1Infinity)
			for chunk := nbChunks - 1; chunk >= 0; chunk-- {
				if chunk != nbChunks-1 {
					for j := uint64(0); j < c; j++ {
						p.DoubleAssign()
					}
				}
				offset := chunk * len(scalars)
				digit := digits[i+offset]

				if digit == 0 {
					continue
				}

				// if msbWindow bit is set, we need to subtract
				if digit&1 == 0 {
					// add
					p.AddMixed(&baseTableAff[(digit>>1)-1])
				} else {
					// sub
					t := baseTableAff[digit>>1]
					t.Neg(&t)
					p.AddMixed(&t)
				}
			}

			// set our result point
			toReturn[i] = p

		}
	})
	toReturnAff := BatchJacobianToAffineG1(toReturn)
	return toReturnAff
}

// batchAddG1Affine adds affine points using the Montgomery batch inversion trick.
// Special cases (doubling, infinity) must be filtered out before this call.
func batchAddG1Affine[TP pG1Affine, TPP ppG1Affine, TC cG1Affine](R *TPP, P *TP, batchSize int) {
	var lambda, lambdain TC

	// from https://docs.zkproof.org/pages/standards/accepted-workshop3/proposal-turbo_plonk.pdf
	// affine point addition formula
	// R(X1, Y1) + P(X2, Y2) = Q(X3, Y3)
	// λ  = (Y2 - Y1) / (X2 - X1)
	// X3 = λ² - (X1 + X2)
	// Y3 = λ * (X1 - X3) - Y1

	// first we compute the 1 / (X2 - X1) for all points using Montgomery batch inversion trick

	// X2 - X1
	for j := 0; j < batchSize; j++ {
		lambdain[j].Sub(&(*P)[j].X, &(*R)[j].X)
	}

	// montgomery batch inversion;
	// lambda[0] = 1 / (P[0].X - R[0].X)
	// lambda[1] = 1 / (P[1].X - R[1].X)
	// ...
	{
		var accumulator fp.Element
		lambda[0].SetOne()
		accumulator.Set(&lambdain[0])

		for i := 1; i < batchSize; i++ {
			lambda[i] = accumulator
			accumulator.Mul(&accumulator, &lambdain[i])
		}

		accumulator.Inverse(&accumulator)

		for i := batchSize - 1; i > 0; i-- {
			lambda[i].Mul(&lambda[i], &accumulator)
			accumulator.Mul(&accumulator, &lambdain[i])
		}
		lambda[0].Set(&accumulator)
	}

	var t fp.Element
	var Q G1Affine

	for j := 0; j < batchSize; j++ {
		// λ  = (Y2 - Y1) / (X2 - X1)
		t.Sub(&(*P)[j].Y, &(*R)[j].Y)
		lambda[j].Mul(&lambda[j], &t)

		// X3 = λ² - (X1 + X2)
		Q.X.Square(&lambda[j])
		Q.X.Sub(&Q.X, &(*R)[j].X)
		Q.X.Sub(&Q.X, &(*P)[j].X)

		// Y3 = λ * (X1 - X3) - Y1
		t.Sub(&(*R)[j].X, &Q.X)
		Q.Y.Mul(&lambda[j], &t)
		Q.Y.Sub(&Q.Y, &(*R)[j].Y)

		(*R)[j].Set(&Q)
	}
}
