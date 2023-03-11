package usecase

import (
	"animal-chipization/internal/domain"
)

type accountRepository interface {
	GetByID(id int) (*domain.Account, error)
	Search(params *domain.SearchAccount) ([]domain.Account, error)
	Update(newAccount *domain.Account) error
	Delete(accountID int) error

	Create(account *domain.Account) (int, error)
	GetByEmail(email string) (*domain.Account, error)
}

type AccountUsecase struct {
	repo accountRepository
}

func NewAccountUsecase(repo accountRepository) *AccountUsecase {
	return &AccountUsecase{repo: repo}
}

func (u *AccountUsecase) Get(id int) (*domain.Account, error) {
	return u.repo.GetByID(id)
}

func (u *AccountUsecase) Search(params *domain.SearchAccount) ([]domain.Account, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return u.repo.Search(params)
}

func (u *AccountUsecase) Update(old *domain.Account, newAccount *domain.UpdateAccount) (*domain.Account, error) {
	if old.ID != newAccount.ID {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrForbidden,
			Description:   "update not your account",
		}
	}

	account := &domain.Account{
		ID:        newAccount.ID,
		FirstName: newAccount.FirstName,
		LastName:  newAccount.LastName,
		Email:     newAccount.Email,
		Password:  newAccount.Password,
	}

	return account, u.repo.Update(account)
}

func (u *AccountUsecase) Delete(executor *domain.Account, id int) error {
	if executor.ID != id {
		return &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrForbidden,
		}
	}
	return u.repo.Delete(id)
}

func (u *AccountUsecase) Register(dto domain.RegistrationParams) (*domain.Account, error) {

	account := domain.NewAccount(dto)

	id, err := u.repo.Create(account)
	account.ID = id

	return account, err
}

func (u *AccountUsecase) Login(email, password string) (*domain.Account, error) {
	account, err := u.repo.GetByEmail(email)

	if err != nil {
		return nil, err
	}

	if account.Password == password {
		return account, nil
	}

	return nil, &domain.ApplicationError{
		OriginalError: nil,
		SimplifiedErr: domain.ErrInvalidInput,
		Description:   "invalid credentials",
	}
}
