package api

import "regexp"

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9]" +
	"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(email string) bool {
	const minEmailLen = 3
	const maxEmailLen = 254
	if len(email) < minEmailLen || len(email) > maxEmailLen {
		return false
	}
	return emailRegex.MatchString(email)
}

func IsPasswordValid(password string) bool {
	const minPasswordLen = 6
	const maxPasswordLen = 50
	return len(password) >= minPasswordLen && len(password) <= maxPasswordLen
}
