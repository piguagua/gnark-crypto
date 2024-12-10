// Copyright 2020 Consensys Software Inc.
// Licensed under the Apache License, Version 2.0. See the LICENSE file for details.

// Code generated by consensys/gnark-crypto DO NOT EDIT

// Package fp contains field arithmetic operations for modulus = 0x126633...70000d.
//
// The API is similar to math/big (big.Int), but the operations are significantly faster (up to 20x for the modular multiplication on amd64, see also https://hackmd.io/@gnark/modular_multiplication)
//
// The modulus is hardcoded in all the operations.
//
// Field elements are represented as an array, and assumed to be in Montgomery form in all methods:
//
//	type Element [10]uint64
//
// # Usage
//
// Example API signature:
//
//	// Mul z = x * y (mod q)
//	func (z *Element) Mul(x, y *Element) *Element
//
// and can be used like so:
//
//	var a, b Element
//	a.SetUint64(2)
//	b.SetString("984896738")
//	a.Mul(a, b)
//	a.Sub(a, a)
//	 .Add(a, b)
//	 .Inv(a)
//	b.Exp(b, new(big.Int).SetUint64(42))
//
// Modulus q =
//
//	q[base10] = 20494478644167774678813387386538961497669590920908778075528754551012016751717791778743535050360001387419576570244406805463255765034468441182772056330021723098661967429339971741066259394985997
//	q[base16] = 0x126633cc0f35f63fc1a174f01d72ab5a8fcd8c75d79d2c74e59769ad9bbda2f8152a6c0fadea490b8da9f5e83f57c497e0e8850edbda407d7b5ce7ab839c2253d369bd31147f73cd74916ea4570000d
//
// # Warning
//
// This code has not been audited and is provided as-is. In particular, there is no security guarantees such as constant time implementation or side-channel attack resistance.
package fp
