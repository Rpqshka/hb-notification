package service

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	hb "hb-notification"
	"hb-notification/pkg/repository"
	"net/smtp"
	"os"
	"time"
)

type NotificationService struct {
	repo repository.Notification
}

func NewNotificationService(repo repository.Notification) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) GetUsers() ([]hb.UserBirthday, error) {
	return s.repo.GetUsers()
}

func (s *NotificationService) Subscribe(userId, subscribeId int) error {
	return s.repo.Subscribe(userId, subscribeId)
}
func (s *NotificationService) Unsubscribe(userId, subscribeId int) error {
	return s.repo.Unsubscribe(userId, subscribeId)
}

func (s *NotificationService) GetSubscriptions(userId int) ([]hb.UserBirthday, error) {
	return s.repo.GetSubscriptions(userId)
}

func (s *NotificationService) CheckBirthday() {
	users, err := s.repo.GetUsers()
	if err != nil {
		logrus.Errorf("error fetching users: %s", err.Error())
		return
	}

	today := time.Now().Format("02-01")

	for _, user := range users {
		userDoB, err := time.Parse("02-01-2006", user.DoB)
		if err != nil {
			logrus.Errorf("error parsing user DoB: %s", err.Error())
			return
		}

		dayBeforeDoB := userDoB.AddDate(0, 0, -1).Format("02-01")

		if dayBeforeDoB == today {
			logrus.Printf("Tomorrow is %s's birthday! Send notification.\n", user.NickName)
			emails, err := s.repo.GetEmailsForNotification(user.Id)
			if err != nil {
				logrus.Errorf("error fetching user emails: %s", err.Error())
				return
			}
			err = sendNotification(user, emails)
			if err != nil {
				logrus.Errorf("error sending notification: %s", err.Error())
				return
			}
		}
	}
}

func sendNotification(user hb.UserBirthday, toEmails []string) error {
	message := fmt.Sprintf("Напоминание о ДР!!!!\r\n\r\nУ %s совсем скоро день рождения (%s)",
		user.NickName, user.DoB)

	from := os.Getenv("SMTP_EMAIL_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	url := os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT")
	auth := smtp.PlainAuth("", from, password, host)

	for _, toEmail := range toEmails {
		to := []string{toEmail}
		err := smtp.SendMail(
			url,
			auth,
			from,
			to,
			[]byte(message),
		)
		if err != nil {
			return errors.New("error with sending recovery mail")
		}
	}
	return nil
}
