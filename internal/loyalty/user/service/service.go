package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	bstorage "github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/balance/storage"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/user/model"
	ustorage "github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty/user/storage"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/db"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound informs that there is no user with the specified login.
	ErrUserNotFound = errors.New("user not found")
	// ErrWrongPassword informs that the password provided is wrong.
	ErrWrongPassword = errors.New("wrong password")
	// ErrDuplicateLogin informs that the login already exists.
	ErrDuplicateLogin = errors.New("login already exists")
)

// hashPassword hashes password for further storage.
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// checkPassword compares already hashed password with a provided non-hashed one.
func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

type userStorage interface {
	AddUser(ctx context.Context, user model.IntUser) (int64, error)
	FindByLogin(ctx context.Context, login string) (model.IntUser, error)
}

type balanceStorage interface {
	AddBalance(ctx context.Context, userID int64) error
}

// UserService is the service layer for user logic.
type UserService struct {
	uStorage userStorage
	bStorage balanceStorage
	jwtGen   *jwt.Generator
}

// New creates a new user service.
func New(db db.DB, jwtGen *jwt.Generator) UserService {
	return UserService{
		uStorage: ustorage.NewStorage(db),
		bStorage: bstorage.NewStorage(db),
		jwtGen:   jwtGen,
	}
}

// Login logs the user in.
func (us *UserService) Login(ctx context.Context, req model.LoginReq) (string, error) {
	user, err := us.uStorage.FindByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("could not login: %w", err)
	}

	if !checkPassword(user.Password, req.Password) {
		return "", ErrWrongPassword
	}

	tok, err := us.jwtGen.GetToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("error getting token: %w", err)
	}

	return tok, nil
}

// Register registers a new user.
func (us *UserService) Register(ctx context.Context, login, password string) (string, error) {
	encryptedPassw, err := hashPassword(password)
	if err != nil {
		return "", fmt.Errorf("error encrypting password: %w", err)
	}

	user := model.IntUser{
		Login:    login,
		Password: encryptedPassw,
	}

	_, err = us.uStorage.FindByLogin(ctx, login)
	if err == nil {
		return "", ErrDuplicateLogin
	}

	userID, err := us.uStorage.AddUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("error adding new user: %w", err)
	}

	err = us.bStorage.AddBalance(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("error adding balance for new user: %w", err)
	}

	tok, err := us.jwtGen.GetToken(userID)
	if err != nil {
		return "", fmt.Errorf("error getting token: %w", err)
	}

	return tok, nil
}
