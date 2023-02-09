package usecase

import (
	"animal-chipization/internal/domain"
	"errors"
)

type accountRepository interface {
	InsertAccount(account *domain.Account) (*domain.Account, error)
	GetAccountByID(id int32) (*domain.Account, error)
	GetAccountByEmail(email string) (*domain.Account, error)
	SearchAccount(firstName, lastName, email string, from, size int) []*domain.Account
	FullUpdateAccount(newAccount *domain.Account) error
	DeleteAccount(accoundID int) error
}

type AccountUsecase struct {
	repo accountRepository
}

func NewAccountUsecase(repo accountRepository) *AccountUsecase {
	return &AccountUsecase{
		repo: repo,
	}
}

func (u *AccountUsecase) CreateUser(firstName, lastName, email, password string) (*domain.Account, error) {
	account := domain.NewAccount(firstName, lastName, email, password)

	return u.repo.InsertAccount(account)
}

func (u *AccountUsecase) GetAccountByID(id int32) (*domain.Account, error) {
	return u.repo.GetAccountByID(id)
}

func (u *AccountUsecase) Login(email, password string) (*domain.Account, error) {
	account, err := u.repo.GetAccountByEmail(email)

	if err != nil {
		return nil, err
	}

	if account.Password == password {
		return account, nil
	}

	return nil, errors.New("invalid credentials")
}

func (u *AccountUsecase) SearchAccount(firstName, lastName, email string, from, size int) []*domain.AccountResponse {
	accounts := u.repo.SearchAccount(firstName, lastName, email, from, size)

	var responseAccounts []*domain.AccountResponse

	for _, account := range accounts {
		responseAccounts = append(responseAccounts, domain.MakeAccountResponse(*account))
	}

	return responseAccounts
}

func (u *AccountUsecase) FullUpdateAccount(currentAccount *domain.Account, newAccount *domain.Account) (*domain.AccountResponse, error) {
	newAccount.ID = currentAccount.ID

	if err := u.repo.FullUpdateAccount(newAccount); err != nil {
		return nil, err
	}

	return domain.MakeAccountResponse(*newAccount), nil
}

func (u *AccountUsecase) DeleteAccount(accoundID int) error {
	return u.repo.DeleteAccount(accoundID)
}
