package psql

import (
	"animal-chipization/internal/domain"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	accountTable                 = "public.account"
	accountEmailUniqueConstraint = "account_email_key"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(account *domain.Account) (int, error) {
	query := fmt.Sprintf(`insert into %s(firstName, lastName, email, password) values ($1, $2, $3, $4) returning id`, accountTable)

	var id int
	err := r.db.Get(&id, query, account.FirstName, account.LastName, account.Email, account.Password)
	if err != nil {
		if strings.Contains(err.Error(), accountEmailUniqueConstraint) {
			return 0, &domain.ApplicationError{
				OriginalError: err,
				SimplifiedErr: domain.ErrAlreadyExist,
				Description:   "account with given email already exists",
			}
		}
		return 0, err
	}

	return id, nil
}

func (r *AccountRepository) GetByID(id int) (*domain.Account, error) {
	query := fmt.Sprintf(`select id, firstName, lastName, email from %s where id=$1`, accountTable)

	var account domain.Account
	if err := r.db.QueryRow(query, id).Scan(&account.ID, &account.FirstName, &account.LastName, &account.Email); err != nil {
		return nil, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "account not found by id",
		}
	}

	return &account, nil
}

func (r *AccountRepository) GetByEmail(email string) (*domain.Account, error) {
	query := fmt.Sprintf(`select id, firstname, lastname, email, password from %s where email=$1`, accountTable)

	var account domain.Account
	if err := r.db.Get(&account, query, email); err != nil {
		return nil, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "account not found by id",
		}
	}

	return &account, nil
}

func (r *AccountRepository) Search(params *domain.SearchAccount) ([]domain.Account, error) {
	var searchQuery []string
	var searchArgs []interface{}
	searchArgs = append(searchArgs, params.Size)
	searchArgs = append(searchArgs, params.From)

	ph := 3
	isSearch := "where"

	if params.FirstName != nil {
		searchQuery = append(searchQuery, fmt.Sprintf("(LOWER(firstname) like '%%' || LOWER($%d) || '%%') ", ph))
		searchArgs = append(searchArgs, params.FirstName)
		ph++
	}

	if params.LastName != nil {
		searchQuery = append(searchQuery, fmt.Sprintf("(LOWER(lastname) like '%%' || LOWER($%d) || '%%')", ph))
		searchArgs = append(searchArgs, params.LastName)
		ph++
	}

	if params.Email != nil {
		searchQuery = append(searchQuery, fmt.Sprintf("(LOWER(email) like '%%' || LOWER($%d) || '%%') ", ph))
		searchArgs = append(searchArgs, params.Email)
	}

	if len(searchQuery) == 0 {
		isSearch = ""
	}

	query := fmt.Sprintf(`
		select id, firstname, lastname, email from %s 
		%s
			%s
		ORDER BY id
		LIMIT $1
		OFFSET $2`,
		accountTable,
		isSearch,
		strings.Join(searchQuery, " and "),
	)

	var accounts []domain.Account
	rows, err := r.db.Query(query, searchArgs...)

	if err != nil {
		return nil, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "database error",
		}
	}

	for rows.Next() {
		var account domain.Account

		err = rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Email)
		accounts = append(accounts, account)
	}

	return accounts, err
}

func (r *AccountRepository) Update(newAccount *domain.Account) error {

	query := fmt.Sprintf(`
		update %s
		set firstname = $1,
			lastname = $2,
			email = $3,
			password = $4
		where id = $5
		`, accountTable)

	_, err := r.db.Exec(query, newAccount.FirstName, newAccount.LastName, newAccount.Email, newAccount.Password, newAccount.ID)
	if err != nil {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrConflict,
			Description:   "account already exist",
		}
	}

	return nil
}

func (r *AccountRepository) Delete(accountID int) error {

	query := fmt.Sprintf(`
	delete from %s
	where id = $1
	`, accountTable)

	_, err := r.db.Exec(query, accountID)
	if err != nil {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "account linked with animal",
		}
	}
	return nil
}
