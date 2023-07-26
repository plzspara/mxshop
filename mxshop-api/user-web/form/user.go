package form

type PasswordLoginForm struct {
	Mobile    string `json:"mobile" form:"mobile" binding:"required,mobile"`
	Password  string `json:"password" form:"password" binding:"required,min=3,max=10"`
	Captcha   string `json:"captcha" form:"captcha" binding:"required,len=5"`
	CaptchaId string `json:"captchaid" form:"captchaid" binding:"required"`
}

type GetSmsRequest struct {
	Mobile string `json:"mobile" form:"mobile" binding:"required,mobile"`
}

type MobileLoginForm struct {
	Mobile string `json:"mobile" form:"mobile" binding:"required,mobile"`
	Code   string `json:"code" form:"code" binding:"required,len=6"`
}

type UserRegister struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"`
	Password string `json:"password" form:"password" binding:"required,min=3,max=10"`
	Code     string `json:"code" form:"code" binding:"required,len=6"`
}
