// Copyright 2020 Consensys Software Inc.
// Licensed under the Apache License, Version 2.0. See the LICENSE file for details.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package fft

import (
	"fmt"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bw6-633/fr"
)

type bitReverseVariant struct {
	name string
	buf  []fr.Element
	fn   func([]fr.Element)
}

const maxSizeBitReverse = 1 << 23

var bitReverse = []bitReverseVariant{
	{name: "bitReverseNaive", buf: make([]fr.Element, maxSizeBitReverse), fn: bitReverseNaive},
	{name: "BitReverse", buf: make([]fr.Element, maxSizeBitReverse), fn: BitReverse},
	{name: "bitReverseCobraInPlace", buf: make([]fr.Element, maxSizeBitReverse), fn: bitReverseCobraInPlace},
}

func TestBitReverse(t *testing.T) {

	// generate a random []fr.Element array of size 2**20
	pol := make([]fr.Element, maxSizeBitReverse)
	one := fr.One()
	pol[0].SetRandom()
	for i := 1; i < maxSizeBitReverse; i++ {
		pol[i].Add(&pol[i-1], &one)
	}

	// for each size, check that all the bitReverse functions fn compute the same result.
	for size := 2; size <= maxSizeBitReverse; size <<= 1 {

		// copy pol into the buffers
		for _, data := range bitReverse {
			copy(data.buf, pol[:size])
		}

		// compute bit reverse shuffling
		for _, data := range bitReverse {
			data.fn(data.buf[:size])
		}

		// all bitReverse.buf should hold the same result
		for i := 0; i < size; i++ {
			for j := 1; j < len(bitReverse); j++ {
				if !bitReverse[0].buf[i].Equal(&bitReverse[j].buf[i]) {
					t.Fatalf("bitReverse %s and %s do not compute the same result", bitReverse[0].name, bitReverse[j].name)
				}
			}
		}

		// bitReverse back should be identity
		for _, data := range bitReverse {
			data.fn(data.buf[:size])
		}

		for i := 0; i < size; i++ {
			for j := 1; j < len(bitReverse); j++ {
				if !bitReverse[0].buf[i].Equal(&bitReverse[j].buf[i]) {
					t.Fatalf("(fn-1) bitReverse %s and %s do not compute the same result", bitReverse[0].name, bitReverse[j].name)
				}
			}
		}
	}

}

func BenchmarkBitReverse(b *testing.B) {
	// generate a random []fr.Element array of size 2**22
	pol := make([]fr.Element, maxSizeBitReverse)
	one := fr.One()
	pol[0].SetRandom()
	for i := 1; i < maxSizeBitReverse; i++ {
		pol[i].Add(&pol[i-1], &one)
	}

	// copy pol into the buffers
	for _, data := range bitReverse {
		copy(data.buf, pol[:maxSizeBitReverse])
	}

	// benchmark for each size, each bitReverse function
	for size := 1 << 18; size <= maxSizeBitReverse; size <<= 1 {
		for _, data := range bitReverse {
			b.Run(fmt.Sprintf("name=%s/size=%d", data.name, size), func(b *testing.B) {
				b.ResetTimer()
				for j := 0; j < b.N; j++ {
					data.fn(data.buf[:size])
				}
			})
		}
	}
}
