package qsparser

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Order specifies the order of fields
// Directions must be either nil or equals in length to that of Fields
type SearchOrder struct {
	// Fields is the field name to order by
	Fields []string
	// Directions is the direction to order by
	// A value of nil indicates that all fields are sorted in the same direction
	Directions []string
}

// Field hold values for a single field
type SearchField struct {
	// SearchVal is the search value for the field
	SearchVal *string
	// SearchOp is the search operation for the field
	SearchOp *string
	// Visible specifies whether the field is visible to the consumer or not
	Visible *bool
}

// SearchParams holds the search parameters
// nil values indicate the parameter was not specified in the query string
type SearchParams struct {
	// Draw can be used to ensure order of requests and responses for asynchronous requests
	Draw *int

	// Page is used for pagination and represents the current page number
	Page *int

	// PerPage is used for pagination and represents the number of rows to return for each page
	PerPage *int

	// GlobalSearchVal represents a global search value that can be used to search across multiple fields
	GlobalSearchVal *string

	// Order represents the field order
	Order *SearchOrder

	// Fields represents the fields and their options
	Fields map[string]*SearchField
}

func (sp *SearchParams) SetRawQuery(req *http.Request) {
	q := req.URL.Query()
	if sp.Draw != nil {
		q.Add("draw", fmt.Sprintf("%d", *sp.Draw))
	}
	if sp.PerPage != nil {
		q.Add("per_page", fmt.Sprintf("%d", *sp.PerPage))
	}
	if sp.Page != nil {
		q.Add("page", fmt.Sprintf("%d", *sp.Page))
	}
	if sp.GlobalSearchVal != nil {
		q.Add("search", *sp.GlobalSearchVal)
	}
	if sp.Order != nil {
		q.Add("ord[fields]", strings.Join(sp.Order.Fields, ","))
		q.Add("ord[dir]", strings.Join(sp.Order.Directions, ","))
	}
	for k, v := range sp.Fields {
		if v == nil {
			continue
		}
		if v.Visible != nil {
			q.Add(fmt.Sprintf("f[%s][visible]", k), strconv.FormatBool(*v.Visible))
		}
		if v.SearchOp != nil {
			q.Add(fmt.Sprintf("f[%s][search][op]", k), *v.SearchOp)
		}
		if v.SearchVal != nil {
			q.Add(fmt.Sprintf("f[%s][search][value]", k), *v.SearchVal)
		}
	}
	req.URL.RawQuery = q.Encode()
}

// ParseSearchQuery parses a HTTP request's query string into a SeachParams
// It supports global search, field search and pagination.
//
// The supported query string parameters are:
//
// * draw - an optional integer value that can be used to ensure the order of quests and responses of asynchronous requests
// * page - an optional integer value that can be used as the requested page in pagination
// * per_page - an optional integer value that can be used as the requested number of rows returned in pagination
// * search - an optional string value that can be used as a global search value
// * ord - a complex value (see below) that specifies the order of the requested fields for field search
// * f - a complex value (see below) that specifies the field names, search values and search operator of the requested fields for field search
//
// The package doesn't assume anything about the validity or possible values of the parameter values. Is is the job of the consumer to validate those.
//
// The 'ord' parameter format is any one of:
//
// `ord[fields]=value1,value2,...`
//
// Where `value1,value2,...` is a comma-separated list of strings representing the field names to order by.
//
// `ord[dir]=value1,value2,...`
//
// Where `value1,value2,...` is a comma-separated list of strings representing the direction in which the fields should be ordered.
//
// If `ord[dir]` is specified the number of values must be be equal to the number of values in ord[fields].
//
// None of the values may contain commas.
//
// Example:
//
// Order by name ascending and age descending:
//
// `ord[fields]=name,age&ord[dir]=asc,desc`
//
// The 'f' parameter format is any of:
//
// `f[field_name][search][value]=value`
//
// Where `field_name` is a string representing a field name and `value` a string representing the search value for that field.
//
// `f[field_name][search][op]=value`
//
// Where `field_name` is a string representing a field name and `value` a string representing the search operator for that field.
//
// `f[field_name][visible]=value`
//
// Where `field_name` is a string representing a field name and `value` is a boolean (true or false).
//
// This parameter is used to inform about the visibility of a field, e.g. if it should be hidden from view.
//
// Examples:
//
// 1) Show records where 'name' eq(==) 'Peter' indicating that the field should be made visible:
//
// `f[name][search][value]=Peter&f[name][search][op]=eq&f[name][search][visible]=true`
//
// 2) Show records where 'age' ge(>=) '20' indicating that the field should not be made visible:
//
// `f[age][search][value]=20&f[age][search][op]=ge&f[age][search][visible]=false`
func ParseSearchQuery(r *http.Request) (sp *SearchParams, err error) {
	fSearchFieldRegex := regexp.MustCompile(`^f\[([^]]+)\]\[([^]]+)\](?:\[([^]]+)\])?$`)

	sp = &SearchParams{}

	if err = r.ParseForm(); err != nil {
		return
	}

	for key, val := range r.Form {
		var (
			searchFieldKey    string
			searchFieldName   string
			searchFieldParam  string
			searchFieldParam2 string
			field             *SearchField
		)

		val0 := val[0]

		if fSearchField := fSearchFieldRegex.FindStringSubmatch(key); fSearchField != nil {
			searchFieldKey = fSearchField[0]
			searchFieldName = fSearchField[1]
			searchFieldParam = fSearchField[2]
			searchFieldParam2 = fSearchField[3]
			if sp.Fields == nil {
				sp.Fields = make(map[string]*SearchField)
			}
			var found bool
			field, found = sp.Fields[searchFieldName]
			if !found {
				sp.Fields[searchFieldName] = &SearchField{}
				field = sp.Fields[searchFieldName]
			}
		}

		var e error
		switch key {
		case "draw":
			sp.Draw, e = parseStrIntPtr(val0)
			err = errors.Join(err, e)
		case "page":
			sp.Page, e = parseStrIntPtr(val0)
			err = errors.Join(err, e)
		case "per_page":
			sp.PerPage, e = parseStrIntPtr(val0)
			err = errors.Join(err, e)
		case "search":
			globalSearchVal := escape(strings.Trim(val0, " "))
			sp.GlobalSearchVal = &globalSearchVal
		case "ord[fields]":
			if val0 == "" {
				err = errors.Join(err, fmt.Errorf("parsing ord parameter: %v", key))
				continue
			}
			fields := strings.Split(val0, ",")
			var lenAllFields int
			for i := range fields {
				fields[i] = strings.TrimSpace(fields[i])
				lenAllFields += len(fields[i])
			}
			if lenAllFields == 0 {
				err = errors.Join(err, fmt.Errorf("parsing ord parameter: %v", key))
				continue
			}
			if sp.Order == nil {
				sp.Order = &SearchOrder{}
			}
			sp.Order.Fields = fields
		case "ord[dir]":
			if val0 == "" {
				err = errors.Join(err, fmt.Errorf("parsing ord parameter: %v", key))
				continue
			}
			dirs := strings.Split(val0, ",")
			for i := range dirs {
				dirs[i] = strings.TrimSpace(dirs[i])
			}
			if sp.Order == nil {
				sp.Order = &SearchOrder{}
			}
			sp.Order.Directions = dirs
		case searchFieldKey:
			switch searchFieldParam {
			case "search":
				switch searchFieldParam2 {
				case "value":
					v := escape(val0)
					field.SearchVal = &v
				case "op":
					if val0 != "" {
						op := escape(val0)
						field.SearchOp = &op
					}
				}
			case "visible":
				if val0 != "" {
					visible := val0 == "true"
					field.Visible = &visible
				}
			}
		}
	}
	if sp.Order != nil && len(sp.Order.Directions) > 0 && len(sp.Order.Fields) != len(sp.Order.Directions) {
		err = errors.Join(err, errors.New("parsing ord parameter: more directions than fields"))
		sp.Order.Directions = nil
	}
	return
}

func parseStrIntPtr(s string) (*int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return &val, err
}

func escape(s string) (e string) {
	e = strings.Replace(s, "\\", "\\\\", -1)
	e = strings.Replace(e, "'", "\\'", -1)
	return
}
