package parser

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func ParseURL(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("EmptyInput")
	}
	u, err := url.Parse(input)
	if err != nil {
		log.Printf("Invalid URL")
		return "", err
	}
	u.Scheme = strings.ToLower(u.Scheme)
	if u.Scheme == "" {
		return "", fmt.Errorf("No scheme\n ")
	}
	u.Host = strings.ToLower(u.Host)
	if u.Scheme == "" {
		return "", fmt.Errorf("No host\n ")
	}
	if len(u.Path) > 1 && strings.HasSuffix(u.Path, "/") {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}
	return u.String(), nil
}
