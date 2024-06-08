package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	hb "hb-notification"
	"hb-notification/pkg/repository"
	"time"
)

const (
	signingKey = "adfa6464aE"
	tokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user hb.User) (int, error) {
	return s.repo.CreateUser(user)
}

func (s *AuthService) CheckNickNameAndEmail(nickname, email string) (int, error) {
	_, err := s.repo.CheckNickNameAndEmail(nickname, email)
	if err == nil {
		return 0, errors.New("this user already registered")
	}
	return 0, nil
}

func (s *AuthService) GetPasswordHash(nickname string) (string, error) {
	passwordHash, err := s.repo.GetPasswordHash(nickname)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (s *AuthService) GenerateToken(nickname, passwordHash string) (string, error) {
	user, err := s.repo.GetUser(nickname, passwordHash)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) CheckEmail(email string) (int, error) {
	_, err := s.repo.CheckEmail(email)
	if err != nil {
		return 0, errors.New("user does not exists")
	}
	return 0, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(accessToken *jwt.Token) (interface{}, error) {
		if _, ok := accessToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil
}
