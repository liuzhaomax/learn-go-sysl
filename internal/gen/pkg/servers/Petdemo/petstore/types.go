// Code generated by sysl DO NOT EDIT.
package petstore

import (
	"fmt"

	"github.com/anz-bank/sysl-go/validator"
)

// Error_ ...
type Error_ struct {
	Code    int64  `json:"code" url:"code"`
	Message string `json:"message" url:"message"`
}

// Error fulfills the error interface.
func (s Error_) Error() string {
	type plain Error_
	return fmt.Sprintf("%+v", plain(s))
}

// GetPetListRequest ...
type GetPetListRequest struct {
}

// *Error_ validator
func (s *Error_) Validate() error {
	return validator.Validate(s)
}

// Pet ...
type Pet = string
