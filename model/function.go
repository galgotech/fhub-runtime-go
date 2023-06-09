// Copyright 2023 The fhub-runtime-go Authors
// This file is part of fhub-runtime-go.
//
// This file is part of fhub-runtime-go.
// fhub-runtime-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// fhub-runtime-go is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with fhub-runtime-go. If not, see <https://www.gnu.org/licenses/>.

package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

type Function struct {
	inputValue  cue.Value `fhub:"input" fhub-unmarshal:"true"`
	outputValue cue.Value `fhub:"output" fhub-unmarshal:"true"`

	Inputs  []string `validate:"min=1"`
	Outputs []string `validate:"min=1"`
}

func (f *Function) Unmarshal(field string, value cue.Value) (err error) {
	if field == "input" {
		f.inputValue = value
		f.Inputs, err = type_struct(f.inputValue)
		if err != nil {
			return err
		}

	} else if field == "output" {
		f.outputValue = value
		f.Outputs, err = type_struct(f.outputValue)
		if err != nil {
			return err
		}

	} else {
		return errors.New("invalid field")
	}

	return nil
}

func (f Function) UnmarshalInput(data []byte) (map[string]any, error) {
	return unmarshalToMap(data, f.inputValue)
}

func (f Function) UnmarshalOutput(data []byte) (map[string]any, error) {
	return unmarshalToMap(data, f.outputValue)
}

func (f Function) ValidateInput(data []byte) error {
	return validate(data, f.inputValue)
}

func (f Function) ValidateOutput(data []byte) error {
	return validate(data, f.outputValue)
}

type Parameter struct {
	Label   string
	cueType cue.Kind
}

func (f *Parameter) Type(k reflect.Kind) bool {
	switch k {
	case reflect.String:
		return f.cueType.IsAnyOf(cue.StringKind)
	case reflect.Bool:
		return f.cueType.IsAnyOf(cue.StringKind)
	}

	return false
}

func type_struct(value cue.Value) ([]string, error) {
	fields, err := value.Fields()
	if err != nil {
		return nil, err
	}

	parameters := []string{}
	for fields.Next() {
		parameters = append(parameters, fields.Label())
	}

	return parameters, nil
}

func unmarshalToMap(data []byte, value cue.Value) (map[string]any, error) {
	valueUnify, err := unifyDataAndValue(data, value)
	if err != nil {
		return nil, fmt.Errorf("error in cue value unify: %s", err)
	}

	dataMap := map[string]any{}
	fields, err := valueUnify.Fields()
	if err != nil {
		return nil, err
	}

	for fields.Next() {
		var value any
		fields.Value().Decode(&value)
		dataMap[fields.Label()] = value
	}

	return dataMap, nil
}

func validate(data []byte, value cue.Value) error {
	valueUnify, err := unifyDataAndValue(data, value)
	if err != nil {
		return err
	}

	return valueUnify.Validate(
		cue.Attributes(true),
		cue.Optional(true),
		cue.Hidden(true),
		cue.Concrete(true),
	)
}

func unifyDataAndValue(data []byte, value cue.Value) (cue.Value, error) {
	valid := json.Valid(data)
	if !valid {
		return value, nil
	}

	dataValue := cuecontext.New().CompileBytes(data)
	valueUnify := value.UnifyAccept(dataValue, value)
	if valueUnify.Err() != nil {
		return value, valueUnify.Err()
	}

	return valueUnify, nil
}
