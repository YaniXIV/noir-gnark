package r1cs

import (
	"fmt"
	cs "github.com/consensys/gnark/constraint/bn254"
	"noirgnark/internal/acir"
)

func BuildR1CS(c string) {
	p := acir.DecodeAcir(c)
	// im only handling main for now.
	fn := p.Functions[0]

	r1cs := cs.NewR1CS(0)
	witnessMap := make(map[int]int)

	ONE := r1cs.AddPublicVariable("1")
	witnessMap[-1] = ONE

	for _, idx := range fn.PublicParameters {
		varID := r1cs.AddPublicVariable(fmt.Sprintf("w%d", idx))
		witnessMap[idx] = varID
	}

	for _, idx := range fn.PrivateParameters {
		varID := r1cs.AddSecretVariable(fmt.Sprintf("w%d", idx))
		witnessMap[idx] = varID
	}

	for i := 0; i <= fn.CurrentWitnessIndex; i++ {
		if _, exists := witnessMap[i]; !exists {
			varID := r1cs.AddInternalVariable()
			witnessMap[i] = varID
		}
	}

	fmt.Println("witnessMap:", witnessMap)
}
