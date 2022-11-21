package validator

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin/binding"
	en_locales "github.com/go-playground/locales/en"
	zh_locales "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni *ut.UniversalTranslator
	// ErrInvalidValidator ...
	ErrInvalidValidator = errors.New("invalid validator")

	registerMu sync.Mutex
)

// InitBindingValidator ...
func InitBindingValidator() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var err error
		uni, err = Register(v)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrInvalidValidator
}

// GetUniversalTranslator ...
func GetUniversalTranslator() (*ut.UniversalTranslator, error) {
	if uni != nil {
		return uni, nil
	}

	err := InitBindingValidator()
	if err != nil {
		return nil, err
	}

	return uni, nil
}

// Register ...
func Register(v *validator.Validate) (*ut.UniversalTranslator, error) {
	registerMu.Lock()
	defer registerMu.Unlock()

	err := v.RegisterValidation("objectid", ValidateObjectID)
	if err != nil {
		return nil, err
	}

	err = v.RegisterValidation("date", ValidateDate)
	if err != nil {
		return nil, err
	}

	err = v.RegisterValidation("time", ValidateTime)
	if err != nil {
		return nil, err
	}

	err = v.RegisterValidation("idcard", ValidateIdCard)
	if err != nil {
		return nil, err
	}

	err = v.RegisterValidation("mobile", ValidateMobile)
	if err != nil {
		return nil, err
	}

	err = v.RegisterValidation("real_mobile", ValidateRealMobile)
	if err != nil {
		return nil, err
	}

	v.RegisterAlias("weekday", "oneof=Sunday Monday Tuesday Wednesday Thursday Friday Saturday")

	en := en_locales.New()
	zh := zh_locales.New()
	uni := ut.New(en, en, zh)

	enTrans, _ := uni.GetTranslator("en")
	err = en_translations.RegisterDefaultTranslations(v, enTrans)
	if err != nil {
		return nil, err
	}
	zhTrans, _ := uni.GetTranslator("zh")
	err = zh_translations.RegisterDefaultTranslations(v, zhTrans)
	if err != nil {
		return nil, err
	}

	err = v.RegisterTranslation("objectid", enTrans, func(ut ut.Translator) error {
		return ut.Add("objectid", "{0} must be ObjectID (24hex)!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("objectid", fe.Field())
		return t
	})
	if err != nil {
		return nil, err
	}

	err = v.RegisterTranslation("objectid", zhTrans, func(ut ut.Translator) error {
		return ut.Add("objectid", "{0} 必须是 ObjectID (24位十六进制)!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("objectid", fe.Field())
		return t
	})
	if err != nil {
		return nil, err
	}

	err = v.RegisterTranslation("date", enTrans, func(ut ut.Translator) error {
		return ut.Add("date", "{0} must be date (YYYY-MM-DD)!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("date", fe.Field())
		return t
	})
	if err != nil {
		return nil, err
	}

	err = v.RegisterTranslation("date", zhTrans, func(ut ut.Translator) error {
		return ut.Add("date", "{0} 必须是日期格式 (YYYY-MM-DD)!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("date", fe.Field())
		return t
	})
	if err != nil {
		return nil, err
	}

	err = v.RegisterTranslation("time", enTrans, func(ut ut.Translator) error {
		return ut.Add("time", "{0} must be time (HH:mm:ss)!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("time", fe.Field())
		return t
	})
	if err != nil {
		return nil, err
	}

	err = v.RegisterTranslation("time", zhTrans, func(ut ut.Translator) error {
		return ut.Add("time", "{0} 必须是时间格式 (HH:mm:ss)!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("time", fe.Field())
		return t
	})
	if err != nil {
		return nil, err
	}

	err = v.RegisterTranslation("weekday", enTrans, func(ut ut.Translator) error {
		return ut.Add("weekday", "{0} must be weekday (oneof Sunday/Monday/Tuesday/Wednesday/Thursday/Friday/Saturday)!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("weekday", fe.Field())
		return t
	})
	if err != nil {
		return nil, err
	}

	err = v.RegisterTranslation("weekday", zhTrans, func(ut ut.Translator) error {
		return ut.Add("weekday", "{0} 必须是星期几 (oneof Sunday/Monday/Tuesday/Wednesday/Thursday/Friday/Saturday)!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("weekday", fe.Field())
		return t
	})
	if err != nil {
		return nil, err
	}

	return uni, nil
}
