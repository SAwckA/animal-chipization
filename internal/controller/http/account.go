package http

import (
	"animal-chipization/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type accountUsecase interface {
	GetAccountByID(id int32) (*domain.Account, error)
	SearchAccount(firstName, lastName, email string, from, size int) []*domain.AccountResponse
	FullUpdateAccount(currentAccount *domain.Account, newAccount *domain.Account) (*domain.AccountResponse, error)
	DeleteAccount(accoundID int) error
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

func (h *AccountHandler) getAccountByID(ctx *gin.Context) {
	accountIDParam := ctx.Param("accountId")

	accountID, err := strconv.ParseInt(accountIDParam, 10, 32)

	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "Invalid id", nil)
		return
	}

	account, err := h.usecase.GetAccountByID(int32(accountID))

	if err != nil {
		newErrorResponse(ctx, http.StatusNotFound, "Account not found", nil)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":        account.ID,
		"firstName": account.FirstName,
		"lastName":  account.LastName,
		"email":     account.Email,
	})
}

func (h *AccountHandler) searchAccount(ctx *gin.Context) {

	var err error

	firstName := ctx.Query("firstName")
	lastName := ctx.Query("lastName")
	email := ctx.Query("email")

	var from int = 0

	fromString := ctx.Query("from")
	if fromString != "" {
		from, err = strconv.Atoi(fromString)

		if err != nil || from < 0 {
			newErrorResponse(ctx, http.StatusBadRequest, "invalid from param", err)
			return
		}
	} else {
		from = 0
	}

	var size int = 10

	sizeString := ctx.Query("size")

	if sizeString != "" {
		size, err = strconv.Atoi(sizeString)
		if err != nil || size <= 0 {
			newErrorResponse(ctx, http.StatusBadRequest, "invalid size param", err)
			return
		}
	} else {
		size = 10
	}

	logrus.Warnln(firstName, lastName, email, from, size)
	result := h.usecase.SearchAccount(firstName, lastName, email, from, size)

	ctx.JSON(http.StatusOK, result)
}

func (h *AccountHandler) updateAccount(c *gin.Context) {
	accountIDString := c.Param("accountId")

	var newAccount *domain.Account

	if err := c.BindJSON(&newAccount); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request data", nil)
		return
	}

	accountID, err := strconv.Atoi(accountIDString)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid accountId", nil)
		return
	}

	currentAccount, _ := c.Get(accountCtx)

	if currentAccount.(*domain.Account).ID != int32(accountID) {
		newErrorResponse(c, http.StatusForbidden, "You are not owner or account not found", nil)
		return
	}

	if newAccount.Email == "" || newAccount.Email == "null" || !validateEmail(newAccount.Email) || contentsOnlySpaces(newAccount.Email) {
		newErrorResponse(c, http.StatusBadRequest, "Invalid email", nil)
		return
	}
	if newAccount.FirstName == "" || newAccount.FirstName == "null" || contentsOnlySpaces(newAccount.FirstName) {
		newErrorResponse(c, http.StatusBadRequest, "Invalid firstName", nil)
		return
	}
	if newAccount.LastName == "" || newAccount.LastName == "null" || contentsOnlySpaces(newAccount.LastName) {
		newErrorResponse(c, http.StatusBadRequest, "Invalid lastName", nil)
		return
	}
	if newAccount.Password == "" || newAccount.Password == "null" || contentsOnlySpaces(newAccount.Password) {
		newErrorResponse(c, http.StatusBadRequest, "Invalid password", nil)
		return
	}

	result, err := h.usecase.FullUpdateAccount(currentAccount.(*domain.Account), newAccount)

	if err != nil {
		newErrorResponse(c, http.StatusConflict, "This email already used", nil)
		return
	}

	c.JSON(http.StatusOK, result)
}

// contentsOnlySpaces проверяет состоит ли строка только из пробелов
func contentsOnlySpaces(s string) bool {
	for _, c := range s {
		if c != ' ' {
			return false
		}
	}
	return true
}

func (h *AccountHandler) deleteAccount(c *gin.Context) {
	accountIDString := c.Param("accountId")
	if accountIDString == "" || accountIDString == "null" {
		newErrorResponse(c, http.StatusBadRequest, "Invalid accountId", nil)
		return
	}

	accountID, err := strconv.Atoi(accountIDString)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid accountId", nil)
		return
	}

	if accountID <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid accountId", nil)
		return
	}

	account, _ := c.Get(accountCtx)

	if accountID != int(account.(*domain.Account).ID) {
		newErrorResponse(c, http.StatusForbidden, "You are not owner or account not found", nil)
		return
	}

	if err = h.usecase.DeleteAccount(accountID); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Account linked with animal", nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}
