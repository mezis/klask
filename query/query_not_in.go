package query

import (
	"github.com/juju/errgo"
	"github.com/mezis/klask/index"
)

// An exclusion filter (returns only values not in the list)
type query_filter_not_in_t struct {
	query_filter_membership_t
}

func (self *query_filter_not_in_t) Run(idx index.Index, targetKey string) error {
	return errgo.New("not implemented")
}
