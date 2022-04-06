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

package fptower

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

// ------------------------------------------------------------
// tests

func TestE6ReceiverIsOperand(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := GenE6()
	genB := GenE6()
	genE2 := GenE2()

	properties.Property("[BLS12-377] Having the receiver as operand (addition) should output the same result", prop.ForAll(
		func(a, b *E6) bool {
			var c, d E6
			d.Set(a)
			c.Add(a, b)
			a.Add(a, b)
			b.Add(&d, b)
			return a.Equal(b) && a.Equal(&c) && b.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("[BLS12-377] Having the receiver as operand (sub) should output the same result", prop.ForAll(
		func(a, b *E6) bool {
			var c, d E6
			d.Set(a)
			c.Sub(a, b)
			a.Sub(a, b)
			b.Sub(&d, b)
			return a.Equal(b) && a.Equal(&c) && b.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("[BLS12-377] Having the receiver as operand (mul) should output the same result", prop.ForAll(
		func(a, b *E6) bool {
			var c, d E6
			d.Set(a)
			c.Mul(a, b)
			a.Mul(a, b)
			b.Mul(&d, b)
			return a.Equal(b) && a.Equal(&c) && b.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("[BLS12-377] Having the receiver as operand (square) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Square(a)
			a.Square(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS12-377] Having the receiver as operand (neg) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Neg(a)
			a.Neg(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS12-377] Having the receiver as operand (double) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Double(a)
			a.Double(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS12-377] Having the receiver as operand (mul by non residue) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.MulByNonResidue(a)
			a.MulByNonResidue(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS12-377] Having the receiver as operand (Inverse) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Inverse(a)
			a.Inverse(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS12-377] Having the receiver as operand (mul by E2) should output the same result", prop.ForAll(
		func(a *E6, b *E2) bool {
			var c E6
			c.MulByE2(a, b)
			a.MulByE2(a, b)
			return a.Equal(&c)
		},
		genA,
		genE2,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestE6Ops(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := GenE6()
	genB := GenE6()
	genE2 := GenE2()

	properties.Property("[BLS12-377] sub & add should leave an element invariant", prop.ForAll(
		func(a, b *E6) bool {
			var c E6
			c.Set(a)
			c.Add(&c, b).Sub(&c, b)
			return c.Equal(a)
		},
		genA,
		genB,
	))

	properties.Property("[BLS12-377] mul & inverse should leave an element invariant", prop.ForAll(
		func(a, b *E6) bool {
			var c, d E6
			d.Inverse(b)
			c.Set(a)
			c.Mul(&c, b).Mul(&c, &d)
			return c.Equal(a)
		},
		genA,
		genB,
	))

	properties.Property("[BLS12-377] inverse twice should leave an element invariant", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Inverse(a).Inverse(&b)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS12-377] neg twice should leave an element invariant", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Neg(a).Neg(&b)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS12-377] square and mul should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b, c E6
			b.Mul(a, a)
			c.Square(a)
			return b.Equal(&c)
		},
		genA,
	))

	properties.Property("[BLS12-377] Double and add twice should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Add(a, a)
			a.Double(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("[BLS12-377] Mul by non residue should be the same as multiplying by (0,1,0)", prop.ForAll(
		func(a *E6) bool {
			var b, c E6
			b.B1.A0.SetOne()
			c.Mul(a, &b)
			a.MulByNonResidue(a)
			return a.Equal(&c)
		},
		genA,
	))

	properties.Property("[BLS12-377] MulByE2 MulByE2 inverse should leave an element invariant", prop.ForAll(
		func(a *E6, b *E2) bool {
			var c E6
			var d E2
			d.Inverse(b)
			c.MulByE2(a, b).MulByE2(&c, &d)
			return c.Equal(a)
		},
		genA,
		genE2,
	))

	properties.Property("[BLS12-377] Mul and MulBy01 should output the same result", prop.ForAll(
		func(a *E6, c0, c1 *E2) bool {
			var b E6
			b.B0.Set(c0)
			b.B1.Set(c1)
			b.Mul(&b, a)
			a.MulBy01(c0, c1)
			return b.Equal(a)
		},
		genA,
		genE2,
		genE2,
	))

	properties.Property("[BLS12-377] Mul and MulBy1 should output the same result", prop.ForAll(
		func(a *E6, c1 *E2) bool {
			var b E6
			b.B1.Set(c1)
			b.Mul(&b, a)
			a.MulBy1(c1)
			return b.Equal(a)
		},
		genA,
		genE2,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

// ------------------------------------------------------------
// benches

func BenchmarkE6Add(b *testing.B) {
	var a, c E6
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Add(&a, &c)
	}
}

func BenchmarkE6Sub(b *testing.B) {
	var a, c E6
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Sub(&a, &c)
	}
}

func BenchmarkE6Mul(b *testing.B) {
	var a, c E6
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Mul(&a, &c)
	}
}

func BenchmarkE6Square(b *testing.B) {
	var a E6
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Square(&a)
	}
}

func BenchmarkE6Inverse(b *testing.B) {
	var a E6
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Inverse(&a)
	}
}

func TestE6Div(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)

	genA := GenE6()
	genB := GenE6()

	properties.Property("[BLS12-377] dividing then multiplying by the same element does nothing", prop.ForAll(
		func(a, b *E6) bool {
			var c E6
			c.Div(a, b)
			c.Mul(&c, b)
			return c.Equal(a)
		},
		genA,
		genB,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}
