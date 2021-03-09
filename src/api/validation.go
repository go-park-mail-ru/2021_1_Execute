package api

import "regexp"

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9]" +
	"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

const minPasswordLen = 6
const maxPasswordLen = 50
const minEmailLen = 3
const maxEmailLen = 254

func IsEmailValid(email string) bool {
	if len(email) < minEmailLen || len(email) > maxEmailLen {
		return false
	}
	return emailRegex.MatchString(email)
}

func IsPasswordValid(password string) bool {
	return len(password) >= minPasswordLen && len(password) <= maxPasswordLen
}
