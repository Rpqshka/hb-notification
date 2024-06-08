package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	hb "hb-notification"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user hb.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (nickname, email, password_hash, dob) VALUES ($1, $2, $3, $4) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.NickName, user.Email, user.Password, user.DoB)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) CheckNickNameAndEmail(nickname, email string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE nickname = $1 OR email = $2", usersTable)
	err := r.db.Get(&id, query, nickname, email)
	return id, err
}

func (r *AuthPostgres) GetUser(nickname, password string) (hb.User, error) {
	var user hb.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE nickname = $1 AND password_hash = $2", usersTable)
	err := r.db.Get(&user, query, nickname, password)
	return user, err
}

func (r *AuthPostgres) GetPasswordHash(nickname string) (string, error) {
	var hash string
	query := fmt.Sprintf("SELECT password_hash FROM %s WHERE nickname = $1", usersTable)
	err := r.db.Get(&hash, query, nickname)
	return hash, err
}

func (r *AuthPostgres) CheckEmail(email string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1", usersTable)
	err := r.db.Get(&id, query, email)
	return id, err
}
