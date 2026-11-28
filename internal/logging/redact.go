package logging

import (
	"context"
	"errors"
	"net"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// redactedPassword replaces a masked password.
const redactedPassword = "xxxxx"

// A query or fragment parameter is secret, and dropped from logged URLs, when
// its name ends in one of secretQuerySuffixes or is one of secretQueryParams.
// Both are compared without regard to case. The suffixes cover names such as
// x-api-key, subscription-key, private_token, client_secret, X-Amz-Signature,
// X-Amz-Credential and sslpassword.
var (
	secretQuerySuffixes = []string{"key", "token", "secret", "signature", "password", "credential"}
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
// secret query and fragment parameters are dropped; and a key/value DSN has its
// password masked. A URL that does not parse, or whose userinfo may hold an
// unescaped /, ? or # (see hidesUserinfo), also loses its user name and whole
// query.
func RedactURL(raw string) string {
	if !strings.Contains(raw, "://") {
		if keywordPasswordPattern.MatchString(raw) || strings.ContainsAny(raw, " \t\r\n") || !strings.ContainsAny(raw, "@?#") {
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
	if query, dropped := dropSecretParams(u.RawQuery, "&;", isSecretQueryPair); dropped {
		changed = true
		u.RawQuery = query
	}
	// The fragment is read from raw, because u.EscapedFragment may re-encode it
	// and so move where its pairs begin and end. A ? separates the parameters in
	// a fragment such as "/callback?access_token=...".
	_, rawFragment, _ := strings.Cut(raw, "#")
	fragment, droppedFromFragment := dropSecretParams(rawFragment, "&;?", isSecretFragmentPair)
	if droppedFromFragment {
		changed = true
		u.Fragment, u.RawFragment = "", ""
	}
	if !changed {
		return raw
	}
	if droppedFromFragment && fragment != "" {
		return u.String() + "#" + fragment
	}
	return u.String()
}

// dropSecretParams returns params without the pairs that isSecret matches, and
// reports whether it dropped any. Any byte of separators ends a pair. A secret
// may hold an unescaped ; or ?, so what follows a dropped pair after one of
// those, up to the next name=value pair or &, is dropped with it.
func dropSecretParams(params, separators string, isSecret func(pair string) bool) (string, bool) {
	var kept strings.Builder
	var separator byte // the one before the current pair
	var started, dropped, inSecret, startsQuery bool
	for {
		end := strings.IndexAny(params, separators)
		pair := params
		if end >= 0 {
			pair = params[:end]
		}
		// A dropped pair leaves its ? to the next pair that is kept.
		startsQuery = startsQuery || separator == '?'
		switch {
		case isSecret(pair):
			dropped, inSecret = true, true
		case inSecret && (separator == ';' || separator == '?') && !strings.Contains(pair, "="):
		default:
			if startsQuery {
				separator = '?'
			}
			if separator != 0 && (started || separator == '?') {
				kept.WriteByte(separator)
			}
			kept.WriteString(pair)
			started, inSecret, startsQuery = true, false, false
		}
		if end < 0 {
			return kept.String(), dropped
		}
		separator, params = params[end], params[end+1:]
	}
}

func isSecretQueryPair(pair string) bool {
	name, _, _ := strings.Cut(pair, "=")
	if decoded, err := url.QueryUnescape(name); err == nil {
		name = decoded
	}
	return isSecretQueryParam(name)
}

// isSecretFragmentPair is isSecretQueryPair for a fragment, where a part
// without = is an anchor or a route, not a parameter.
func isSecretFragmentPair(pair string) bool {
	return strings.Contains(pair, "=") && isSecretQueryPair(strings.TrimPrefix(pair, "/"))
}

func isSecretQueryParam(name string) bool {
	name = strings.ToLower(name)
	return secretQueryParams[name] || slices.ContainsFunc(secretQuerySuffixes, func(suffix string) bool {
		return strings.HasSuffix(name, suffix)
	})
}

// hidesUserinfo reports whether an @ after the apparent end of raw's authority
// may be the real end of its userinfo. url.Parse accepts
// "http://user:12/34@host" as host "user", port "12" and path "/34@host". In a
// well-formed URL, userinfo ends before the first /, ? or #, so a later @ is
// part of the path, query or fragment when the authority is unambiguously a host
// (see isCompleteHost). Otherwise it is read that way only when the authority is
// plainly a host and no colon comes before the @, so no password can be there.
func hidesUserinfo(raw string) bool {
	_, rest, found := strings.Cut(raw, "://")
	end := strings.IndexAny(rest, "/?#")
	at := strings.LastIndex(rest, "@")
	if !found || end < 0 || at < end {
		return false
	}
	authority := rest[:end]
	if isCompleteHost(authority) {
		return false
	}
	return !isPlainHost(authority) || strings.Contains(rest[:at], ":")
}

// isCompleteHost reports whether authority, which would otherwise be read as
// user:password, can only be a host: a bracketed IP address with or without a
// port, or a plain host with a port from 0 to 65535. A single-label name such
// as "user:12" can be either, and so is not one. A user name with a dot and a
// password that starts with digits and a /, ? or # looks the same as a host and
// port, and is taken for one.
func isCompleteHost(authority string) bool {
	if address, found := strings.CutPrefix(authority, "["); found {
		if address, found = strings.CutSuffix(address, "]"); found {
			return net.ParseIP(address) != nil
		}
	}
	host, port, err := net.SplitHostPort(authority)
	if err != nil || strings.Contains(host, "@") {
		return false
	}
	if _, err := strconv.ParseUint(port, 10, 16); err != nil {
		return false
	}
	return isPlainHost(host) || net.ParseIP(host) != nil
}

// isPlainHost reports whether host is unlikely to be a user name: it has a dot,
// or is localhost.
func isPlainHost(host string) bool {
	return host == "localhost" || strings.Contains(host, ".")
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
