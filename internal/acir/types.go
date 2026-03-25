package acir

import (
	"encoding/json"
	"fmt"
)

type Program struct {
	Functions              []Function `json:"functions"`
	UnconstrainedFunctions []Function `json:"unconstrained_functions"`
}

type Function struct {
	FunctionName        string      `json:"function_name"`
	CurrentWitnessIndex int         `json:"current_witness_index"`
	Opcodes             []RawOpcode `json:"opcodes"`
	PrivateParameters   []int       `json:"private_parameters"`
	PublicParameters    []int       `json:"public_parameters"`
	ReturnValues        []int       `json:"return_values"`
}

type RawOpcode = Opcode

type Opcode struct {
	AssertZero *AssertZero `json:"AssertZero,omitempty"`
}

type AssertZero struct {
	MulTerms           []MulTerm           `json:"mul_terms"`
	LinearCombinations []LinearCombination `json:"linear_combinations"`
	QC                 FieldElement        `json:"q_c"`
}

type FieldElement [32]byte

type MulTerm struct {
	Coeff FieldElement
	LHS   int
	RHS   int
}

type LinearCombination struct {
	Coeff   FieldElement
	Witness int
}

func (m *MulTerm) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 3 {
		return fmt.Errorf("mul term: expected 3 elements, got %d", len(raw))
	}

	if err := json.Unmarshal(raw[0], &m.Coeff); err != nil {
		return fmt.Errorf("mul term coeff: %w", err)
	}
	if err := json.Unmarshal(raw[1], &m.LHS); err != nil {
		return fmt.Errorf("mul term lhs: %w", err)
	}
	if err := json.Unmarshal(raw[2], &m.RHS); err != nil {
		return fmt.Errorf("mul term rhs: %w", err)
	}

	return nil
}

func (l *LinearCombination) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 2 {
		return fmt.Errorf("linear combination: expected 2 elements, got %d", len(raw))
	}

	if err := json.Unmarshal(raw[0], &l.Coeff); err != nil {
		return fmt.Errorf("linear combination coeff: %w", err)
	}
	if err := json.Unmarshal(raw[1], &l.Witness); err != nil {
		return fmt.Errorf("linear combination witness: %w", err)
	}

	return nil
}

func (op *Opcode) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if payload, ok := raw["AssertZero"]; ok {
		var assertZero AssertZero
		if err := json.Unmarshal(payload, &assertZero); err != nil {
			return fmt.Errorf("AssertZero: %w", err)
		}
		op.AssertZero = &assertZero
	}

	return nil
}
