package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func comparePasswordHash(hash, inputPassword string) error {
	hashedPassword := []byte(hash)
	inputPasswordBytes := []byte(inputPassword)

	return bcrypt.CompareHashAndPassword(hashedPassword, inputPasswordBytes)
}

func validateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

func validateDoB(dobStr string) (string, error) {
	layout := "02-01-2006"

	dob, err := time.Parse(layout, dobStr)
	if err != nil {
		return "", err
	}

	now := time.Now()
	maxDate := now.AddDate(-120, 0, 0)

	if dob.Before(maxDate) {
		return "", errors.New("incorrect date of birth")
	}

	dobFormat := dob.Format(layout)
	return dobFormat, nil
}

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	//parse token
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
