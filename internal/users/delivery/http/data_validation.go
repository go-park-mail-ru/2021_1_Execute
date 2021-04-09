package http

import (
	"regexp"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.ParamTagRegexMap["password"] = regexp.MustCompile("^([а-яА-Яa-zA-Z0-9!_?]{6,30})$")
	govalidator.ParamTagRegexMap["username"] = regexp.MustCompile("^([а-яА-Яa-zA-Z0-9 ]{3,30})$")
}
