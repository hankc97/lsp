package tcpserver

import (
	"strings"
	"regexp"
	"fmt"
)

func getBlacklisted(text string) error {

	blacklist := []string{
		"foo",
		"bar",
		"baz",
	}
	regexpString := fmt.Sprintf("\\b(%s)\\b", strings.Join(blacklist, "|"))
	valid, err := regexp.MatchString(regexpString, text)
	if err != nil {
		return err
	}

	if valid {
		fmt.Println("valid")
	}
	return nil
}