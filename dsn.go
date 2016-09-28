package utils

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type DSN struct {
	Host     string
	User     string
	Password string
	Database string
	Timeout  time.Duration
}

func (dsn *DSN) String() string {
	u := url.URL{
		Host: fmt.Sprintf("tcp(%s)", dsn.Host),
		Path: "/" + dsn.Database,
		RawQuery: url.Values{
			"timeout": {dsn.Timeout.String()},
		}.Encode(),
	}
	if dsn.Password == "" {
		u.User = url.User(dsn.User)
	} else {
		u.User = url.UserPassword(dsn.User, dsn.Password)
	}
	return strings.TrimPrefix(u.String(), "//")
}
