package query

import (
	"github.com/juju/errgo"
	"github.com/mezis/klask/index"
)

// A comparison filter (less than, greater than, or both)
type query_filter_between_t struct {
	field        index.Field
	less_than    interface{}
	greater_than interface{}
}

func (self *query_filter_between_t) parse(field index.Field, parsed map[string]interface{}) error {
	self.field = field

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

	// we don't check values, or operator/operand compatibility at this point;
	// it will be done lazily when applying the filter

	return nil
}
