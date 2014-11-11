package query

import (
	"github.com/juju/errgo"
	// "github.com/mezis/klask/index"
	"regexp"
)

// A query that orders results.
// Represented a field name, preceded by "+" (default, for
// "ascending) or "-" (for "descending").
type query_order_t struct {
	name      string
	ascending bool
	limit     uint
	offset    uint
}

var gOrderRE = regexp.MustCompile("^([+-]?)(.*)$")

func (self *query_order_t) parse(parsed interface{}) error {
	switch val := parsed.(type) {
	case string:
		matches := gOrderRE.FindStringSubmatch(val)
		if matches == nil {
			return errgo.Newf("bad order '%s', expected to match /[+-].*/", val)
		}
		self.name = matches[2]
		self.ascending = (matches[1] != "-")
	default:
		return errgo.Newf("bad order '%v' (%T), expected a string", val, val)
	}
	return nil
}

func (self *query_order_t) Run(records string, ctx Context) (string, error) {
	return "", errgo.New("not implemented")
}
