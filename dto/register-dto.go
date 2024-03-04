package dto

type RegisterDTO struct {
	ID       uint64 `json:"id" form:"id"`
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
	Profile  string `json:"profile" form:"profile"`
	Jk       string `json:"jk"`
}
