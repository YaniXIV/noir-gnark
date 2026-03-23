package r1cs

import (
	"fmt"
	"math/big"

	fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/constraint"
	cs "github.com/consensys/gnark/constraint/bn254"
	"noirgnark/internal/acir"
)

// TODO: Make this a method for acir.ACIR.
// qL⋅xa + qR⋅xb + qO⋅xc + qM⋅(xa⋅xb) + qC == 0
// (qL)i · xai + (qR)i · xbi + (qO)i · xci + (qM)i · (xai xbi ) + (qC)i = 0
func BuildSparseR1CS(circuit acir.ACIR, values fr_bn254.Vector) (*cs_bn254.SparseR1CS, fr_bn254.Vector, fr_bn254.Vector) {
	sparseR1CS := cs_bn254.NewSparseR1CS(int(circuit.CurrentWitness) - 1)

	publicVariables, secretVariables, indexMap := backend.HandleValues(circuit, sparseR1CS, values)
	handleOpcodes(circuit, sparseR1CS, indexMap)

	return sparseR1CS, publicVariables, secretVariables
}

// TODO:
// (qL)i · xai + (qR)i · xbi + (qO)i · xci + (qM)i · (xai xbi ) + (qC)i = 0
func handleAssertZeroOpcode(a *acir.ACIR) {

}

func handleArithmeticOpcode(a *acir_opcode.ArithmeticOpcode, sparseR1CS constraint.SparseR1CS, indexMap map[string]int) {
	var xa, xb, xc int
	var qL, qR, qO, qC, qM1, qM2 constraint.Coeff

	// Case qM⋅(xa⋅xb)
	if len(a.MulTerms) != 0 {
		mulTerm := a.MulTerms[0]
		qM1 = sparseR1CS.FromInterface(mulTerm.Coefficient)
		qM2 = sparseR1CS.FromInterface(1)
		xa = indexMap[fmt.Sprint(int(mulTerm.MultiplicandIndex))]
		xb = indexMap[fmt.Sprint(int(mulTerm.MultiplierIndex))]
	}

	// Case qO⋅xc
	if len(a.SimpleTerms) == 1 {
		qOwOTerm := a.SimpleTerms[0]
		qO = sparseR1CS.FromInterface(qOwOTerm.Coefficient)
		xc = indexMap[fmt.Sprint(int(qOwOTerm.VariableIndex))]
	}

	// Case qL⋅xa + qR⋅xb
	if len(a.SimpleTerms) == 2 {
		// qL⋅xa
		qLwLTerm := a.SimpleTerms[0]
		qL = sparseR1CS.FromInterface(qLwLTerm.Coefficient)
		xa = indexMap[fmt.Sprint(int(qLwLTerm.VariableIndex))]
		// qR⋅xb
		qRwRTerm := a.SimpleTerms[1]
		qR = sparseR1CS.FromInterface(qRwRTerm.Coefficient)
		xb = indexMap[fmt.Sprint(int(qRwRTerm.VariableIndex))]
	}

	// Case qL⋅xa + qR⋅xb + qO⋅xc
	if len(a.SimpleTerms) == 3 {
		// qL⋅xa
		qLwLTerm := a.SimpleTerms[0]
		qL = sparseR1CS.FromInterface(qLwLTerm.Coefficient)
		xa = indexMap[fmt.Sprint(int(qLwLTerm.VariableIndex))]
		// qR⋅xb
		qRwRTerm := a.SimpleTerms[1]
		qR = sparseR1CS.FromInterface(qRwRTerm.Coefficient)
		xb = indexMap[fmt.Sprint(int(qRwRTerm.VariableIndex))]
		// qO⋅xc
		qOwOTerm := a.SimpleTerms[2]
		qO = sparseR1CS.FromInterface(qOwOTerm.Coefficient)
		xc = indexMap[fmt.Sprint(int(qOwOTerm.VariableIndex))]
	}

	// Add the qC term
	qC = sparseR1CS.FromInterface(a.QC)

	K := sparseR1CS.MakeTerm(&qC, 0)
	K.MarkConstant()

	constraint := constraint.SparseR1C{
		L: sparseR1CS.MakeTerm(&qL, xa),
		R: sparseR1CS.MakeTerm(&qR, xb),
		O: sparseR1CS.MakeTerm(&qO, xc),
		M: [2]constraint.Term{sparseR1CS.MakeTerm(&qM1, xa), sparseR1CS.MakeTerm(&qM2, xb)},
		K: K.CoeffID(),
	}

	sparseR1CS.AddConstraint(constraint)
}
