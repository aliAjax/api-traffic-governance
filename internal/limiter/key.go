package limiter

import (
	"net/http"
	"strings"
)

func Key(req *http.Request, parts ...string) string {
	values := []string{}
	for _, p := range parts {
		switch strings.ToLower(p) {
		case "tenant":
			values = append(values, req.Header.Get("X-Tenant-ID"))
		case "ip":
			values = append(values, strings.Split(req.RemoteAddr, ":")[0])
		case "path":
			values = append(values, req.URL.Path)
		case "method":
			values = append(values, req.Method)
		case "authorization":
			values = append(values, req.Header.Get("Authorization"))
		default:
			values = append(values, req.Header.Get(p))
		}
	}
	return strings.Join(values, "|")
}
func NormalizeKey(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
