package utils

import (
	"regexp"

	"github.com/XZ0730/hertz-scaffold/pkg/constants"
)

func CheckCardID(id string) bool {
	r := regexp.MustCompile(constants.CardIdRegexp)
	return r.MatchString(id)
}

func CheckPhone(id string) bool {
	r := regexp.MustCompile(constants.PhoneRegexp)
	return r.MatchString(id)
}
