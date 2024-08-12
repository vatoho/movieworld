package userrepo

import (
	"database/sql"
	"errors"
	auth "kinopoisk/service_auth/proto"
)

type UserRepo interface {
	LoginRepo(username, password string) (*auth.User, error)
	RegisterRepo(username, password string) (*auth.User, error)
	FindUserByUsername(username string) (*auth.User, error)
}

type UserRepoMySQL struct {
	db *sql.DB
}

func NewUserRepoMySQL(db *sql.DB) *UserRepoMySQL {
	return &UserRepoMySQL{
		db: db,
	}
}

func (u *UserRepoMySQL) LoginRepo(username, password string) (*auth.User, error) {
	foundUser := &auth.User{}
	err := u.db.
		QueryRow("SELECT id, username FROM users WHERE username = ? AND password = ?", username, password).
		Scan(&foundUser.ID, &foundUser.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}

func (u *UserRepoMySQL) RegisterRepo(username, password string) (*auth.User, error) {
	res, err := u.db.Exec(
		"INSERT INTO users (`username`, `password`) VALUES (?, ?)",
		username,
		password,
	)
	if err != nil {
		return nil, err
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &auth.User{
		ID:       uint64(userID),
		Username: username,
	}, nil
}

func (u *UserRepoMySQL) FindUserByUsername(username string) (*auth.User, error) {
	foundUser := &auth.User{}
	err := u.db.
		QueryRow("SELECT id, username FROM users WHERE username = ?", username).
		Scan(&foundUser.ID, &foundUser.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}
