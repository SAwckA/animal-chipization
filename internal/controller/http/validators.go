package http

import (
	"errors"
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

// GetIntQuery достаёт значение из Query ?{key} из копии gin.Context,
// При этом подставляет deafult_ значение, если задано null или значение отсутствует
func GetIntQuery(cCp *gin.Context, key string, default_ int) (int, error) {
	val := cCp.Query(key)

	if val == "" || val == "null" {
		return default_, nil
	}

	res, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return res, err
}

func validateID(cCp *gin.Context, name string) (int, error) {
	paramString := cCp.Param(name)

	if paramString == "" || paramString == "null" {
		return 0, errors.New("invalid id")
	}

	res, err := strconv.Atoi(paramString)

	if err != nil {
		return 0, errors.New("invalid id")
	}

	if res <= 0 {
		return 0, errors.New("invalid id")
	}

	return res, nil
}
