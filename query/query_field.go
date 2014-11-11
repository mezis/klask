package query

import (
	"github.com/juju/errgo"
	// "github.com/mezis/klask/index"
)

// A generic field filter, will point to one of the filtering queries
// (query_filter_*_t).
// Represented by a key/value pair in a generic query (query_generic_t).
type query_field_t struct {
	name  string
	query Query
}

func (self *query_field_t) parse(name string, parsed interface{}) error {
	var (
		err error = nil
	)

	// type coercion
	switch parsed.(type) {
	case map[string]interface{}:
	default:
		return errgo.Newf("bad node type '%T' expected object (%v)", parsed, parsed)
	}
	node := parsed.(map[string]interface{})

	// figure out the filter type
	for key, _ := range node {
		switch key {
		case "$gt", "$lt":
			q := new(query_filter_between_t)
			err = q.parse(name, node)
			self.query = q
		case "$in", "$eq":
			q := new(query_filter_in_t)
			err = q.parse(name, node)
			self.query = q
		case "$ni", "$neq":
			q := new(query_filter_not_in_t)
			err = q.parse(name, node)
			self.query = q
		default:
			return errgo.Newf("unknown filter type '%s'", key)
		}
		break
	}
	self.name = name

	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (self *query_field_t) Run(records string, ctx Context) (string, error) {
	return "", errgo.New("not implemented")
}
