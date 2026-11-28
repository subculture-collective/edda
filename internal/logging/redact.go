package logging

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"slices"
	"strings"
)

// redactedPassword replaces a masked password.
const redactedPassword = "xxxxx"

// A query parameter is secret, and dropped from logged URLs, when its name ends
// in one of secretQuerySuffixes or is one of secretQueryParams. Both are
// compared without regard to case. The suffixes cover names such as x-api-key,
// subscription-key, private_token, client_secret, X-Amz-Signature and
// sslpassword.
var (
	secretQuerySuffixes = []string{"key", "token", "secret", "signature", "password"}
	secretQueryParams   = map[string]bool{
		"passwd": true, "pass": true, "auth": true, "authorization": true, "sig": true,
	}
)

var (
	urlPattern = regexp.MustCompile("[a-zA-Z][a-zA-Z0-9+.-]*://[^\\s\"'`<>]+")
	// keywordPasswordPattern matches password=value in a key/value PostgreSQL
	// DSN, with the value bare or single-quoted.
	keywordPasswordPattern = regexp.MustCompile(`(?i)\b((?:ssl)?password)=('(?:[^'\\]|\\.)*'|[^\s'` + "`" + `]+)`)
	// openAuthorityPattern matches a URL that ends the text with nothing to end
	// its host. A trailing ellipsis or space is not part of it. The group is what
	// follows "://".
	openAuthorityPattern = regexp.MustCompile("[a-zA-Z][a-zA-Z0-9+.-]*://([^\\s\"'`<>/?#]+?)(?:\\.\\.\\.|…)?\\s*$")
)

// RedactURL returns raw, a URL, configured endpoint or DSN, safe to log or to
// put in an error: userinfo is removed, except that a postgres:// or
// postgresql:// URL keeps its user name and has its password masked;
// secret query parameters are dropped; and a key/value DSN has its password
// masked. A URL that does not parse, or whose userinfo may hold an unescaped
// /, ? or # (see hidesUserinfo), also loses its user name and whole query.
func RedactURL(raw string) string {
	if !strings.Contains(raw, "://") {
		if keywordPasswordPattern.MatchString(raw) || strings.ContainsAny(raw, " \t\r\n") || !strings.ContainsAny(raw, "@?") {
			return keywordPasswordPattern.ReplaceAllString(raw, "${1}="+redactedPassword)
		}
		// An endpoint with its scheme left out.
		return strings.TrimPrefix(RedactURL("http://"+raw), "http://")
	}
	u, err := url.Parse(raw)
	if err != nil || hidesUserinfo(raw) {
		return redactUnparsed(raw)
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
			if isSecretQueryParam(name) {
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

func isSecretQueryParam(name string) bool {
	name = strings.ToLower(name)
	return secretQueryParams[name] || slices.ContainsFunc(secretQuerySuffixes, func(suffix string) bool {
		return strings.HasSuffix(name, suffix)
	})
}

// hidesUserinfo reports whether an @ after the apparent end of raw's authority
// may be the real end of its userinfo. url.Parse accepts
// "http://user:12/34@host" as host "user", port "12" and path "/34@host". Such
// an @ is read as part of a path or query only when the authority is plainly a
// host (it has a dot, or is localhost) and no colon comes before the @, so no
// password can be there. A valid URL with a port and a later @ is cut too.
func hidesUserinfo(raw string) bool {
	_, rest, found := strings.Cut(raw, "://")
	end := strings.IndexAny(rest, "/?#")
	at := strings.LastIndex(rest, "@")
	if !found || end < 0 || at < end {
		return false
	}
	host := rest[:end]
	plainHost := host == "localhost" || strings.Contains(host, ".")
	return !plainHost || strings.Contains(rest[:at], ":")
}

// redactUnparsed cuts raw, a URL that url.Parse rejects or may misread, down to
// its host and path. It drops everything up to the last @, because an unescaped
// password may hold a slash. When a ? or # comes before that @, the password
// cannot be told apart from the query, so only the scheme is kept.
func redactUnparsed(raw string) string {
	scheme, rest, _ := strings.Cut(raw, "://")
	if i := strings.LastIndex(rest, "@"); i >= 0 {
		if strings.ContainsAny(rest[:i], "?#") {
			return scheme + "://" + redactedPassword
		}
		rest = rest[i+1:]
	}
	if i := strings.IndexAny(rest, "?#"); i >= 0 {
		rest = rest[:i]
	}
	return scheme + "://" + rest
}

// RedactURLError returns err with credentials removed when it is a *url.Error,
// the type that net/url and net/http return with the whole URL inside. The
// result is still a *url.Error with the same cause, so errors.Is, errors.As and
// Timeout keep working. Causes that quote the URL are replaced: the reason
// for a URL that held something to remove and either failed to parse or had no
// scheme (it can quote part of the credentials), a cause for a URL that
// hidesUserinfo (net/http takes the user name for the host and part of the
// password for the port; only a canceled or expired context is kept), and a
// cause that itself names a URL with credentials. Call it on the error as
// net/url or net/http returned it; a wrapped one is not changed.
func RedactURLError(err error) error {
	urlErr, ok := err.(*url.Error)
	if !ok {
		return err
	}
	clean := RedactURL(urlErr.URL)
	cause := urlErr.Err
	// With the scheme left out, "user:password@host" parses with the user name
	// as its scheme, and net/http quotes that scheme.
	quotesURL := urlErr.Op == "parse" || !strings.Contains(urlErr.URL, "://")
	switch {
	case quotesURL && clean != urlErr.URL:
		cause = errors.New("invalid URL")
	case hidesUserinfo(urlErr.URL):
		cause = contextCauseOnly(cause)
	case cause != nil:
		if text := Redact(cause.Error()); text != cause.Error() {
			cause = errors.New(text)
		}
	}
	return &url.Error{Op: urlErr.Op, URL: clean, Err: cause}
}

// contextCauseOnly returns the context error that cause matches, which quotes
// nothing, or a fixed error in its place.
func contextCauseOnly(cause error) error {
	for _, known := range []error{context.Canceled, context.DeadlineExceeded} {
		if errors.Is(cause, known) {
			return known
		}
	}
	return errors.New("invalid URL")
}

// RedactReply is Redact for text that a server sent, such as a reply body or an
// error message in it. Such text can be cut off anywhere, so a URL that ends it
// before anything ends the host may be the start of a password: everything
// after its "://" is masked. A complete URL that ends the text with no path is
// masked too.
func RedactReply(text string) string {
	if m := openAuthorityPattern.FindStringSubmatchIndex(text); m != nil {
		text = text[:m[2]] + redactedPassword + text[m[3]:]
	}
	return Redact(text)
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
