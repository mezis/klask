package query

import (
	"github.com/juju/errgo"
	// "github.com/mezis/klask/index"
)

// A comparison filter (less than, greater than, or both)
type query_filter_between_t struct {
	name         string
	less_than    interface{}
	greater_than interface{}
}

func (self *query_filter_between_t) parse(name string, parsed map[string]interface{}) error {
	self.name = name

	// parse
	for key, val := range parsed {
		switch key {
		case "$gt":
			self.greater_than = val
		case "$lt":
			self.less_than = val
		default:
			return errgo.Newf("unexpected key '%s' for range filter in '%v'", key, parsed)
		}
	}

	// we don't check values, field names, or operator/operand compatibility at this point;
	// it will be done lazily when applying the filter

	return nil
}

func (self *query_filter_between_t) Run(records string, ctx Context) (string, error) {
	return "", errgo.New("not implemented")
}
