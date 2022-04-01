package Users

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"strings"
)

type UserRepositoryDb struct {
	Client *sqlx.DB
}

func NewUserRepositoryDb(client *sqlx.DB) UserRepositoryDb {
	return UserRepositoryDb{
		Client: client,
	}
}

func (repo *UserRepositoryDb) CreateUser(u UserModel) (*UserModel, error) {
	// Inserting a new user
	_, err := repo.Client.NamedExec("INSERT INTO users (login, password) VALUES (:login, :password)", u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (repo *UserRepositoryDb) SelectUserByUUID(id uuid.UUID) (*UserModel, error) {
	// Creating user model
	u := UserModel{}

	// Filling the user model
	err := repo.Client.Get(&u, "SELECT * FROM users WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &u, nil
}

func (repo *UserRepositoryDb) SelectUserByLogin(login string) (*UserModel, error) {
	// Creating user model
	u := UserModel{}

	// Filling the user model
	err := repo.Client.Get(&u, "SELECT * FROM users WHERE lower(login) = $1", strings.ToLower(login))
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &u, nil
}
