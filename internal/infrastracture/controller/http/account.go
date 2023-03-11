package http

import (
	"animal-chipization/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

const accountIDParam = "accountId"

type accountUsecase interface {
	Get(id int) (*domain.Account, error)
	Search(params *domain.SearchAccount) ([]domain.Account, error)
	Update(old *domain.Account, newAccount *domain.UpdateAccount) (*domain.Account, error)
	Delete(executor *domain.Account, id int) error
}

type AccountHandler struct {
	usecase    accountUsecase
	middleware authMiddleware
}

func NewAccountHandler(usecase accountUsecase, middleware authMiddleware) *AccountHandler {
	return &AccountHandler{usecase: usecase, middleware: middleware}
}

func (h *AccountHandler) InitRoutes(router *gin.Engine) *gin.Engine {
	account := router.Group("/accounts")
	{
		account.Use(h.middleware.checkAuthHeaderMiddleware)
		account.GET("/:accountId",
			errorHandlerWrap(h.accountByID),
		)

		account.GET("/search",
			errorHandlerWrap(h.search),
		)
		account.PUT("/:accountId",
			h.middleware.authMiddleware,
			errorHandlerWrap(h.update),
		)
		account.DELETE("/:accountId",
			h.middleware.authMiddleware,
			errorHandlerWrap(h.delete),
		)

	}

	return router
}

func (h *AccountHandler) accountByID(c *gin.Context) error {

	accountID, err := validateID(c.Copy(), accountIDParam)
	if err != nil {
		return err
	}

	account, err := h.usecase.Get(accountID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, account.MapResponse())
	return nil
}

func (h *AccountHandler) search(c *gin.Context) error {

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
			resp = append(resp, v.MapResponse())
		}
	}

	c.JSON(http.StatusOK, resp)
	return nil

}

func (h *AccountHandler) update(c *gin.Context) error {

	accountID, err := validateID(c.Copy(), accountIDParam)
	if err != nil {
		return err
	}

	var input *domain.UpdateAccount
	if err = c.BindJSON(&input); err != nil {
		return NewErrBind(err)
	}

	input.ID = accountID
	currentAccount := c.MustGet(accountCtx)

	result, err := h.usecase.Update(currentAccount.(*domain.Account), input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, result.MapResponse())
	return nil
}

func (h *AccountHandler) delete(c *gin.Context) error {

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
