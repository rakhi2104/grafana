package datasource

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/grafana/grafana/pkg/infra/log"
)

var logger = log.New("datasource")

// URLValidationError represents an error from validating a data source URL.
type URLValidationError struct {
	error

	url string
}

// Error returns the error message.
func (e URLValidationError) Error() string {
	return fmt.Sprintf("Validation of data source URL %q failed: %s", e.url, e.error.Error())
}

// Unwrap returns the wrapped error.
func (e URLValidationError) Unwrap() error {
	return e.error
}

// reURL is a regexp to detect if a URL specifies the protocol. We match also strings where the actual protocol is
// missing (i.e., "://"), in order to catch these as invalid when parsing.
var reURL = regexp.MustCompile("^[^:]*://")

// ValidateURL validates a data source's URL.
//
// The data source's type and URL must be provided. If successful, the valid URL object is returned, otherwise an
// error is returned.
func ValidateURL(typeName, urlStr string) (*url.URL, error) {
	switch strings.ToLower(typeName) {
	case "mssql":
		return validateMSSQLURL(urlStr)
	default:
		logger.Debug("Applying default URL parsing for this data source type", "type", typeName, "url", urlStr)
	}

	// Make sure the URL starts with a protocol specifier, so parsing is unambiguous
	if !reURL.MatchString(urlStr) {
		logger.Debug(
			"Data source URL doesn't specify protocol, so prepending it with http:// in order to make it unambiguous",
			"type", typeName, "url", urlStr)
		urlStr = fmt.Sprintf("http://%s", urlStr)
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, URLValidationError{error: err, url: urlStr}
	}

	return u, nil
}

func validateMSSQLURL(u string) (*url.URL, error) {
	logger.Debug("Validating MSSQL URL", "url", u)

	// Recognize ODBC connection strings like host\instance:1234
	reODBC := regexp.MustCompile(`^[^\\]+(:?\\[^:]+)?(:?:\d+)?$`)
	var host string
	switch {
	case reODBC.MatchString(u):
		logger.Debug("Recognized as ODBC URL format", "url", u)
		host = u
	default:
		logger.Debug("Couldn't recognize as valid MSSQL URL", "url", u)
		return nil, fmt.Errorf("unrecognized MSSQL URL format: %q", u)
	}
	return &url.URL{
		Scheme: "sqlserver",
		Host:   host,
	}, nil
}
