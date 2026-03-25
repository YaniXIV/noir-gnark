package noirgnark_test

import (
	"fmt"
	"github.com/YaniXIV/noir-go"
	"noirgnark/internal/acir"
	//"noirgnark/internal/sparse_r1cs"
	"testing"
)

const circuit_constraints_20k = "testdata/circuit_constraints_20k "
const circuit_int_heavy_5k = "testdata/circuit_int_heavy_5k"
const circuit_small = "testdata/circuit_small"
const circuit_wide_4k = "testdata/circuit_wide_4k"

//really should pull this stuff out of the test, but some good experimentation.

var validStatuses = map[string]struct{}{
	// only supporting assertzero for now.
	"AssertZero": {},
}

func TestFoobar(t *testing.T) {

	comp, err := noirgo.Compile(circuit_small)
	if err != nil {
		panic(err)
	}
	c := comp.ACIR.JSON
	//fmt.Println(comp.ACIR.JSON)
	// r1cs.BuildInitialWitness(comp.ACIR.JSON)
	p := acir.DecodeAcir(c)
	mainFn, err := getEntry(p)
	if err != nil {
		panic(err)
	}
	getValidOpcodes(mainFn.Opcodes)

}

func getValidOpcodes(opcodes []acir.Opcode) {
	for _, opcode := range opcodes {
		if opcode.AssertZero != nil {

			az := opcode.AssertZero
			fmt.Printf("AssertZero\n")
			fmt.Printf("      QC: %x\n", az.QC)
			fmt.Printf("      MulTerms (%d):\n", len(az.MulTerms))
			for j, mt := range az.MulTerms {
				fmt.Printf("        [%d] coeff=%x  lhs=%d  rhs=%d\n", j, mt.Coeff, mt.LHS, mt.RHS)
			}
			fmt.Printf("      LinearCombinations (%d):\n", len(az.LinearCombinations))
			for j, lc := range az.LinearCombinations {
				fmt.Printf("        [%d] coeff=%x  witness=%d\n", j, lc.Coeff, lc.Witness)
			}
		} else {
			fmt.Printf("(unhandled/nil opcode)\n")

		}

	}

}
func getEntry(p *acir.Program) (*acir.Function, error) {
	for _, fn := range p.Functions {
		if fn.FunctionName == "main" {
			fmt.Println("caught main")
			fmt.Println(fn)
			return &fn, nil
		}
	}
	return nil, fmt.Errorf("main function doesnt exist")

}

func PrintProgram(p *acir.Program) {
	for _, fn := range p.Functions {
		fmt.Printf("=== Function: %s ===\n", fn.FunctionName)
		fmt.Printf("  CurrentWitnessIndex: %d\n", fn.CurrentWitnessIndex)
		fmt.Printf("  PrivateParameters:   %v\n", fn.PrivateParameters)
		fmt.Printf("  PublicParameters:    %v\n", fn.PublicParameters)
		fmt.Printf("  ReturnValues:        %v\n", fn.ReturnValues)
		fmt.Printf("  Opcodes (%d):\n", len(fn.Opcodes))

		for i, op := range fn.Opcodes {
			fmt.Printf("    [%d] ", i)
			if op.AssertZero != nil {
				az := op.AssertZero
				fmt.Printf("AssertZero\n")
				fmt.Printf("      QC: %x\n", az.QC)
				fmt.Printf("      MulTerms (%d):\n", len(az.MulTerms))
				for j, mt := range az.MulTerms {
					fmt.Printf("        [%d] coeff=%x  lhs=%d  rhs=%d\n", j, mt.Coeff, mt.LHS, mt.RHS)
				}
				fmt.Printf("      LinearCombinations (%d):\n", len(az.LinearCombinations))
				for j, lc := range az.LinearCombinations {
					fmt.Printf("        [%d] coeff=%x  witness=%d\n", j, lc.Coeff, lc.Witness)
				}
			} else {
				fmt.Printf("(unhandled/nil opcode)\n")
			}
		}
	}
}
