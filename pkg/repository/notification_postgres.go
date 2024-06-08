package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	hb "hb-notification"
)

type NotificationPostgres struct {
	db *sqlx.DB
}

func NewNotificationPostgres(db *sqlx.DB) *NotificationPostgres {
	return &NotificationPostgres{db: db}
}

func (r *NotificationPostgres) GetUsers() ([]hb.UserBirthday, error) {
	var users []hb.UserBirthday
	query := fmt.Sprintf("SELECT id, nickname, email, dob FROM %s", usersTable)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user hb.UserBirthday
		err := rows.Scan(&user.Id, &user.NickName, &user.Email, &user.DoB)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *NotificationPostgres) Subscribe(userId, subscribeId int) error {
	var userExists bool
	userCheckQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", usersTable)
	err := r.db.QueryRow(userCheckQuery, subscribeId).Scan(&userExists)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user does not exist")
	}

	var existingId int
	checkQuery := fmt.Sprintf("SELECT id FROM %s WHERE user_id = $1 AND subscribed_to_id = $2", subscriptionsTable)
	if err = r.db.QueryRow(checkQuery, userId, subscribeId).Scan(&existingId); err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingId != 0 {
		return errors.New("already subscribed")
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (user_id, subscribed_to_id) VALUES ($1, $2)", subscriptionsTable)
	_, err = r.db.Exec(insertQuery, userId, subscribeId)
	if err != nil {
		return err
	}
	return nil
}

func (r *NotificationPostgres) Unsubscribe(userId, subscribeId int) error {
	var userExists bool
	userCheckQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", usersTable)
	err := r.db.QueryRow(userCheckQuery, subscribeId).Scan(&userExists)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user does not exist")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND subscribed_to_id = $2", subscriptionsTable)
	result, err := r.db.Exec(query, userId, subscribeId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("cannot unsubscribe")
	}

	return nil
}

func (r *NotificationPostgres) GetSubscriptions(userId int) ([]hb.UserBirthday, error) {
	var users []hb.UserBirthday
	query := fmt.Sprintf(`
		SELECT u.id, u.nickname, u.email, u.dob
		FROM %s u
		JOIN %s s ON u.id = s.subscribed_to_id
		WHERE s.user_id = $1
	`, usersTable, subscriptionsTable)

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user hb.UserBirthday
		err := rows.Scan(&user.Id, &user.NickName, &user.Email, &user.DoB)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *NotificationPostgres) GetEmailsForNotification(subscribeId int) ([]string, error) {
	var emails []string
	query := fmt.Sprintf(`
		SELECT u.email
		FROM %s u
		JOIN %s s ON u.id = s.user_id
		WHERE s.subscribed_to_id = $1
	`, usersTable, subscriptionsTable)

	rows, err := r.db.Query(query, subscribeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}
