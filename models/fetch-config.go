package models

import (
	"net/url"
	"strconv"
)

type FetchConfig struct {
	Limit  int
	Cursor int
}

func FetchConfigFromQuery(u url.Values) (f *FetchConfig) {

	f = &FetchConfig{
		Limit:  10,
		Cursor: 0,
	}

	cursor := u.Get("_cursor")
	if cursor != "" {
		f.Cursor, _ = strconv.Atoi(cursor)
	}

	limit := u.Get("_limit")
	if limit != "" {
		f.Limit, _ = strconv.Atoi(limit)
	}

	return f

}
