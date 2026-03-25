package acir

import (
	"encoding/json"
	"fmt"
)

func DecodeProgram(data []byte) (*Program, error) {
	var program Program
	if err := json.Unmarshal(data, &program); err != nil {
		return nil, err
	}
	return &program, nil
}

func DecodeAcir(data string) *Program {
	program, err := DecodeProgram([]byte(data))
	if err != nil {
		panic(err)
	}
	return program
}

func (f *Function) AssertZeroOpcodes() []*AssertZero {
	var out []*AssertZero
	for _, op := range f.Opcodes {
		if op.AssertZero != nil {
			out = append(out, op.AssertZero)
		}
	}
	return out
}

func (f *Function) MainFunction(p *Program) (*Function, error) {
	for i := range p.Functions {
		if p.Functions[i].FunctionName == "main" {
			return &p.Functions[i], nil
		}
	}
	return nil, fmt.Errorf("main function not found")
}
