package main

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

type input struct {
	Email    string `validate:"email"`
	ChildAge int    `validate:"min=0,max=17"`
}

func main() {
	validate := validator.New()

	v := input{
		Email:    "test@example.com",
		ChildAge: 9,
	}
	errs := validate.Struct(v)
	if errs != nil {
		if ve, ok := errs.(validator.ValidationErrors); ok {
			for _, fe := range ve {
				fmt.Printf("error: %v\n", fe)
				fmt.Printf("  field: %v\n", fe.Field())
				fmt.Printf("  tag: %v\n", fe.Tag())
			}
		}
	}

	iv := input{
		Email:    "test",
		ChildAge: -1,
	}
	errs = validate.Struct(iv)
	if errs != nil {
		if ve, ok := errs.(validator.ValidationErrors); ok {
			for _, fe := range ve {
				fmt.Printf("error: %v\n", fe)
				fmt.Printf("  field: %v\n", fe.Field())
				fmt.Printf("  tag: %v\n", fe.Tag())
			}
		}
	}
}
