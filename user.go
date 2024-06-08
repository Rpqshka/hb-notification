package hb_notification

type User struct {
	Id       int    `json:"id" db:"id"`
	NickName string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	DoB      string `json:"dob" binding:"required"`
}

type UserBirthday struct {
	Id       int    `json:"id" db:"id"`
	NickName string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	DoB      string `json:"dob" binding:"required"`
}
