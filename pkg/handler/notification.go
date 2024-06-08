package handler

import (
	"github.com/gin-gonic/gin"
	hb "hb-notification"
	"net/http"
)

type usersDataResponse struct {
	Users []hb.UserBirthday `json:"users"`
}

type getAllUsersResponse struct {
	Data usersDataResponse `json:"data"`
}

func (h *Handler) getUsers(c *gin.Context) {

	users, err := h.services.Notification.GetUsers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllUsersResponse{
		Data: usersDataResponse{
			Users: users,
		},
	})
}

type subscribeInput struct {
	SubscribeId int `json:"subscribe_id" binding:"required"`
}

func (h *Handler) subscribe(c *gin.Context) {
	var input subscribeInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	if userId == input.SubscribeId {
		newErrorResponse(c, http.StatusBadRequest, "cannot subscribe at yourself")
		return
	}

	if err = h.services.Notification.Subscribe(userId, input.SubscribeId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "subscribed",
	})
}

type unsubscribeInput struct {
	UnsubscribeId int `json:"unsubscribe_id" binding:"required"`
}

func (h *Handler) unsubscribe(c *gin.Context) {
	var input unsubscribeInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Notification.Unsubscribe(userId, input.UnsubscribeId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "unsubscribed",
	})
}

func (h *Handler) getSubscriptions(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.services.Notification.GetSubscriptions(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllUsersResponse{
		Data: usersDataResponse{
			Users: users,
		},
	})
}

type updateCronInput struct {
	CronInterval string `json:"cron_interval"`
}

func (h *Handler) updateCron(c *gin.Context) {
	var input updateCronInput
	
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.Cron.Remove(h.JobID)

	newJobID, err := h.Cron.AddFunc(input.CronInterval, func() { h.services.Notification.CheckBirthday() })
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.JobID = newJobID

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "cron updated",
	})
}
