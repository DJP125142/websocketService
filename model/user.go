package model

type UserInfo struct {
	UserInfo User `json:"userInfo"`
}

type User struct {
	UserId   int    `json:"userId"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}
