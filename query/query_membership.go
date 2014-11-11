package query

import (
	"github.com/juju/errgo"
	// "github.com/mezis/klask/index"
)

// Generic membership query, used to share parsing code between
// `query_filter_in_t` and `query_filter_not_in_t`.
type query_filter_membership_t struct {
	name   string
	values []interface{}
}

func (self *query_filter_membership_t) parse(name string, parsed map[string]interface{}) error {
	self.name = name

	for key, val := range parsed {
		switch key {
		case "$in", "$ni":
			switch v := val.(type) {
			case []interface{}:
				self.values = v
			default:
				return errgo.Newf("bad filter, %s requires an array argument (got '%v')", key, v)
			}
		case "$eq", "$neq":
			values := make([]interface{}, 1)
			values[0] = val
			self.values = values
		default:
			return errgo.Newf("bad filter '%s', "+
				"membership filters cannot be combined with others (in '%v')", key, parsed)
		}
	}
	return nil
}
