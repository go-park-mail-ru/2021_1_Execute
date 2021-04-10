package http

import (
	"regexp"

	"github.com/asaskevich/govalidator"
)

const passwordRegexpr = "^([а-яА-Яa-zA-Z0-9!_?]{6,30})$"
const usernameRegexpr = "^([а-яА-Яa-zA-Z0-9 ]{3,30})$"

func init() {
	govalidator.TagMap["password"] = govalidator.Validator(func(str string) bool {
		result, err := regexp.MatchString(passwordRegexpr, str)
		if err != nil {
			return false
		}
		return result
	})
	govalidator.TagMap["username"] = govalidator.Validator(func(str string) bool {
		result, err := regexp.MatchString(usernameRegexpr, str)
		if err != nil {
			return false
		}
		return result
	})
}
