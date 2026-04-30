package main

import (
	"core-ticket/base/helpers/struct_validator"
	"core-ticket/base/utility/nullable"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gopkg.in/guregu/null.v4"
)

func GetValidator() (validate *validator.Validate) {
	validate = binding.Validator.Engine().(*validator.Validate)
	setCustomValidator(validate)
	return
}

func setCustomValidator(validate *validator.Validate) {
	validate.RegisterCustomTypeFunc(
		struct_validator.NullFlake,
		null.String{},
		null.Int{},
		null.Bool{},
		null.Float{},
		null.Time{},
		nullable.Int{},
		nullable.Float{},
		nullable.Bool{},
		nullable.Time{},
		nullable.String{},
	)
}
