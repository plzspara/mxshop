package form

type PasswordLoginForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required"`
	Password string `json:"password" form:"password" binding:"required,min=3,max=10"`
}
