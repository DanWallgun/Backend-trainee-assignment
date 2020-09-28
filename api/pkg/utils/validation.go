package utils

import "net/url"

// ValidateURL validates url
func ValidateURL(URLString string) bool {
	_, err := url.ParseRequestURI(URLString)
	return err == nil
}
