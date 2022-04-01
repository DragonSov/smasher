package services

import (
	"errors"
	"github.com/DragonSov/smasher/server/domain/Users"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type DefaultUserService struct {
	repo Users.UserRepositoryDb
}

type UserService interface {
	CreateUser(Users.UserModel) (*Users.UserModel, error)
	SelectUserByUUID(uuid.UUID) (*Users.UserModel, error)
	SelectUserByLogin(string) (*Users.UserModel, error)
	SignIn(string, string) (*string, error)
}

func NewUserService(repository Users.UserRepositoryDb) DefaultUserService {
	return DefaultUserService{repo: repository}
}

func HashPassword(password string) (string, error) {
	// Hashing a password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CheckPassword(hash, password string) bool {
	// Checking password hash and password for a match
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SignJWT(sub uuid.UUID) (*string, error) {
	// Duration calculation and creation of a JWT token
	expireDuration := 20 * time.Minute
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expireDuration).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   sub.String(),
	})

	// Signing a token
	signedString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return nil, err
	}

	return &signedString, err
}

func DecodeJWT(tokenString string) (*jwt.StandardClaims, error) {
	// Creating standard claims and parsing the token
	claims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func (db DefaultUserService) CreateUser(u Users.UserModel) (*Users.UserModel, error) {
	// Hashing user password
	password, err := HashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = password

	// Creating user
	user, err := db.repo.CreateUser(u)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db DefaultUserService) SelectUserByUUID(id uuid.UUID) (*Users.UserModel, error) {
	// Selecting user by ID
	user, err := db.repo.SelectUserByUUID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db DefaultUserService) SelectUserByLogin(login string) (*Users.UserModel, error) {
	// Selecting user by login
	user, err := db.repo.SelectUserByLogin(login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db DefaultUserService) SignIn(login, password string) (*string, error) {
	// Check if the user exists
	user, err := db.SelectUserByLogin(login)
	if err != nil || user == nil {
		return nil, errors.New("Invalid username or password")
	}

	// Checking user password
	correctPassword := CheckPassword(user.Password, password)
	if !correctPassword {
		return nil, errors.New("Invalid username or password")
	}

	// Signing JWT token (access)
	token, err := SignJWT(user.ID)
	if err != nil {
		return nil, errors.New("An unexpected error has occurred")
	}

	return token, nil
}
