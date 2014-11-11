package query

import (
	"github.com/juju/errgo"
	// "github.com/mezis/klask/index"
)

// A selection filter (returns only values in the list)
type query_filter_in_t struct {
	query_filter_membership_t
}

func (self *query_filter_in_t) Run(records string, ctx Context) (string, error) {
	// err := self.field.Filter("in", self.values, idx.RecordsKey(), targetKey)
	// if err != nil {
	// return errgo.Mask(err)
	// }
	// return nil
	return "", errgo.New("not implemented")
}
