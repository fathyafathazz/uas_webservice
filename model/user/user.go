package user

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password_hash"`
	Role string `json:"role"`
}