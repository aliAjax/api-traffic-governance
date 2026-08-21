package domain

import (
	"net/url"
	"regexp"
	"strings"
)

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9.-]{1,62}$`)

func ValidName(v string) bool { return namePattern.MatchString(v) }
func ValidURL(v string) bool {
	u, e := url.Parse(v)
	return e == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}
func NormalizePath(v string) string {
	v = strings.TrimSpace(v)
	if v == "" || v[0] != '/' {
		v = "/" + v
	}
	for strings.Contains(v, "//") {
		v = strings.ReplaceAll(v, "//", "/")
	}
	return v
}
func HeaderMatch(headers map[string]string, provided map[string]string) bool {
	if headers == nil {
		return true
	}
	for k, v := range headers {
		found := false
		for pk, pv := range provided {
			if strings.EqualFold(k, pk) && v == pv {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
