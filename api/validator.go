package api

import (
	"interview/db/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	currency, ok := fieldLevel.Field().Interface().(string)
	if !ok {
		return false
	}

	return util.IsSupportedCurrency(currency)
}
