package models

import (
	"net/url"
	"strconv"
)

type FetchConfig struct {
	Limit  int
	Cursor int
}

func FetchConfigFromQuery(u *url.Values) (f *FetchConfig, err error) {

	f = &FetchConfig{
		Limit:  10,
		Cursor: 0,
	}

	cursor := u.Get("_cursor")
	if cursor != "" {
		f.Cursor, err = strconv.Atoi(cursor)
	}

	limit := u.Get("_limit")
	if limit != "" {
		f.Limit = strconv.Atoi(limit)
	}

	return f, err

}
