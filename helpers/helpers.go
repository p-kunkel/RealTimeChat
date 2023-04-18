package helpers

import (
	"log"
	"net"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func IsValidEmail(e string) bool {
	if len(e) < 5 && len(e) > 254 {
		return false
	}
	re, err := regexp.Compile(`[^\w@.-]`)
	if err != nil {
		log.Fatal(err)
	}

	if re.MatchString(e) {
		return false
	}

	parts := strings.Split(e, "@")
	if len(parts) != 2 {
		return false
	}
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}
