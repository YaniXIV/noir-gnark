package noirgnark

import (
	"fmt"

	"github.com/YaniXIV/noir-go"

	"noirgnark/internal/acir"
)

func TestDecode() {
	fmt.Println("inital commit")

	comp, err := noirgo.Compile("noirtest")
	if err != nil {
		panic(err)
	}

	fmt.Println(comp.ACIR.JSON)
	acir.DecodeAcir(comp.ACIR.JSON)

}
