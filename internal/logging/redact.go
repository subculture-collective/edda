package logging

import (
	"net/url"
	"regexp"
	"strings"
)

// redactedPassword replaces a masked password.
const redactedPassword = "xxxxx"

// secretQueryParams are query parameters dropped from logged URLs, compared
// without regard to case.
var secretQueryParams = map[string]bool{
	"password": true, "sslpassword": true, "passwd": true, "pass": true,
	"token": true, "access_token": true, "refresh_token": true, "id_token": true,
	"api_key": true, "apikey": true, "key": true,
	"secret": true, "client_secret": true,
	"auth": true, "authorization": true, "sig": true, "signature": true,
}

var (
	urlPattern = regexp.MustCompile("[a-zA-Z][a-zA-Z0-9+.-]*://[^\\s\"'`<>]+")
	// keywordPasswordPattern matches password=value in a key/value PostgreSQL
	// DSN, with the value bare or single-quoted.
	keywordPasswordPattern = regexp.MustCompile(`(?i)\b((?:ssl)?password)=('(?:[^'\\]|\\.)*'|[^\s'` + "`" + `]+)`)
	// userinfoPattern matches the userinfo of a URL that url.Parse rejects.
	userinfoPattern = regexp.MustCompile(`^([a-zA-Z][a-zA-Z0-9+.-]*://)[^/?#]*@`)
)

// RedactURL returns raw, a configured endpoint or DSN, safe to log: userinfo
// is removed, except that a postgres:// or postgresql:// URL keeps its user
// name and has its password masked; secret-bearing query parameters are
// dropped; and a key/value DSN has its password masked.
func RedactURL(raw string) string {
	if !strings.Contains(raw, "://") {
		return keywordPasswordPattern.ReplaceAllString(raw, "${1}="+redactedPassword)
	}
	u, err := url.Parse(raw)
	if err != nil {
		return userinfoPattern.ReplaceAllString(raw, "$1")
	}
	changed := false
	if u.User != nil {
		changed = true
		_, hasPassword := u.User.Password()
		switch {
		case (u.Scheme == "postgres" || u.Scheme == "postgresql") && hasPassword:
			u.User = url.UserPassword(u.User.Username(), redactedPassword)
		case u.Scheme == "postgres" || u.Scheme == "postgresql":
		default:
			u.User = nil
		}
	}
	if u.RawQuery != "" {
		kept := make([]string, 0, strings.Count(u.RawQuery, "&")+1)
		for _, pair := range strings.Split(u.RawQuery, "&") {
			name, _, _ := strings.Cut(pair, "=")
			if decoded, err := url.QueryUnescape(name); err == nil {
				name = decoded
			}
			if secretQueryParams[strings.ToLower(name)] {
				changed = true
			} else {
				kept = append(kept, pair)
			}
		}
		u.RawQuery = strings.Join(kept, "&")
	}
	if !changed {
		return raw
	}
	return u.String()
}

// Redact returns text with RedactURL applied to every URL in it and the
// password of any key/value DSN masked. The log handler applies it to every
// message and field value, so an error that quotes a configured endpoint or
// database URL does not reach a log destination with its credentials.
func Redact(text string) string {
	if strings.Contains(text, "://") {
		text = urlPattern.ReplaceAllStringFunc(text, RedactURL)
	}
	if strings.Contains(strings.ToLower(text), "password=") {
		text = keywordPasswordPattern.ReplaceAllString(text, "${1}="+redactedPassword)
	}
	return text
}
