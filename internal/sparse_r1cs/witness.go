package r1cs

import (
	"fmt"

	cs "github.com/consensys/gnark/constraint/bn254"
	"noirgnark/internal/acir"
)

// BuildWitnessMap constructs a gnark SparseR1CS with all variables registered,
// and returns the r1cs + a map from Noir witness index to gnark variable ID.
// Noir witness index -1 maps to the gnark constant "1" variable.
func BuildWitnessMap(fn *acir.Function) (*cs.SparseR1CS, map[int]int) {
	r1cs := cs.NewSparseR1CS(len(fn.Opcodes))
	witnessMap := make(map[int]int, fn.CurrentWitnessIndex+2)

	one := r1cs.AddPublicVariable("1")
	witnessMap[-1] = one

	publicSet := make(map[int]struct{}, len(fn.PublicParameters)+len(fn.ReturnValues))
	for _, idx := range fn.PublicParameters {
		publicSet[idx] = struct{}{}
	}
	for _, idx := range fn.ReturnValues {
		publicSet[idx] = struct{}{}
	}

	for _, idx := range fn.PublicParameters {
		if _, exists := witnessMap[idx]; exists {
			continue
		}
		witnessMap[idx] = r1cs.AddPublicVariable(fmt.Sprintf("w%d", idx))
	}

	for _, idx := range fn.ReturnValues {
		if _, exists := witnessMap[idx]; exists {
			continue
		}
		witnessMap[idx] = r1cs.AddPublicVariable(fmt.Sprintf("w%d", idx))
	}

	for _, idx := range fn.PrivateParameters {
		if _, exists := witnessMap[idx]; exists {
			continue
		}
		witnessMap[idx] = r1cs.AddSecretVariable(fmt.Sprintf("w%d", idx))
	}

	for i := 0; i <= fn.CurrentWitnessIndex; i++ {
		if _, isPublic := publicSet[i]; isPublic {
			continue
		}
		if _, exists := witnessMap[i]; exists {
			continue
		}
		witnessMap[i] = r1cs.AddInternalVariable()
	}

	return r1cs, witnessMap
}
