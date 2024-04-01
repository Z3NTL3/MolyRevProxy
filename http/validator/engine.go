/*
Written by Efdal Sancak (aka z3ntl3)

github.com/z3ntl3

Disclaimer: Educational purposes only
License: GNU
*/
package validator

import (
	"errors"

	vld "github.com/go-playground/validator/v10"
)

type Validator struct {
	*vld.Validate
}

func (v *Validator) ValidateStruct(obj any) error {
	var error_ error
	if err := v.Struct(obj); err != nil {
		for _, err := range err.(vld.ValidationErrors) {
			error_ = errors.New(err.Error())
			break
		}
	}
	return error_
}

func (v *Validator) Engine() any {
	return v.Validate
}
