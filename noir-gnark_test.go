package noirgnark_test

import (
	"fmt"
	"github.com/YaniXIV/noir-go"
	"noirgnark/internal/acir"
	"noirgnark/internal/r1cs"
	"testing"
)

func TestFoobar(t *testing.T) {

	comp, err := noirgo.Compile("noirtest")
	if err != nil {
		panic(err)
	}
	r1cs.BuildR1CS(comp.ACIR.JSON)

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
