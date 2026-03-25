package r1cs

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark/constraint"
	cs "github.com/consensys/gnark/constraint/bn254"
	"noirgnark/internal/acir"
)

// fieldToBigInt converts a [32]byte big-endian field element to *big.Int.
func fieldToBigInt(fe acir.FieldElement) *big.Int {
	return new(big.Int).SetBytes(fe[:])
}

func ensureGenericSparseBlueprint(r1cs *cs.SparseR1CS) constraint.BlueprintID {
	for i, blueprint := range r1cs.Blueprints {
		if _, ok := blueprint.(*constraint.BlueprintGenericSparseR1C[constraint.U64]); ok {
			return constraint.BlueprintID(i)
		}
	}
	return r1cs.AddBlueprint(&constraint.BlueprintGenericSparseR1C[constraint.U64]{})
}

func coeffIDFromFieldElement(r1cs *cs.SparseR1CS, fe acir.FieldElement) uint32 {
	return r1cs.AddCoeff(r1cs.FromInterface(fieldToBigInt(fe)))
}

func constantCoeffID(r1cs *cs.SparseR1CS, fe acir.FieldElement) uint32 {
	term := r1cs.MakeTerm(r1cs.FromInterface(fieldToBigInt(fe)), 0)
	term.MarkConstant()
	return uint32(term.CoeffID())
}

// AddAssertZeroConstraint translates one AssertZero opcode into a gnark SparseR1C
// and adds it to the r1cs. Uses witnessMap to resolve Noir witness indices to gnark var IDs.
func AddAssertZeroConstraint(r1cs *cs.SparseR1CS, op *acir.AssertZero, witnessMap map[int]int) {
	blueprint := ensureGenericSparseBlueprint(r1cs)

	var (
		xa = uint32(witnessMap[-1])
		xb = uint32(witnessMap[-1])
		xc = uint32(witnessMap[-1])
		qL = constantCoeffID(r1cs, acir.FieldElement{})
		qR = constantCoeffID(r1cs, acir.FieldElement{})
		qO = constantCoeffID(r1cs, acir.FieldElement{})
		qM = constantCoeffID(r1cs, acir.FieldElement{})
		qC = constantCoeffID(r1cs, op.QC)
	)

	if len(op.MulTerms) > 0 {
		mt := op.MulTerms[0]
		xa = uint32(witnessMap[mt.LHS])
		xb = uint32(witnessMap[mt.RHS])
		qM = coeffIDFromFieldElement(r1cs, mt.Coeff)
	}

	if len(op.MulTerms) == 0 {
		if len(op.LinearCombinations) > 0 {
			xa = uint32(witnessMap[op.LinearCombinations[0].Witness])
			qL = coeffIDFromFieldElement(r1cs, op.LinearCombinations[0].Coeff)
		}
		if len(op.LinearCombinations) > 1 {
			xb = uint32(witnessMap[op.LinearCombinations[1].Witness])
			qR = coeffIDFromFieldElement(r1cs, op.LinearCombinations[1].Coeff)
		}
		if len(op.LinearCombinations) > 2 {
			xc = uint32(witnessMap[op.LinearCombinations[2].Witness])
			qO = coeffIDFromFieldElement(r1cs, op.LinearCombinations[2].Coeff)
		}
	} else {
		for _, lc := range op.LinearCombinations {
			switch {
			case qL == constraint.CoeffIdZero && lc.Witness == int(xa):
				qL = coeffIDFromFieldElement(r1cs, lc.Coeff)
			case qR == constraint.CoeffIdZero && lc.Witness == int(xb):
				qR = coeffIDFromFieldElement(r1cs, lc.Coeff)
			case qO == constraint.CoeffIdZero:
				xc = uint32(witnessMap[lc.Witness])
				qO = coeffIDFromFieldElement(r1cs, lc.Coeff)
			case qL == constraint.CoeffIdZero:
				qL = coeffIDFromFieldElement(r1cs, lc.Coeff)
			case qR == constraint.CoeffIdZero:
				qR = coeffIDFromFieldElement(r1cs, lc.Coeff)
			}
		}

		// Noir commonly emits the qL/qR witnesses as original witness indices.
		// If they matched the mul wires by witness index, the checks above may not have
		// triggered because xa/xb are gnark IDs. Re-run against the mul witnesses directly.
		if len(op.MulTerms) > 0 {
			mt := op.MulTerms[0]
			for _, lc := range op.LinearCombinations {
				switch {
				case qL == constraint.CoeffIdZero && lc.Witness == mt.LHS:
					qL = coeffIDFromFieldElement(r1cs, lc.Coeff)
				case qR == constraint.CoeffIdZero && lc.Witness == mt.RHS:
					qR = coeffIDFromFieldElement(r1cs, lc.Coeff)
				case qO == constraint.CoeffIdZero:
					xc = uint32(witnessMap[lc.Witness])
					qO = coeffIDFromFieldElement(r1cs, lc.Coeff)
				}
			}
		}
	}

	r1cs.AddSparseR1C(constraint.SparseR1C{
		XA: xa,
		XB: xb,
		XC: xc,
		QL: qL,
		QR: qR,
		QO: qO,
		QM: qM,
		QC: qC,
	}, blueprint)
}

// BuildSparseR1CS is the top-level entry point.
// Takes raw ACIR JSON, decodes it, builds witness map, translates all AssertZero opcodes.
// Returns the complete SparseR1CS ready for gnark proving.
func BuildSparseR1CS(acirJSON []byte) (*cs.SparseR1CS, map[int]int, error) {
	program, err := acir.DecodeProgram(acirJSON)
	if err != nil {
		return nil, nil, err
	}

	var receiver acir.Function
	mainFn, err := receiver.MainFunction(program)
	if err != nil {
		return nil, nil, err
	}

	r1cs, witnessMap := BuildWitnessMap(mainFn)
	ensureGenericSparseBlueprint(r1cs)

	for _, op := range mainFn.Opcodes {
		if op.AssertZero == nil {
			continue
		}
		AddAssertZeroConstraint(r1cs, op.AssertZero, witnessMap)
	}

	if r1cs.GetNbConstraints() == 0 {
		return nil, nil, fmt.Errorf("no AssertZero constraints found in main function")
	}

	return r1cs, witnessMap, nil
}
