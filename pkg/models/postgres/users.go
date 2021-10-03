package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"toDoGolangProject/pkg/models"
)

const (
	insertUserSql     = "INSERT INTO users (name, email, password) VALUES($1, $2, $3)"
	selectUser        = "SELECT id, password FROM users WHERE email = $1 AND active = TRUE"
	selectUserAllInfo = `SELECT id, name, email, active FROM users WHERE id = $1`
)

type UserModel struct {
	Pool *pgxpool.Pool
}

var UserId int

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	_, err = m.Pool.Exec(context.Background(), insertUserSql, name, email, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	row := m.Pool.QueryRow(context.Background(), selectUser, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	UserId = id
	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}
	err := m.Pool.QueryRow(context.Background(), selectUserAllInfo, id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}
