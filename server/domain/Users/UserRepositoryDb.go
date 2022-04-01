package Users

import (
	"github.com/google/uuid"
	"time"
)

type UserModel struct {
	ID        uuid.UUID `db:"id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

type UserRepository interface {
	CreateUser(UserModel) (*UserModel, error)
	SelectUserByUUID(uuid.UUID) (*UserModel, error)
	SelectUserByLogin(string) (*UserModel, error)
}
