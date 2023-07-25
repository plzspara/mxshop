package inittialize

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"mxshop-api/global"
	va "mxshop-api/validator"
)

func InitTrans(locale string) error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		zhT := zh.New()
		enT := en.New()
		translator := ut.New(enT, zhT, enT)
		global.Trans, ok = translator.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("GetTranslator error: %v", locale)
		}
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, global.Trans)
		default:
			en_translations.RegisterDefaultTranslations(v, global.Trans)
		}

		_ = v.RegisterValidation("mobile", va.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
	return nil
}
