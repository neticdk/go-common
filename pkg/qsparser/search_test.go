package qsparser

import (
	"fmt"
	"math"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRequest(qs string) (*SearchParams, error) {
	tgt := fmt.Sprintf("http://localhost?%s", qs)
	r := httptest.NewRequest("GET", tgt, nil)
	return ParseSearchQuery(r)
}

func ptr[T any](v T) *T {
	return &v
}

func Test_ParseSearchQuery_Params(t *testing.T) {
	// t.Parallel()

	t.Run("parameters", func(t *testing.T) {
		tests := map[string]struct {
			in     string
			out    *SearchParams
			expErr bool
		}{
			"empty query string must return empty struct": {
				in:     "",
				out:    &SearchParams{},
				expErr: false,
			},
			"draw must accept a MinInt value": {
				in:     fmt.Sprintf("draw=%d", math.MinInt),
				out:    &SearchParams{Draw: ptr(math.MinInt)},
				expErr: false,
			},
			"draw must accept a value of 1": {
				in:     "draw=1",
				out:    &SearchParams{Draw: ptr(1)},
				expErr: false,
			},
			"draw must accept a MaxInt value": {
				in:     fmt.Sprintf("draw=%d", math.MaxInt),
				out:    &SearchParams{Draw: ptr(math.MaxInt)},
				expErr: false,
			},
			"draw must not accept a letter": {
				in:     "draw=a",
				out:    &SearchParams{},
				expErr: true,
			},
			"page must accept a MinInt value": {
				in:     fmt.Sprintf("page=%d", math.MinInt),
				out:    &SearchParams{Page: ptr(math.MinInt)},
				expErr: false,
			},
			"page must accept a value of 1": {
				in:     "page=1",
				out:    &SearchParams{Page: ptr(1)},
				expErr: false,
			},
			"page must accept a MaxInt value": {
				in:     fmt.Sprintf("page=%d", math.MaxInt),
				out:    &SearchParams{Page: ptr(math.MaxInt)},
				expErr: false,
			},
			"page must not accept a letter": {
				in:     "page=a",
				out:    &SearchParams{},
				expErr: true,
			},
			"perPage must accept a MinInt value": {
				in:     fmt.Sprintf("per_page=%d", math.MinInt),
				out:    &SearchParams{PerPage: ptr(math.MinInt)},
				expErr: false,
			},
			"perPage must accept a value of 1": {
				in:     "per_page=1",
				out:    &SearchParams{PerPage: ptr(1)},
				expErr: false,
			},
			"perPage must accept a MaxInt value": {
				in:     fmt.Sprintf("per_page=%d", math.MaxInt),
				out:    &SearchParams{PerPage: ptr(math.MaxInt)},
				expErr: false,
			},
			"perPage must not accept a letter": {
				in:     "per_page=a",
				out:    &SearchParams{},
				expErr: true,
			},
			"search must accept an empty value": {
				in:     "search=",
				out:    &SearchParams{GlobalSearchVal: ptr("")},
				expErr: false,
			},
			"search must accept a single alphanumeric word": {
				in:     "search=simple123",
				out:    &SearchParams{GlobalSearchVal: ptr("simple123")},
				expErr: false,
			},
			"search must accept urlencoded spaces": {
				in:     "search=simple%20with%20spaces",
				out:    &SearchParams{GlobalSearchVal: ptr("simple with spaces")},
				expErr: false,
			},
			"search must accept and escape backspace": {
				in:     "search=simple%20with%20\\backspace",
				out:    &SearchParams{GlobalSearchVal: ptr(`simple with \\backspace`)},
				expErr: false,
			},
			"search must accept and escape single quote": {
				in:     "search=simple%20with%20'quote'",
				out:    &SearchParams{GlobalSearchVal: ptr(`simple with \'quote\'`)},
				expErr: false,
			},
			"ord must include a list of fields": {
				in:     "ord[fields]=",
				out:    &SearchParams{},
				expErr: true,
			},
			"ord must accept a single field name": {
				in: "ord[fields]=name",
				out: &SearchParams{Order: &SearchOrder{
					Fields:     []string{"name"},
					Directions: nil,
				}},
				expErr: false,
			},
			"ord must accept multiple field names": {
				in: "ord[fields]=name,age",
				out: &SearchParams{Order: &SearchOrder{
					Fields:     []string{"name", "age"},
					Directions: nil,
				}},
				expErr: false,
			},
			"ord must accept multiple field names together with the same amount of order directions": {
				in: "ord[fields]=name,age&ord[dir]=asc,desc",
				out: &SearchParams{Order: &SearchOrder{
					Fields:     []string{"name", "age"},
					Directions: []string{"asc", "desc"},
				}},
				expErr: false,
			},
			"ord must not accept multiple field names without the same amount of order directions": {
				in: "ord[fields]=name,age&ord[dir]=asc",
				out: &SearchParams{Order: &SearchOrder{
					Fields:     []string{"name", "age"},
					Directions: nil,
				}},
				expErr: true,
			},
			"ord must not accept directions if no fields name are given": {
				in: "ord[dir]=asc",
				out: &SearchParams{Order: &SearchOrder{
					Fields:     nil,
					Directions: nil,
				}},
				expErr: true,
			},
			"ord must not accept an empty list of field names": {
				in:     "ord[fields]=,",
				out:    &SearchParams{},
				expErr: true,
			},
			"f must accept a visible value of true": {
				in:     "f[name][visible]=true",
				out:    &SearchParams{Fields: map[string]*SearchField{"name": {Visible: ptr(true)}}},
				expErr: false,
			},
			"f must accept a visible value of false": {
				in:     "f[name][visible]=false",
				out:    &SearchParams{Fields: map[string]*SearchField{"name": {Visible: ptr(false)}}},
				expErr: false,
			},
			"empty visible value for f must result in nil Visible value": {
				in:     "f[name][visible]=",
				out:    &SearchParams{Fields: map[string]*SearchField{"name": {Visible: nil}}},
				expErr: false,
			},
			"f must accept an empty search value": {
				in:     "f[name][search][value]=",
				out:    &SearchParams{Fields: map[string]*SearchField{"name": {SearchVal: ptr("")}}},
				expErr: false,
			},
			"f must accept urlencoded spaces": {
				in:     "f[name][search][value]=simple%20with%20space",
				out:    &SearchParams{Fields: map[string]*SearchField{"name": {SearchVal: ptr("simple with space")}}},
				expErr: false,
			},
			"f must accept and escape backspace": {
				in:     "f[name][search][value]=simple%20with%20\\backspace",
				out:    &SearchParams{Fields: map[string]*SearchField{"name": {SearchVal: ptr(`simple with \\backspace`)}}},
				expErr: false,
			},
			"f must accept and escape single quote": {
				in:     "f[name][search][value]=simple%20with%20'quote'",
				out:    &SearchParams{Fields: map[string]*SearchField{"name": {SearchVal: ptr(`simple with \'quote\'`)}}},
				expErr: false,
			},
			"empty op subparameter value for the search parameter must result in SearchOp: nil and SearchVal: nil": {
				in:     "f[name][search][op]=",
				out:    &SearchParams{Fields: map[string]*SearchField{"name": {SearchOp: nil, SearchVal: nil}}},
				expErr: false,
			},
			"f search op must accept a single word value": {
				in:     "f[name][search][op]=eq",
				out:    &SearchParams{Fields: map[string]*SearchField{"name": {SearchOp: ptr("eq")}}},
				expErr: false,
			},
			"a full query string with all options must return a full struct with all options": {
				in: "draw=1&page=1&per_page=10&search=a%20search&ord[fields]=name,age&ord[dir]=desc,asc&f[name][visible]=true&f[name][search][value]=peter&f[name][search][op]=ilike&f[age][visible]=false&f[age][search][value]=42&f[age][search][op]=ge",
				out: &SearchParams{
					Draw:            ptr(1),
					Page:            ptr(1),
					PerPage:         ptr(10),
					GlobalSearchVal: ptr("a search"),
					Order: &SearchOrder{
						Fields:     []string{"name", "age"},
						Directions: []string{"desc", "asc"},
					},
					Fields: map[string]*SearchField{
						"name": {Visible: ptr(true), SearchVal: ptr("peter"), SearchOp: ptr("ilike")},
						"age":  {Visible: ptr(false), SearchVal: ptr("42"), SearchOp: ptr("ge")},
					},
				},
				expErr: false,
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				want, err := testRequest(tt.in)
				if !tt.expErr {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
				assert.Equal(t, tt.out, want)
			})
		}
	})
}

func Test_SearchParams_SetRawQuery(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		in  *SearchParams
		out map[string]string
	}{
		"empty query string must return empty struct": {
			in:  &SearchParams{},
			out: map[string]string{},
		},
		"draw must accept a MinInt value": {
			in:  &SearchParams{Draw: ptr(math.MinInt)},
			out: map[string]string{"draw": fmt.Sprintf("%d", math.MinInt)},
		},
		"draw must accept a value of 1": {
			in:  &SearchParams{Draw: ptr(1)},
			out: map[string]string{"draw": "1"},
		},
		"draw must accept a MaxInt value": {
			in:  &SearchParams{Draw: ptr(math.MaxInt)},
			out: map[string]string{"draw": fmt.Sprintf("%d", math.MaxInt)},
		},
		"page must accept a MinInt value": {
			in:  &SearchParams{Page: ptr(math.MinInt)},
			out: map[string]string{"page": fmt.Sprintf("%d", math.MinInt)},
		},
		"page must accept a value of 1": {
			in:  &SearchParams{Page: ptr(1)},
			out: map[string]string{"page": "1"},
		},
		"page must accept a MaxInt value": {
			in:  &SearchParams{Page: ptr(math.MaxInt)},
			out: map[string]string{"page": fmt.Sprintf("%d", math.MaxInt)},
		},
		"perPage must accept a MinInt value": {
			in:  &SearchParams{PerPage: ptr(math.MinInt)},
			out: map[string]string{"per_page": fmt.Sprintf("%d", math.MinInt)},
		},
		"perPage must accept a value of 1": {
			in:  &SearchParams{PerPage: ptr(1)},
			out: map[string]string{"per_page": "1"},
		},
		"perPage must accept a MaxInt value": {
			in:  &SearchParams{PerPage: ptr(math.MaxInt)},
			out: map[string]string{"per_page": fmt.Sprintf("%d", math.MaxInt)},
		},
		"search must accept an empty value": {
			in:  &SearchParams{GlobalSearchVal: ptr("")},
			out: map[string]string{"search": ""},
		},
		"search must accept a single alphanumeric word": {
			in:  &SearchParams{GlobalSearchVal: ptr("simple123")},
			out: map[string]string{"search": "simple123"},
		},
		"search must accept urlencoded spaces": {
			in:  &SearchParams{GlobalSearchVal: ptr("simple with spaces")},
			out: map[string]string{"search": "simple with spaces"},
		},
		"search must accept and escape backspace": {
			in:  &SearchParams{GlobalSearchVal: ptr(`simple with \\backspace`)},
			out: map[string]string{"search": "simple with \\\\backspace"},
		},
		"search must accept and escape single quote": {
			in:  &SearchParams{GlobalSearchVal: ptr(`simple with \'quote\'`)},
			out: map[string]string{"search": "simple with \\'quote\\'"},
		},
		"ord must accept a single field name": {
			in: &SearchParams{Order: &SearchOrder{
				Fields:     []string{"name"},
				Directions: nil,
			}},
			out: map[string]string{"ord[fields]": "name"},
		},
		"ord must accept multiple field names": {
			in: &SearchParams{Order: &SearchOrder{
				Fields:     []string{"name", "age"},
				Directions: nil,
			}},
			out: map[string]string{"ord[fields]": "name,age"},
		},
		"ord must accept multiple field names together with the same amount of order directions": {
			in: &SearchParams{Order: &SearchOrder{
				Fields:     []string{"name", "age"},
				Directions: []string{"asc", "desc"},
			}},
			out: map[string]string{"ord[fields]": "name,age", "ord[dir]": "asc,desc"},
		},
		"ord must not accept multiple field names without the same amount of order directions": {
			in: &SearchParams{Order: &SearchOrder{
				Fields:     []string{"name", "age"},
				Directions: nil,
			}},
			out: map[string]string{},
		},
		"ord must not accept directions if no fields name are given": {
			in: &SearchParams{Order: &SearchOrder{
				Fields:     nil,
				Directions: nil,
			}},
			out: map[string]string{},
		},
		"f must accept a visible value of true": {
			in:  &SearchParams{Fields: map[string]*SearchField{"name": {Visible: ptr(true)}}},
			out: map[string]string{"f[name][visible]": "true"},
		},
		"f must accept a visible value of false": {
			in:  &SearchParams{Fields: map[string]*SearchField{"name": {Visible: ptr(false)}}},
			out: map[string]string{"f[name][visible]": "false"},
		},
		"empty visible value for f must result in nil Visible value": {
			in:  &SearchParams{Fields: map[string]*SearchField{"name": {Visible: nil}}},
			out: map[string]string{"f[name][visible]": ""},
		},
		"f must accept an empty search value": {
			in:  &SearchParams{Fields: map[string]*SearchField{"name": {SearchVal: ptr("")}}},
			out: map[string]string{"f[name][search][value]": ""},
		},
		"f must accept urlencoded spaces": {
			in:  &SearchParams{Fields: map[string]*SearchField{"name": {SearchVal: ptr("simple with space")}}},
			out: map[string]string{"f[name][search][value]": "simple with space"},
		},
		"f must accept and escape backspace": {
			in:  &SearchParams{Fields: map[string]*SearchField{"name": {SearchVal: ptr(`simple with \\backspace`)}}},
			out: map[string]string{"f[name][search][value]": "simple with \\\\backspace"},
		},
		"f must accept and escape single quote": {
			in:  &SearchParams{Fields: map[string]*SearchField{"name": {SearchVal: ptr(`simple with \'quote\'`)}}},
			out: map[string]string{"f[name][search][value]": "simple with \\'quote\\'"},
		},
		"empty op subparameter value for the search parameter must result in SearchOp: nil and SearchVal: nil": {
			in:  &SearchParams{Fields: map[string]*SearchField{"name": {SearchOp: nil, SearchVal: nil}}},
			out: map[string]string{},
		},
		"f search op must accept a single word value": {
			in:  &SearchParams{Fields: map[string]*SearchField{"name": {SearchOp: ptr("eq")}}},
			out: map[string]string{"f[name][search][op]": "eq"},
		},
		"a full query string with all options must return a full struct with all options": {
			in: &SearchParams{
				Draw:            ptr(1),
				Page:            ptr(1),
				PerPage:         ptr(10),
				GlobalSearchVal: ptr("a search"),
				Order: &SearchOrder{
					Fields:     []string{"name", "age"},
					Directions: []string{"desc", "asc"},
				},
				Fields: map[string]*SearchField{
					"name": {Visible: ptr(true), SearchVal: ptr("peter"), SearchOp: ptr("ilike")},
					"age":  {Visible: ptr(false), SearchVal: ptr("42"), SearchOp: ptr("ge")},
				},
			},
			out: map[string]string{
				"draw":                   "1",
				"page":                   "1",
				"per_page":               "10",
				"search":                 "a search",
				"ord[fields]":            "name,age",
				"ord[dir]":               "desc,asc",
				"f[name][visible]":       "true",
				"f[name][search][value]": "peter",
				"f[name][search][op]":    "ilike",
				"f[age][visible]":        "false",
				"f[age][search][value]":  "42",
				"f[age][search][op]":     "ge",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://localhost", nil)
			tt.in.SetRawQuery(req)
			for k, v := range tt.out {
				got := req.URL.Query().Get(k)
				want := v
				assert.Equal(t, want, got)
			}
		})
	}
}
