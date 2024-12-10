// Copyright 2020 Consensys Software Inc.
// Licensed under the Apache License, Version 2.0. See the LICENSE file for details.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package eddsa

import (
	"crypto/sha256"
	"math/big"
	"math/rand"
	"testing"

	crand "crypto/rand"

	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark-crypto/hash"
)

func Example() {
	// instantiate hash function
	hFunc := hash.MIMC_BN254.New()

	// create a eddsa key pair
	privateKey, _ := GenerateKey(crand.Reader)
	publicKey := privateKey.PublicKey

	// generate a message (the size must be a multiple of the size of Fr)
	var _msg fr.Element
	_msg.SetRandom()
	msg := _msg.Marshal()

	// sign the message
	signature, _ := privateKey.Sign(msg, hFunc)

	// verifies signature
	isValid, _ := publicKey.Verify(signature, msg, hFunc)
	if !isValid {
		fmt.Println("1. invalid signature")
	} else {
		fmt.Println("1. valid signature")
	}

	// Output: 1. valid signature
}

func TestNonMalleability(t *testing.T) {

	// buffer too big
	t.Run("buffer_overflow", func(t *testing.T) {
		bsig := make([]byte, 2*sizeFr+1)
		var sig Signature
		_, err := sig.SetBytes(bsig)
		if err != errWrongSize {
			t.Fatal("should raise wrong size error")
		}
	})

	// R overflows p_mod
	t.Run("R_overflow", func(t *testing.T) {
		bsig := make([]byte, 2*sizeFr)
		frMod := fr.Modulus()
		r := big.NewInt(1)
		r.Add(frMod, r)
		buf := r.Bytes()
		for i := 0; i < sizeFr; i++ {
			bsig[sizeFr-1-i] = buf[i]
		}

		var sig Signature
		_, err := sig.SetBytes(bsig)
		if err != errRBiggerThanPMod {
			t.Fatal("should raise error r >= p_mod")
		}
	})

	// S overflows r_mod
	t.Run("S_overflow", func(t *testing.T) {
		bsig := make([]byte, 2*sizeFr)
		o := big.NewInt(1)
		cp := twistededwards.GetEdwardsCurve()
		o.Add(&cp.Order, o)
		buf := o.Bytes()
		copy(bsig[sizeFr:], buf[:])
		big.NewInt(1).FillBytes(bsig[:sizeFr])

		var sig Signature
		_, err := sig.SetBytes(bsig)
		if err != errSBiggerThanRMod {
			t.Fatal("should raise error s >= r_mod")
		}
	})

}

func TestNoZeros(t *testing.T) {
	t.Run("R.Y=0", func(t *testing.T) {
		// R points are 0
		var sig Signature
		sig.R.X.SetInt64(1)
		sig.R.Y.SetInt64(0)
		s := big.NewInt(1)
		s.FillBytes(sig.S[:])
		bts := sig.Bytes()
		var newSig Signature
		_, err := newSig.SetBytes(bts)
		if err != errZero {
			t.Fatal("expected error for zero R.Y")
		}
	})
	t.Run("S=0", func(t *testing.T) {
		// S is 0
		var R twistededwards.PointAffine
		cp := twistededwards.GetEdwardsCurve()
		R.ScalarMultiplication(&cp.Base, big.NewInt(1))
		var sig Signature
		sig.R.Set(&R)
		bts := sig.Bytes()
		var newSig Signature
		_, err := newSig.SetBytes(bts)
		if err != errZero {
			t.Fatal("expected error for zero S")
		}
	})
}

func TestSerialization(t *testing.T) {

	src := rand.NewSource(0)
	r := rand.New(src) //#nosec G404 weak rng is fine here

	privKey1, err := GenerateKey(r)
	if err != nil {
		t.Fatal(err)
	}
	pubKey1 := privKey1.PublicKey

	privKey2, err := GenerateKey(r)
	if err != nil {
		t.Fatal(err)
	}
	pubKey2 := privKey2.PublicKey

	pubKeyBin1 := pubKey1.Bytes()
	pubKey2.SetBytes(pubKeyBin1)
	pubKeyBin2 := pubKey2.Bytes()
	if len(pubKeyBin1) != len(pubKeyBin2) {
		t.Fatal("Inconsistent size")
	}
	for i := 0; i < len(pubKeyBin1); i++ {
		if pubKeyBin1[i] != pubKeyBin2[i] {
			t.Fatal("Error serialize(deserialize(.))")
		}
	}

	privKeyBin1 := privKey1.Bytes()
	privKey2.SetBytes(privKeyBin1)
	privKeyBin2 := privKey2.Bytes()
	if len(privKeyBin1) != len(privKeyBin2) {
		t.Fatal("Inconsistent size")
	}
	for i := 0; i < len(privKeyBin1); i++ {
		if privKeyBin1[i] != privKeyBin2[i] {
			t.Fatal("Error serialize(deserialize(.))")
		}
	}
}

func TestEddsaMIMC(t *testing.T) {

	src := rand.NewSource(0)
	r := rand.New(src) //#nosec G404 weak rng is fine here

	// create eddsa obj and sign a message
	privKey, err := GenerateKey(r)
	if err != nil {
		t.Fatal(nil)
	}
	pubKey := privKey.PublicKey
	hFunc := hash.MIMC_BN254.New()

	var frMsg fr.Element
	frMsg.SetString("44717650746155748460101257525078853138837311576962212923649547644148297035978")
	msgBin := frMsg.Bytes()
	signature, err := privKey.Sign(msgBin[:], hFunc)
	if err != nil {
		t.Fatal(err)
	}

	// verifies correct msg
	res, err := pubKey.Verify(signature, msgBin[:], hFunc)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatal("Verify correct signature should return true")
	}

	// verifies wrong msg
	frMsg.SetString("44717650746155748460101257525078853138837311576962212923649547644148297035979")
	msgBin = frMsg.Bytes()
	res, err = pubKey.Verify(signature, msgBin[:], hFunc)
	if err != nil {
		t.Fatal(err)
	}
	if res {
		t.Fatal("Verify wrong signature should be false")
	}

}

func TestEddsaSHA256(t *testing.T) {

	src := rand.NewSource(0)
	r := rand.New(src) //#nosec G404 weak rng is fine here

	hFunc := sha256.New()

	// create eddsa obj and sign a message
	// create eddsa obj and sign a message

	privKey, err := GenerateKey(r)
	pubKey := privKey.PublicKey
	if err != nil {
		t.Fatal(err)
	}

	signature, err := privKey.Sign([]byte("message"), hFunc)
	if err != nil {
		t.Fatal(err)
	}

	// verifies correct msg
	res, err := pubKey.Verify(signature, []byte("message"), hFunc)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatal("Verify correct signature should return true")
	}

	// verifies wrong msg
	res, err = pubKey.Verify(signature, []byte("wrong_message"), hFunc)
	if err != nil {
		t.Fatal(err)
	}
	if res {
		t.Fatal("Verify wrong signature should be false")
	}

}

// benchmarks

func BenchmarkVerify(b *testing.B) {

	src := rand.NewSource(0)
	r := rand.New(src) //#nosec G404 weak rng is fine here

	hFunc := hash.MIMC_BN254.New()

	// create eddsa obj and sign a message
	privKey, err := GenerateKey(r)
	pubKey := privKey.PublicKey
	if err != nil {
		b.Fatal(err)
	}
	var frMsg fr.Element
	frMsg.SetString("44717650746155748460101257525078853138837311576962212923649547644148297035978")
	msgBin := frMsg.Bytes()
	signature, _ := privKey.Sign(msgBin[:], hFunc)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pubKey.Verify(signature, msgBin[:], hFunc)
	}
}
