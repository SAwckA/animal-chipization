package psql

import (
	"animal-chipization/internal/domain"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	accountTable = "public.account"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) InsertAccount(account *domain.Account) (*domain.Account, error) {
	query := fmt.Sprintf(`insert into %s(firstName, lastName, email, password) values ($1, $2, $3, $4) returning id`, accountTable)

	row := r.db.QueryRow(query, account.FirstName, account.LastName, account.Email, account.Password)

	var id int32
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	account.ID = id
	return account, nil
}

func (r *AccountRepository) GetAccountByID(id int32) (*domain.Account, error) {
	query := fmt.Sprintf(`select id, firstName, lastName, email from %s where id=$1`, accountTable)

	var account domain.Account

	if err := r.db.Get(&account, query, id); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) GetAccountByEmail(email string) (*domain.Account, error) {
	query := fmt.Sprintf(`select id, firstname, lastname, email, password from %s where email=$1`, accountTable)

	var account domain.Account

	if err := r.db.Get(&account, query, email); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) SearchAccount(firstName, lastName, email string, from, size int) []*domain.Account {
	query := fmt.Sprintf(`
		select * from %s 
		where 

		(LOWER(firstname) like '%%' || LOWER($1) || '%%') 
		and 
		(LOWER(lastname) like '%%' || LOWER($2) || '%%')
		and 
		(LOWER(email) like '%%' || LOWER($3) || '%%') 

		LIMIT $4
		OFFSET $5;
	`, accountTable)

	logrus.Warn(query)

	var accounts []*domain.Account

	err := r.db.Select(&accounts, query, firstName, lastName, email, size, from)

	if err != nil {
		return nil
	}

	return accounts
}

func (r *AccountRepository) FullUpdateAccount(newAccount *domain.Account) error {

	query := fmt.Sprintf(`
		update %s
		set firstname = $1,
			lastname = $2,
			email = $3,
			password = $4
		where id = $5
		`, accountTable)

	_, err := r.db.Exec(query, newAccount.FirstName, newAccount.LastName, newAccount.Email, newAccount.Password, newAccount.ID)

	return err
}

func (r *AccountRepository) DeleteAccount(accoundID int) error {

	query := fmt.Sprintf(`
	delete from %s
	where id = $1
	`, accountTable)

	_, err := r.db.Exec(query, accoundID)

	return err
}
