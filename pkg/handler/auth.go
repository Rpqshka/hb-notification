package handler

import (
	"github.com/gin-gonic/gin"
	hb "hb-notification"
	"net/http"
)

type signUpInput struct {
	NickName        string `json:"nickname" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	DoB             string `json:"dayOfBirthday" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !validateEmail(input.Email) {
		newErrorResponse(c, http.StatusBadRequest, "enter a different email")
		return
	}

	if input.Password != input.PasswordConfirm {
		newErrorResponse(c, http.StatusBadRequest, "passwords does not match")
		return
	}

	dob, err := validateDoB(input.DoB)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if _, err = h.services.Authorization.CheckNickNameAndEmail(input.NickName, input.Email); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := hb.User{
		NickName: input.NickName,
		Email:    input.Email,
		Password: input.Password,
		DoB:      dob,
	}

	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	NickName string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	passwordHash, err := h.services.GetPasswordHash(input.NickName)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = comparePasswordHash(passwordHash, input.Password); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.NickName, passwordHash)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
