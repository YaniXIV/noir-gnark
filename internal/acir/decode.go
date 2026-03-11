package acir

import (
	"encoding/json"
)

func DecodeAcir(data string) *Program {
	var program Program
	err := json.Unmarshal([]byte(data), &program)
	if err != nil {
		panic(err)
	}
	return &program
}
