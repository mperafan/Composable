package main

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash"
	"github.com/consensys/gnark/std/hash/mimc"
)

const CANTL = 8 // cantidad de hojas
const PROF = 4  //profundidad del arbol

type Circuit struct { //circuito de mi programa que representa un merkle tree
	Root       frontend.Variable       `gnark:",public"` // hojas de mi merkle tree
	Path       [PROF]frontend.Variable `gnark:",public"` // Ruta de la hoja por la que quiero preguntar
	PathDirecc [PROF]frontend.Variable `gnark:",public"` //Si la hoja corresponde al nodo izquierdo, es un 0, si corresponde al nodo derecho, es un 1
	Leafq      frontend.Variable       `gnark:",public"` // Hoja por la que quiero saber si pertenece o no al arbol
}

func nodeSum(api frontend.API, h hash.Hash, izq, der frontend.Variable) frontend.Variable {
	h.Reset()
	h.Write(izq, der)
	newn := h.Sum() //el nodo cuyo valor es el hash de sus hijos
	return newn
}

func (circuit *Circuit) Define(api frontend.API) error {
	var hash2 frontend.Variable
	H, _ := mimc.NewMiMC(api)
	hash2 = circuit.Leafq
	for i := PROF - 1; i >= 0; i-- {
		hash2 = api.Select(circuit.PathDirecc[i], nodeSum(api, &H, circuit.Path[i], hash2), nodeSum(api, &H, hash2, circuit.Path[i]))
	}

	api.AssertIsEqual(circuit.Root, hash2)
	return nil
}

func main() {

}
