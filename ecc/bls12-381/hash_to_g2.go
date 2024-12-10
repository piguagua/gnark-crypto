// Copyright 2020 Consensys Software Inc.
// Licensed under the Apache License, Version 2.0. See the LICENSE file for details.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package bls12381

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fp"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/internal/fptower"

	"math/big"
)

//Note: This only works for simple extensions

func g2IsogenyXNumerator(dst *fptower.E2, x *fptower.E2) {
	g2EvalPolynomial(dst,
		false,
		[]fptower.E2{
			{
				A0: fp.Element{5185457120960601698, 494647221959407934, 8971396042087821730, 324544954362548322, 14214792730224113654, 1405280679127738945},
				A1: fp.Element{5185457120960601698, 494647221959407934, 8971396042087821730, 324544954362548322, 14214792730224113654, 1405280679127738945},
			},
			{
				A0: fp.Element{0},
				A1: fp.Element{6910023028261548496, 9745789443900091043, 7668299866710145304, 2432656849393633605, 2897729527445498821, 776645607375592125},
			},
			{
				A0: fp.Element{724047465092313539, 15783990863276714670, 12824896677063784855, 15246381572572671516, 13186611051602728692, 1485475813959743803},
				A1: fp.Element{12678383550985550056, 4872894721950045521, 13057521970209848460, 10439700461551592610, 10672236800577525218, 388322803687796062},
			},
			{
				A0: fp.Element{4659755689450087917, 1804066951354704782, 15570919779568036803, 15592734958806855601, 7597208057374167129, 1841438384006890194},
				A1: fp.Element{0},
			},
		},
		x)
}

func g2IsogenyXDenominator(dst *fptower.E2, x *fptower.E2) {
	g2EvalPolynomial(dst,
		true,
		[]fptower.E2{
			{
				A0: fp.Element{0},
				A1: fp.Element{2250392438786206615, 17463829474098544446, 14571211649711714824, 4495761442775821336, 258811604141191305, 357646605018048850},
			},
			{
				A0: fp.Element{4933130441833534766, 15904462746612662304, 8034115857496836953, 12755092135412849606, 7007796720291435703, 252692002104915169},
				A1: fp.Element{8469300574244328829, 4752422838614097887, 17848302789776796362, 12930989898711414520, 16851051131888818207, 1621106615542624696},
			},
		},
		x)
}

func g2IsogenyYNumerator(dst *fptower.E2, x *fptower.E2, y *fptower.E2) {
	var _dst fptower.E2
	g2EvalPolynomial(&_dst,
		false,
		[]fptower.E2{
			{
				A0: fp.Element{10869708750642247614, 13056187057366814946, 1750362034917495549, 6326189602300757217, 1140223926335695785, 632761649765668291},
				A1: fp.Element{10869708750642247614, 13056187057366814946, 1750362034917495549, 6326189602300757217, 1140223926335695785, 632761649765668291},
			},
			{
				A0: fp.Element{0},
				A1: fp.Element{13765940311003083782, 5579209876153186557, 11349908400803699438, 11707848830955952341, 199199289641242246, 899896674917908607},
			},
			{
				A0: fp.Element{15562563812347550836, 2436447360975022760, 6528760985104924230, 5219850230775796305, 5336118400288762609, 194161401843898031},
				A1: fp.Element{16286611277439864375, 18220438224251737430, 906913588459157469, 2019487729638916206, 75985378181939686, 1679637215803641835},
			},
			{
				A0: fp.Element{11849179119594500956, 13906615243538674725, 14543197362847770509, 2041759640812427310, 2879701092679313252, 1259985822978576468},
				A1: fp.Element{0},
			},
		},
		x)

	dst.Mul(&_dst, y)
}

func g2IsogenyYDenominator(dst *fptower.E2, x *fptower.E2) {
	g2EvalPolynomial(dst,
		true,
		[]fptower.E2{
			{
				A0: fp.Element{99923616639376095, 10339114964526300021, 6204619029868000785, 1288486622530663893, 14587509920085997152, 272081012460753233},
				A1: fp.Element{99923616639376095, 10339114964526300021, 6204619029868000785, 1288486622530663893, 14587509920085997152, 272081012460753233},
			},
			{
				A0: fp.Element{0},
				A1: fp.Element{6751177316358619845, 15498000274876530106, 6820146801716041242, 13487284328327464010, 776434812423573915, 1072939815054146550},
			},
			{
				A0: fp.Element{7399695662750302149, 14633322083064217648, 12051173786245255430, 9909266166264498601, 1288323043582377747, 379038003157372754},
				A1: fp.Element{6002735353327561446, 6023563502162542543, 13831244861028377885, 15776815867859765525, 4123780734888324547, 1494760614490167112},
			},
		},
		x)
}

func g2Isogeny(p *G2Affine) {

	den := make([]fptower.E2, 2)

	g2IsogenyYDenominator(&den[1], &p.X)
	g2IsogenyXDenominator(&den[0], &p.X)

	g2IsogenyYNumerator(&p.Y, &p.X, &p.Y)
	g2IsogenyXNumerator(&p.X, &p.X)

	den = fptower.BatchInvertE2(den)

	p.X.Mul(&p.X, &den[0])
	p.Y.Mul(&p.Y, &den[1])
}

// g2SqrtRatio computes the square root of u/v and returns 0 iff u/v was indeed a quadratic residue
// if not, we get sqrt(Z * u / v). Recall that Z is non-residue
// If v = 0, u/v is meaningless and the output is unspecified, without raising an error.
// The main idea is that since the computation of the square root involves taking large powers of u/v, the inversion of v can be avoided
func g2SqrtRatio(z *fptower.E2, u *fptower.E2, v *fptower.E2) uint64 {

	// https://www.ietf.org/archive/id/draft-irtf-cfrg-hash-to-curve-16.html#name-sqrt_ratio-for-any-field

	tv1 := fptower.E2{
		A0: fp.Element{8921533702591418330, 15859389534032789116, 3389114680249073393, 15116930867080254631, 3288288975085550621, 1021049300055853010},
		A1: fp.Element{8921533702591418330, 15859389534032789116, 3389114680249073393, 15116930867080254631, 3288288975085550621, 1021049300055853010},
	} //tv1 = c6

	var tv2, tv3, tv4, tv5 fptower.E2
	var exp big.Int
	// c4 = 7 = 2³ - 1
	// q is odd so c1 is at least 1.
	exp.SetBytes([]byte{7})

	tv2.Exp(*v, &exp) // 2. tv2 = vᶜ⁴
	tv3.Square(&tv2)  // 3. tv3 = tv2²
	tv3.Mul(&tv3, v)  // 4. tv3 = tv3 * v
	tv5.Mul(u, &tv3)  // 5. tv5 = u * tv3

	// c3 = 1001205140483106588246484290269935788605945006208159541241399033561623546780709821462541004956387089373434649096260670658193992783731681621012512651314777238193313314641988297376025498093520728838658813979860931248214124593092835
	exp.SetBytes([]byte{42, 67, 122, 75, 140, 53, 252, 116, 189, 39, 142, 170, 34, 242, 94, 158, 45, 201, 14, 80, 231, 4, 107, 70, 110, 89, 228, 147, 73, 232, 189, 5, 10, 98, 207, 209, 109, 220, 166, 239, 83, 20, 147, 48, 151, 142, 240, 17, 214, 134, 25, 200, 97, 133, 199, 178, 146, 232, 90, 135, 9, 26, 4, 150, 107, 249, 30, 211, 231, 27, 116, 49, 98, 195, 56, 54, 33, 19, 207, 215, 206, 214, 177, 215, 99, 130, 234, 178, 106, 160, 0, 1, 199, 24, 227})

	tv5.Exp(tv5, &exp)  // 6. tv5 = tv5ᶜ³
	tv5.Mul(&tv5, &tv2) // 7. tv5 = tv5 * tv2
	tv2.Mul(&tv5, v)    // 8. tv2 = tv5 * v
	tv3.Mul(&tv5, u)    // 9. tv3 = tv5 * u
	tv4.Mul(&tv3, &tv2) // 10. tv4 = tv3 * tv2

	// c5 = 4
	exp.SetBytes([]byte{4})
	tv5.Exp(tv4, &exp)      // 11. tv5 = tv4ᶜ⁵
	isQNr := g2NotOne(&tv5) // 12. isQR = tv5 == 1
	c7 := fptower.E2{
		A0: fp.Element{1921729236329761493, 9193968980645934504, 9862280504246317678, 6861748847800817560, 10375788487011937166, 4460107375738415},
		A1: fp.Element{16821121318233475459, 10183025025229892778, 1779012082459463630, 3442292649700377418, 1061500799026501234, 1352426537312017168},
	}
	tv2.Mul(&tv3, &c7)                 // 13. tv2 = tv3 * c7
	tv5.Mul(&tv4, &tv1)                // 14. tv5 = tv4 * tv1
	tv3.Select(int(isQNr), &tv3, &tv2) // 15. tv3 = CMOV(tv2, tv3, isQR)
	tv4.Select(int(isQNr), &tv4, &tv5) // 16. tv4 = CMOV(tv5, tv4, isQR)
	exp.Lsh(big.NewInt(1), 3-2)        // 18, 19: tv5 = 2ⁱ⁻² for i = c1

	for i := 3; i >= 2; i-- { // 17. for i in (c1, c1 - 1, ..., 2):

		tv5.Exp(tv4, &exp)               // 20.    tv5 = tv4ᵗᵛ⁵
		nE1 := g2NotOne(&tv5)            // 21.    e1 = tv5 == 1
		tv2.Mul(&tv3, &tv1)              // 22.    tv2 = tv3 * tv1
		tv1.Mul(&tv1, &tv1)              // 23.    tv1 = tv1 * tv1    Why not write square?
		tv5.Mul(&tv4, &tv1)              // 24.    tv5 = tv4 * tv1
		tv3.Select(int(nE1), &tv3, &tv2) // 25.    tv3 = CMOV(tv2, tv3, e1)
		tv4.Select(int(nE1), &tv4, &tv5) // 26.    tv4 = CMOV(tv5, tv4, e1)

		if i > 2 {
			exp.Rsh(&exp, 1) // 18, 19. tv5 = 2ⁱ⁻²
		}
	}

	*z = tv3
	return isQNr
}

func g2NotOne(x *fptower.E2) uint64 {

	//Assuming hash is implemented for G1 and that the curve is over Fp
	var one fp.Element
	return one.SetOne().NotEqual(&x.A0) | g1NotZero(&x.A1)

}

// g2MulByZ multiplies x by [-2, -1] and stores the result in z
func g2MulByZ(z *fptower.E2, x *fptower.E2) {

	z.Mul(x, &fptower.E2{
		A0: fp.Element{9794203289623549276, 7309342082925068282, 1139538881605221074, 15659550692327388916, 16008355200866287827, 582484205531694093},
		A1: fp.Element{4897101644811774638, 3654671041462534141, 569769440802610537, 17053147383018470266, 17227549637287919721, 291242102765847046},
	})

}

// https://www.ietf.org/archive/id/draft-irtf-cfrg-hash-to-curve-16.html#name-simplified-swu-method
// MapToCurve2 implements the SSWU map
// No cofactor clearing or isogeny
func MapToCurve2(u *fptower.E2) G2Affine {

	var sswuIsoCurveCoeffA = fptower.E2{
		A0: fp.Element{0},
		A1: fp.Element{16517514583386313282, 74322656156451461, 16683759486841714365, 815493829203396097, 204518332920448171, 1306242806803223655},
	}
	var sswuIsoCurveCoeffB = fptower.E2{
		A0: fp.Element{2515823342057463218, 7982686274772798116, 7934098172177393262, 8484566552980779962, 4455086327883106868, 1323173589274087377},
		A1: fp.Element{2515823342057463218, 7982686274772798116, 7934098172177393262, 8484566552980779962, 4455086327883106868, 1323173589274087377},
	}

	var tv1 fptower.E2
	tv1.Square(u) // 1.  tv1 = u²

	//mul tv1 by Z
	g2MulByZ(&tv1, &tv1) // 2.  tv1 = Z * tv1

	var tv2 fptower.E2
	tv2.Square(&tv1)    // 3.  tv2 = tv1²
	tv2.Add(&tv2, &tv1) // 4.  tv2 = tv2 + tv1

	var tv3 fptower.E2
	var tv4 fptower.E2
	tv4.SetOne()
	tv3.Add(&tv2, &tv4)                // 5.  tv3 = tv2 + 1
	tv3.Mul(&tv3, &sswuIsoCurveCoeffB) // 6.  tv3 = B * tv3

	tv2NZero := g2NotZero(&tv2)

	// tv4 = Z
	tv4 = fptower.E2{
		A0: fp.Element{9794203289623549276, 7309342082925068282, 1139538881605221074, 15659550692327388916, 16008355200866287827, 582484205531694093},
		A1: fp.Element{4897101644811774638, 3654671041462534141, 569769440802610537, 17053147383018470266, 17227549637287919721, 291242102765847046},
	}

	tv2.Neg(&tv2)
	tv4.Select(int(tv2NZero), &tv4, &tv2) // 7.  tv4 = CMOV(Z, -tv2, tv2 != 0)
	tv4.Mul(&tv4, &sswuIsoCurveCoeffA)    // 8.  tv4 = A * tv4

	tv2.Square(&tv3) // 9.  tv2 = tv3²

	var tv6 fptower.E2
	tv6.Square(&tv4) // 10. tv6 = tv4²

	var tv5 fptower.E2
	tv5.Mul(&tv6, &sswuIsoCurveCoeffA) // 11. tv5 = A * tv6

	tv2.Add(&tv2, &tv5) // 12. tv2 = tv2 + tv5
	tv2.Mul(&tv2, &tv3) // 13. tv2 = tv2 * tv3
	tv6.Mul(&tv6, &tv4) // 14. tv6 = tv6 * tv4

	tv5.Mul(&tv6, &sswuIsoCurveCoeffB) // 15. tv5 = B * tv6
	tv2.Add(&tv2, &tv5)                // 16. tv2 = tv2 + tv5

	var x fptower.E2
	x.Mul(&tv1, &tv3) // 17.   x = tv1 * tv3

	var y1 fptower.E2
	gx1NSquare := g2SqrtRatio(&y1, &tv2, &tv6) // 18. (is_gx1_square, y1) = sqrt_ratio(tv2, tv6)

	var y fptower.E2
	y.Mul(&tv1, u) // 19.   y = tv1 * u

	y.Mul(&y, &y1) // 20.   y = y * y1

	x.Select(int(gx1NSquare), &tv3, &x) // 21.   x = CMOV(x, tv3, is_gx1_square)
	y.Select(int(gx1NSquare), &y1, &y)  // 22.   y = CMOV(y, y1, is_gx1_square)

	y1.Neg(&y)
	y.Select(int(g2Sgn0(u)^g2Sgn0(&y)), &y, &y1)

	// 23.  e1 = sgn0(u) == sgn0(y)
	// 24.   y = CMOV(-y, y, e1)

	x.Div(&x, &tv4) // 25.   x = x / tv4

	return G2Affine{x, y}
}

func g2EvalPolynomial(z *fptower.E2, monic bool, coefficients []fptower.E2, x *fptower.E2) {
	dst := coefficients[len(coefficients)-1]

	if monic {
		dst.Add(&dst, x)
	}

	for i := len(coefficients) - 2; i >= 0; i-- {
		dst.Mul(&dst, x)
		dst.Add(&dst, &coefficients[i])
	}

	z.Set(&dst)
}

// g2Sgn0 is an algebraic substitute for the notion of sign in ordered fields
// Namely, every non-zero quadratic residue in a finite field of characteristic =/= 2 has exactly two square roots, one of each sign
// https://www.ietf.org/archive/id/draft-irtf-cfrg-hash-to-curve-16.html#name-the-sgn0-function
// The sign of an element is not obviously related to that of its Montgomery form
func g2Sgn0(z *fptower.E2) uint64 {

	nonMont := z.Bits()

	sign := uint64(0) // 1. sign = 0
	zero := uint64(1) // 2. zero = 1
	var signI uint64
	var zeroI uint64

	// 3. i = 1
	signI = nonMont.A0[0] % 2 // 4.   sign_i = x_i mod 2
	zeroI = g1NotZero(&nonMont.A0)
	zeroI = 1 ^ (zeroI|-zeroI)>>63 // 5.   zero_i = x_i == 0
	sign = sign | (zero & signI)   // 6.   sign = sign OR (zero AND sign_i) # Avoid short-circuit logic ops
	zero = zero & zeroI            // 7.   zero = zero AND zero_i
	// 3. i = 2
	signI = nonMont.A1[0] % 2 // 4.   sign_i = x_i mod 2
	// 5.   zero_i = x_i == 0
	sign = sign | (zero & signI) // 6.   sign = sign OR (zero AND sign_i) # Avoid short-circuit logic ops
	// 7.   zero = zero AND zero_i
	return sign

}

// MapToG2 invokes the SSWU map, and guarantees that the result is in g2
func MapToG2(u fptower.E2) G2Affine {
	res := MapToCurve2(&u)
	//this is in an isogenous curve
	g2Isogeny(&res)
	res.ClearCofactor(&res)
	return res
}

// EncodeToG2 hashes a message to a point on the G2 curve using the SSWU map.
// It is faster than HashToG2, but the result is not uniformly distributed. Unsuitable as a random oracle.
// dst stands for "domain separation tag", a string unique to the construction using the hash function
// https://www.ietf.org/archive/id/draft-irtf-cfrg-hash-to-curve-16.html#roadmap
func EncodeToG2(msg, dst []byte) (G2Affine, error) {

	var res G2Affine
	u, err := fp.Hash(msg, dst, 2)
	if err != nil {
		return res, err
	}

	res = MapToCurve2(&fptower.E2{
		A0: u[0],
		A1: u[1],
	})

	//this is in an isogenous curve
	g2Isogeny(&res)
	res.ClearCofactor(&res)
	return res, nil
}

// HashToG2 hashes a message to a point on the G2 curve using the SSWU map.
// Slower than EncodeToG2, but usable as a random oracle.
// dst stands for "domain separation tag", a string unique to the construction using the hash function
// https://www.ietf.org/archive/id/draft-irtf-cfrg-hash-to-curve-16.html#roadmap
func HashToG2(msg, dst []byte) (G2Affine, error) {
	u, err := fp.Hash(msg, dst, 2*2)
	if err != nil {
		return G2Affine{}, err
	}

	Q0 := MapToCurve2(&fptower.E2{
		A0: u[0],
		A1: u[1],
	})
	Q1 := MapToCurve2(&fptower.E2{
		A0: u[2+0],
		A1: u[2+1],
	})

	//TODO (perf): Add in E' first, then apply isogeny
	g2Isogeny(&Q0)
	g2Isogeny(&Q1)

	var _Q0, _Q1 G2Jac
	_Q0.FromAffine(&Q0)
	_Q1.FromAffine(&Q1).AddAssign(&_Q0)

	_Q1.ClearCofactor(&_Q1)

	Q1.FromJacobian(&_Q1)
	return Q1, nil
}

func g2NotZero(x *fptower.E2) uint64 {
	//Assuming G1 is over Fp and that if hashing is available for G2, it also is for G1
	return g1NotZero(&x.A0) | g1NotZero(&x.A1)

}
