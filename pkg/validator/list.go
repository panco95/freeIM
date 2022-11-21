package validator

import (
	"regexp"

	"github.com/globalsign/mgo/bson"
	validator "github.com/go-playground/validator/v10"
)

var (
	dateRegex   = regexp.MustCompile(`^\d{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$`)
	timeRegex   = regexp.MustCompile(`^(0\d|1\d|2[0-3]):[0-5]\d:[0-5]\d$`)
	mobileRegex = regexp.MustCompile(`^1[3456789]\d{9}$`)
)

// ValidateObjectID ...
func ValidateObjectID(fl validator.FieldLevel) bool {
	// currentField, _, _ := fl.GetStructFieldOK()
	// table := currentField.Type().Name()
	// column := fl.FieldName()
	value := fl.Field().String()

	if value == "" {
		return true
	}

	return bson.IsObjectIdHex(value)
}

// ValidateDate ...
func ValidateDate(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	return dateRegex.MatchString(value)
}

// ValidateTime ...
func ValidateTime(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	return timeRegex.MatchString(value)
}

func ValidateMobile(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value != "" {
		if len(value) > 30 {
			return false
		}
	}
	return true
}

func ValidateRealMobile(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return mobileRegex.MatchString(value)
}

func ValidateIdCard(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value != "" {
		res, err := regexp.Match("(^\\d{15}$)|(^\\d{17}(\\d|X|x)$)", []byte(value))
		if err != nil {
			return false
		}
		return res
	}
	return true
}
