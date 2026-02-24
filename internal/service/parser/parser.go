package parser

import (
	"log"
	"net/url"
	"strings"
)

func ParseURL(input string) (string, error) {
	u, err := url.ParseRequestURI(input)
	if err != nil {
		log.Printf("Invalid URL")
		return "", err
	}
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	if len(u.Path) > 1 && strings.HasSuffix(u.Path, "/") {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}
	return u.String(), nil
}
