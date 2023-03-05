package http

import (
	"animal-chipization/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

const accountIDParam = "accountId"

type accountUsecase interface {
	Get(id int) (*domain.Account, error)
	Search(dto *domain.SearchAccount) ([]domain.Account, error)
	Update(old *domain.Account, newAccount domain.UpdateAccount) (*domain.Account, error)
	Delete(executor *domain.Account, id int) error
}

type AccountHandler struct {
	usecase    accountUsecase
	middleware *AuthMiddleware
}

type AccountResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func NewAccountHandler(usecase accountUsecase, middleware *AuthMiddleware) *AccountHandler {
	return &AccountHandler{usecase: usecase, middleware: middleware}
}

func (h *AccountHandler) InitRoutes(router *gin.Engine) *gin.Engine {
	account := router.Group("/accounts")
	{
		account.Use(h.middleware.checkAuthHeaderMiddleware)
		account.GET("/:accountId",
			errorHandlerWrap(h.getAccountByID),
		)

		account.GET("/search",
			errorHandlerWrap(h.searchAccount),
		)
		account.PUT("/:accountId",
			h.middleware.authMiddleware,
			errorHandlerWrap(h.updateAccount),
		)
		account.DELETE("/:accountId",
			h.middleware.authMiddleware,
			errorHandlerWrap(h.deleteAccount),
		)

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
func (h *AccountHandler) getAccountByID(c *gin.Context) error {

	accountID, err := validateID(c.Copy(), accountIDParam)
	if err != nil {
		return err
	}

	account, err := h.usecase.Get(accountID)

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, account.Response())
	return nil
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
func (h *AccountHandler) searchAccount(c *gin.Context) error {

	var input domain.SearchAccount
	if err := c.BindQuery(&input); err != nil {
		return err
	}

	result, err := h.usecase.Search(&input)
	if err != nil {
		return err
	}

	resp := make([]map[string]interface{}, 0)

	if result != nil {
		for _, v := range result {
			resp = append(resp, v.Response())
		}
	}

	c.JSON(http.StatusOK, resp)
	return nil

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
func (h *AccountHandler) updateAccount(c *gin.Context) error {

	accountID, err := validateID(c.Copy(), accountIDParam)
	if err != nil {
		return err
	}

	var input domain.UpdateAccount

	if err := c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	input.ID = accountID

	currentAccount := c.MustGet(accountCtx)

	result, err := h.usecase.Update(currentAccount.(*domain.Account), input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, result.Response())
	return nil
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
func (h *AccountHandler) deleteAccount(c *gin.Context) error {

	accountID, err := validateID(c.Copy(), accountIDParam)
	if err != nil {
		return err
	}

	account := c.MustGet(accountCtx)

	err = h.usecase.Delete(account.(*domain.Account), accountID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, nil)
	return nil
}
