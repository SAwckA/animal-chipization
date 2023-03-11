package domain

const (
	AccountSearchDefaultFrom = 0
	AccountSearchDefaultSize = 10
)

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (a *Account) MapResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":        a.ID,
		"firstName": a.FirstName,
		"lastName":  a.LastName,
		"email":     a.Email,
	}
}

func NewAccount(params RegistrationParams) *Account {
	return &Account{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  params.Password,
	}
}

type RegistrationParams struct {
	FirstName string `json:"firstName" binding:"required,exclude_whitespace"`
	LastName  string `json:"lastName" binding:"required,exclude_whitespace"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,exclude_whitespace"`
}

type SearchAccount struct {
	FirstName *string `form:"firstName"`
	LastName  *string `form:"lastName"`
	Email     *string `form:"email"`
	Size      *int    `form:"size"`
	From      *int    `form:"from"`
}

func (s *SearchAccount) Validate() error {
	err := &ApplicationError{
		OriginalError: nil,
		SimplifiedErr: ErrInvalidInput,
		Description:   "validation error",
	}
	var defaultFrom, defaultSize = AccountSearchDefaultFrom, AccountSearchDefaultSize

	if s.From == nil {
		s.From = &defaultFrom
	}
	if s.Size == nil {
		s.Size = &defaultSize
	}

	switch {
	case *s.From < 0:
		return err
	case *s.Size <= 0:
		return err

	default:
		return nil
	}
}

type UpdateAccount struct {
	ID        int
	FirstName string `json:"firstName" binding:"required,exclude_whitespace"`
	LastName  string `json:"lastName" binding:"required,exclude_whitespace"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,exclude_whitespace"`
}
