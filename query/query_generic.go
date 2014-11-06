package query

import (
	"github.com/juju/errgo"
	"github.com/mezis/klask/index"
	"strings"
)

// A generic query, which can combine $and, $or, field filters, and a $by
// clause. They will be run in an unspecified order, except the optional $by clause
// which is run last.
// Represented by a JSON object.
type query_generic_t struct {
	queries []Query
}

func (self *query_generic_t) parse(idx index.Index, parsed interface{}) error {
	var (
		err     error = nil
		order   Query = nil
		queries       = make([]Query, 0)
		node    map[string]interface{}
	)

	switch n := parsed.(type) {
	case map[string]interface{}:
		node = n
	default:
		return errgo.Newf("unexpected node  of type %T (%v)", node, node)
	}

	for key, subnode := range node {
		switch {
		case key == "$or":
			q := new(query_or_t)
			queries = append(queries, q)
			err = q.parse(idx, subnode)
		case key == "$and":
			q := new(query_and_t)
			queries = append(queries, q)
			err = q.parse(idx, subnode)
		case key == "$by":
			q := new(query_order_t)
			order = q
			err = q.parse(idx, subnode)
		case strings.HasPrefix(key, "$"):
			err = errgo.Newf("unknown subquery type '%s'", key)
		default:
			q := new(query_field_t)
			queries = append(queries, q)
			err = q.parse(idx, key, subnode)
		}
		if err != nil {
			return errgo.Mask(err)
		}
	}

	// the order query, if any, should be last
	if order != nil {
		queries = append(queries, order)
	}
	self.queries = queries
	return nil
}
