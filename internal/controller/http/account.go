package http

import (
	"animal-chipization/internal/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type accountUsecase interface {
	Get(id int) (*domain.Account, error)
	Search(dto domain.SearchAccountDTO, size int, from int) (*[]domain.Account, error)
	Update(old *domain.Account, newAccount domain.UpdateAccountDTO) (*domain.Account, error)
	Delete(executor *domain.Account, id int) error
}

type AccountHandler struct {
	usecase    accountUsecase
	middleware *Middleware
}

type AccountResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func NewAccountHandler(usecase accountUsecase, middleware *Middleware) *AccountHandler {
	return &AccountHandler{usecase: usecase, middleware: middleware}
}

func (h *AccountHandler) InitRoutes(router *gin.Engine) *gin.Engine {
	account := router.Group("/accounts")
	{
		account.GET("/:accountId",
			h.middleware.ckeckAuthHeaderMiddleware,
			h.getAccountByID,
		)

		account.GET("/search",
			h.middleware.ckeckAuthHeaderMiddleware,
			h.searchAccount,
		)
		account.PUT("/:accountId",
			h.middleware.authMiddleware,
			h.updateAccount)
		account.DELETE("/:accountId",
			h.middleware.authMiddleware,
			h.deleteAccount)

	}

	return router
}

// API 1: Получение информации об аккаунте пользователя
// GET - /accounts/{accountId}
// 	{accountId}: "int"	// Идентификатор аккаунта пользователя
// 	- request
// 		Body {
// 			empty
// 		}
//
// 	- response
// 		Body {
//			"id": "int",		// Идентификатор аккаунта пользователя
//			"firstName": "string",	// Имя пользователя
// 			"lastName": "string",	// Фамилия пользователя
// 			"email": "string"		// Адрес электронной почты
// 		}
func (h *AccountHandler) getAccountByID(c *gin.Context) {

	var input domain.AccountID

	if err := c.BindUri(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	account, err := h.usecase.Get(input.ID)

	switch errors.Unwrap(err) {
	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case nil:
		c.JSON(http.StatusOK, account.Response())

	default:
		unreachableError(c, err)
	}

}

// API 2: Поиск аккаунтов пользователей по параметрам
// 	GET - /accounts/search
// 		?firstName={firstName}
// 		&lastName={lastName}
// 		&email={email}
// 		&from={from}
// 		&size={size}
//
// 		{firstName}: "string",	// Имя пользователя, может использоваться только часть имени без учета регистра, если null, не участвует в фильтрации
// 		{firstName}: "string",	// Фамилия пользователя, может использоваться только часть фамилии без учета регистра, если null, не участвует в фильтрации
// 		{email}: "string",		// Адрес электронной почты, может использоваться только часть адреса электронной почты без учета регистра, если null, не участвует в фильтрации
// 		{from}: "int"		// Количество элементов, которое необходимо пропустить для формирования страницы с результатами (по умолчанию 0)
// 		{size}: "int"		// Количество элементов на странице (по умолчанию 10)
// 	- request
// 	Body  {
// 			empty
// 		}
//
// 	- response
// 	Body [
// 		{
// 			“id”: "int",		// Идентификатор аккаунта пользователя
// 			"firstName": "string",	// Имя пользователя
// 			"lastName": "string",	// Фамилия пользователя
// 			"email": "string"	// Адрес электронной почты
// 		}
// 	]
func (h *AccountHandler) searchAccount(c *gin.Context) {

	var input domain.SearchAccountDTO

	if err := c.ShouldBindQuery(&input); err != nil {
		badRequest(c, err.Error())
		return
	}

	size, err := GetIntQuery(c.Copy(), "size", domain.AccountDefaultSize)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	from, err := GetIntQuery(c.Copy(), "from", domain.AccountDefaultFrom)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	result, err := h.usecase.Search(input, size, from)

	switch errors.Unwrap(err) {
	case domain.ErrInvalidInput:
		badRequest(c, err.Error())

	case nil:
		resp := make([]map[string]interface{}, 0)

		if result != nil {
			for _, v := range *result {
				resp = append(resp, v.Response())
			}
		}

		c.JSON(http.StatusOK, resp)
	default:
		unreachableError(c, err)
	}
}

// API 3: Обновление данных аккаунта пользователя
// 	PUT - /accounts/{accountId}
// 		{accountId}: "int"	// Идентификатор аккаунта пользователя
//
// - request
// 	Body {
// 		"firstName": "string",	// Новое имя пользователя
// 		"lastName": "string",	// Новая фамилия пользователя
// 		"email": "string",		// Новый адрес электронной почты
// 		"password": "string"    	// Пароль от аккаунта
// 	}
//
// - response
// 	Body {
// 		"id": "int",		// Идентификатор аккаунта пользователя
// 		"firstName": "string",	// Новое имя пользователя
// 		"lastName": "string",	// Новая фамилия пользователя
// 		"email": "string"			// Новый адрес электронной почты
// 	}
func (h *AccountHandler) updateAccount(c *gin.Context) {

	var accountID domain.AccountID
	if err := c.BindUri(&accountID); err != nil {
		badRequest(c, err.Error())
		return
	}

	var input domain.UpdateAccountDTO

	if err := c.BindJSON(&input); err != nil {
		badRequest(c, err.Error())
		return
	}

	input.ID = accountID.ID

	currentAccount, _ := c.Get(accountCtx)

	result, err := h.usecase.Update(currentAccount.(*domain.Account), input)

	switch errors.Unwrap(err) {
	case domain.ErrBadDatabaseOut:
		conflictResponse(c, err.Error())

	case domain.ErrNotFound:
		notFoundResponse(c, err.Error())

	case domain.ErrForbidden:
		forbiddenResponse(c, err.Error())

	case nil:
		c.JSON(http.StatusOK, result.Response())

	default:
		unreachableError(c, err)
	}
}

// API 4: Удаление аккаунта пользователя
// 	DELETE - /accounts/{accountId}
// 		{accountId}: "int"	// Идентификатор аккаунта пользователя
//
// - request
// 	Body {
// 		empty
// 	}
//
// - response
// 	Body {
// 		empty
// 	}
func (h *AccountHandler) deleteAccount(c *gin.Context) {

	var accountID domain.AccountID
	if err := c.BindUri(&accountID); err != nil {
		badRequest(c, err.Error())
		return
	}

	account, _ := c.Get(accountCtx)

	err := h.usecase.Delete(account.(*domain.Account), accountID.ID)

	switch errors.Unwrap(err) {
	case domain.ErrBadDatabaseOut:
		badRequest(c, err.Error())

	case domain.ErrForbidden:
		forbiddenResponse(c, err.Error())

	case domain.ErrNotFound:
		forbiddenResponse(c, err.Error())

	case nil:
		c.JSON(http.StatusOK, nil)

	default:
		unreachableError(c, err)
	}
}
