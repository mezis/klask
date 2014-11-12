package query

import (
	"github.com/juju/errgo"
	// "github.com/mezis/klask/index"
	"strings"
)

// A generic query, which can combine $and, $or, field filters, and a $by
// clause. They will be run in an unspecified order, except the optional $by clause
// which is run last.
// Represented by a JSON object.
type query_generic_t struct {
	// queries []Query
	queries *query_and_t
	order   *query_order_t
}

func (self *query_generic_t) parse(parsed interface{}) error {
	var (
		err     error                  = nil
		order   *query_order_t         = nil
		limit   uint                   = 0
		offset  uint                   = 0
		queries []Query                = make([]Query, 0)
		node    map[string]interface{} = nil
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
			err = q.parse(subnode)
		case key == "$and":
			q := new(query_and_t)
			queries = append(queries, q)
			err = q.parse(subnode)
		case key == "$by":
			order := new(query_order_t)
			err = order.parse(subnode)
		case key == "$limit":
			limit, err = self.parseInt(subnode)
		case key == "$offset":
			offset, err = self.parseInt(subnode)
		case strings.HasPrefix(key, "$"):
			err = errgo.Newf("unknown subquery type '%s'", key)
		default:
			q := new(query_field_t)
			queries = append(queries, q)
			err = q.parse(key, subnode)
		}
		if err != nil {
			return errgo.Mask(err)
		}
	}

	// create the and query
	if len(queries) > 0 {
		q := new(query_and_t)
		q.queries = queries
		self.queries = q
	}

	// the order query, if any, should be last
	if order == nil && (limit != 0 || offset != 0) {
		return errgo.New("cannot have $limit or $offset without $by")
	}
	if order != nil {
		order.limit = limit
		order.offset = offset
		self.order = order
	}

	return nil
}

func (self *query_generic_t) parseInt(val interface{}) (uint, error) {
	switch v := val.(type) {
	case float64:
		if v < 0 {
			return 0, errgo.Newf("unexpected negative value '%v'", v)
		}
		return uint(v), nil
	default:
		return 0, errgo.Newf("expected positive int, got %T '%v'", v, v)
	}
}

func (self *query_generic_t) Run(records string, ctx Context) (string, error) {
	var (
		err      error  = nil
		unsorted string = ""
		sorted   string = ""
	)

	if self.queries != nil {
		unsorted, err = self.queries.Run(records, ctx)
		if err != nil {
			return "", errgo.Mask(err)
		}
	} else {
		unsorted = records
	}

	if self.order != nil {
		sorted, err = self.order.Run(unsorted, ctx)
		if err != nil {
			return "", errgo.Mask(err)
		}
	} else {
		sorted = unsorted
	}

	return sorted, nil
}
