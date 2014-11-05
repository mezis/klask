package query

import (
	"github.com/juju/errgo"
	"github.com/mezis/klask/index"
	"regexp"
)

// A query that orders results.
// Represented a field name, preceded by "+" (default, for
// "ascending) or "-" (for "descending").
type query_order_t struct {
	field     index.Field
	ascending bool
}

var gOrderRE = regexp.MustCompile("^([+-]?)(.*)$")

func (self *query_order_t) parse(idx index.Index, parsed interface{}) error {
	switch val := parsed.(type) {
	case string:
		matches := gOrderRE.FindStringSubmatch(val)
		if matches == nil {
			return errgo.Newf("bad order '%s', expected to match /[+-].*/", val)
		}
		field, err := idx.Field(matches[1])
		if err != nil {
			return errgo.Mask(err)
		}
		self.field = field
		self.ascending = (matches[0] != "-")
	default:
		return errgo.Newf("bad order '%v' (%T), expected a string", val, val)
	}
	return nil
}
