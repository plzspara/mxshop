package form

type PasswordLoginForm struct {
	Mobile    string `json:"mobile" form:"mobile" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required,min=3,max=10"`
	Captcha   string `json:"captcha" form:"captcha" binding:"required,len=5"`
	CaptchaId string `json:"captchaid" form:"captchaid" binding:"required"`
}
