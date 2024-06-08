package repository

import (
	"github.com/jmoiron/sqlx"
	hb "hb-notification"
)

type Authorization interface {
	CreateUser(user hb.User) (int, error)
	CheckNickNameAndEmail(nickname, email string) (int, error)
	GetPasswordHash(nickname string) (string, error)
	GetUser(nickname, password string) (hb.User, error)
	CheckEmail(email string) (int, error)
}

type Notification interface {
	GetUsers() ([]hb.UserBirthday, error)
	Subscribe(userId, subscribeId int) error
	Unsubscribe(userId, subscribeId int) error
	GetSubscriptions(userId int) ([]hb.UserBirthday, error)
	GetEmailsForNotification(subscribeId int) ([]string, error)
}

type Repository struct {
	Authorization
	Notification
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Notification:  NewNotificationPostgres(db),
	}
}
