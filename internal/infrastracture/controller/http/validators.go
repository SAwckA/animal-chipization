package http

import (
	"animal-chipization/internal/domain"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ExcludeWhitespace валидатор на наличие пробелов в поле
// доступен как binding:"exclude_whitespace"
var ExcludeWhitespace validator.Func = func(fl validator.FieldLevel) bool {
	return !(strings.Contains(fl.Field().String(), " ") || strings.Contains(fl.Field().String(), "\n") || strings.Contains(fl.Field().String(), "\t"))
}

// AllowedStrings валидатор, разрешающий только перечисленные через ";" строки
// доступен как binding:"allowed_strings=test,strings"
var AllowedStrings validator.Func = func(fl validator.FieldLevel) bool {
	for _, v := range strings.Split(fl.Param(), ";") {
		if v == fl.Field().String() {
			return true
		}
	}
	return false
}

// ParamID Нужен, чтобы доставать из параметров ( /{Id} ) и валидировать
func ParamID(cCp *gin.Context, name string) (int, error) {
	paramString := cCp.Param(name)

	if paramString == "" || paramString == "null" {
		return 0, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "Invalid id param",
		}
	}

	res, err := strconv.Atoi(paramString)

	if err != nil {
		return 0, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "Invalid id param",
		}
	}

	if res <= 0 {
		return 0, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "Invalid id param",
		}
	}

	return res, nil
}
