package main

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/test"
)

const SIZE = 1 << PROF

func Recursiva(merkle [SIZE * 2][]byte, nodo int, otros [][]byte, dir []big.Int, arr [SIZE]*big.Int, t *testing.T) error {
	if nodo*2 >= len(merkle) {
		var Asig Circuit
		for i := 0; i < len(otros); i++ {
			Asig.Path[i] = otros[i]
		}
		for i := 0; i < len(dir); i++ {
			Asig.PathDirecc[i] = dir[i]
		}
		Asig.Root = merkle[1]
		Asig.Leafq = getBytes(arr[nodo-len(merkle)/2])
		var circuit Circuit
		assert := test.NewAssert(t)
		err := test.IsSolved(&circuit, &Asig, ecc.BN254.ScalarField())
		assert.NoError(err)
		return nil
	}
	Recursiva(merkle, nodo*2, append(otros, merkle[nodo*2+1]), append(dir, *big.NewInt(0)), arr, t)
	Recursiva(merkle, nodo*2+1, append(otros, merkle[nodo*2]), append(dir, *big.NewInt(1)), arr, t)
	return nil
}

func getBytes(b *big.Int) []byte {
	const SIZE = 32
	bElement := fr.NewElement(b.Uint64())
	res := make([]byte, SIZE)
	for i := 0; i < SIZE; i++ {
		res[i] = bElement.Bytes()[i]
	}
	return res
}

func nodeSumByte(izq []byte, der []byte) []byte {
	mimc := hash.MIMC_BN254.New()
	mimc.Write(izq)
	mimc.Write(der)
	return mimc.Sum(make([]byte, 0))
}

func Test(t *testing.T) {
	for nt := 0; nt < 1000; nt++ {
		var Merkle [SIZE * 2][]byte
		var arr [SIZE]*big.Int
		for i := 0; i < SIZE; i++ {
			arr[i] = big.NewInt(rand.Int63n(int64(1e18)))
			Merkle[i+SIZE] = getBytes(arr[i])
		}
		for i := SIZE - 1; i > 0; i-- {
			Merkle[i] = nodeSumByte(Merkle[i*2], Merkle[i*2+1])
		}

		Recursiva(Merkle, 1, [][]byte{}, []big.Int{}, arr, t)
	}
}
