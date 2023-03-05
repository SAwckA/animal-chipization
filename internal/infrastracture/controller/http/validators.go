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

// getIntQuery достаёт значение из Query ?{key} из копии gin.Context,
// При этом подставляет deafult_ значение, если задано null или значение отсутствует
func getIntQuery(cCp *gin.Context, key string, default_ int) (int, error) {
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
