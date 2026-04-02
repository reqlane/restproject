package models

import (
	"net/url"
)

var TeacherFieldNames = []string{
	"first_name",
	"last_name",
	"email",
	"class",
	"subject",
}

var StudentFieldNames = []string{
	"first_name",
	"last_name",
	"email",
	"class",
}

var ExecFieldNames = []string{
	"first_name",
	"last_name",
	"email",
	"username",
	"inactive_status",
	"role",
}

type Criteria struct {
	Filters  map[string]string
	Sortings []string
}

func (c *Criteria) AddFilters(query url.Values, fieldNames []string) {
	for _, fieldName := range fieldNames {
		value := query.Get(fieldName)
		if value != "" {
			c.Filters[fieldName] = value
		}
	}
}

type Pagination struct {
	Page, Limit int
}
