package service

import (
	hb "hb-notification"
	"hb-notification/pkg/repository"
)

type Authorization interface {
	CreateUser(user hb.User) (int, error)
	CheckNickNameAndEmail(nickname, email string) (int, error)
	GetPasswordHash(nickname string) (string, error)
	GenerateToken(nickname, passwordHash string) (string, error)
	ParseToken(accessToken string) (int, error)
	CheckEmail(email string) (int, error)
}

type Notification interface {
	GetUsers() ([]hb.UserBirthday, error)
	Subscribe(userId, subscribeId int) error
	Unsubscribe(userId, subscribeId int) error
	GetSubscriptions(userId int) ([]hb.UserBirthday, error)
	CheckBirthday()
}

type Service struct {
	Authorization
	Notification
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Notification:  NewNotificationService(repos.Notification),
	}
}
