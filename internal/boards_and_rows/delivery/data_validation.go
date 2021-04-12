package delivery

import (
	"regexp"

	"github.com/asaskevich/govalidator"
)

const nameRegexpr = "^([а-яА-Яa-zA-Z0-9!_?]{1,30})$"
const descriptionRegexpr = "^([а-яА-Яa-zA-Z0-9][а-яА-Яa-zA-Z0-9 ()?,.!:;\'\'\"\"]*)$"

func init() {
	govalidator.TagMap["name"] = govalidator.Validator(func(str string) bool {
		result, err := regexp.MatchString(nameRegexpr, str)
		if err != nil {
			return false
		}
		return result
	})
	govalidator.TagMap["description"] = govalidator.Validator(func(str string) bool {
		result, err := regexp.MatchString(descriptionRegexpr, str)
		if err != nil {
			return false
		}
		return result
	})
}
