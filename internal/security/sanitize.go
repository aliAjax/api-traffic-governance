package security

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(?i)(authorization|cookie|token|secret|password)[:=][^\s,;]+`)

func Redact(s string) string { return re.ReplaceAllString(s, "$1=[REDACTED]") }
func Headers(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		lk := strings.ToLower(k)
		if lk == "authorization" || lk == "cookie" || strings.Contains(lk, "token") {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}
func ValidateExperimentData(s string) bool { return !strings.ContainsAny(s, "\x00\r\n") }
