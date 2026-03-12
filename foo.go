package noirgnark

import (
	"fmt"

	"github.com/YaniXIV/noir-go"

	"noirgnark/internal/acir"
)

func TestDecode() *acir.Program {
	fmt.Println("inital commit")

	comp, err := noirgo.Compile("noirtest")
	if err != nil {
		panic(err)
	}

	fmt.Println(comp.ACIR.JSON)
	p := acir.DecodeAcir(comp.ACIR.JSON)
	return p

}
