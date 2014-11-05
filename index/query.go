package index

import (
	"encoding/json"
	// "github.com/garyburd/redigo/redis"
	"github.com/juju/errgo"
	"regexp"
)

type Query interface {
	// Results() []Id
	// Run(sourceKey string, targetKey string)
}

// the toplevel type of query, wraps a single query
type query_t struct {
	idx   Index
	query Query
}

// A generic query, which can combine $and, $or, field filters, and a $by
// clause. They will be run in an unspecified order, expect any $order clause
// which is run last.
// Represented by a JSON object.
type query_generic_t struct {
	queries []Query
}

// A query that orders results.
// Represented a field name, preceded by "+" (default, for
// "ascending) or "-" (for "descending").
type query_order_t struct {
	field     Field
	ascending bool
}

// A query that merges results.
// Represented by a JSON array of generic queries, as the value of an $or key.
type query_or_t struct {
	queries []Query
}

// A query that intersects results, progressively transforming the set of IDs.
// Represented by a JSON array of generic queries, as the value of an $and key.
type query_and_t struct {
	queries []Query
}

// A generic field filter, will point to one of the filtering queries
// (query_filter_*_t).
// Represented by a key/value pair in a generic query (query_generic_t).
type query_field_t struct {
	field Field
	query Query
}

// A comparison filter (less than, greater than, or both)
type query_filter_between_t struct {
	field        Field
	less_than    interface{}
	greater_than interface{}
}

// A selection filter (returns only values in the list)
type query_filter_in_t struct {
	field  Field
	values []interface{}
}

// An exclusion filter (returns only values not in the list)
type query_filter_not_in_t struct {
	field  Field
	values []interface{}
}

func NewQuery(idx Index) Query {
	q := new(query_t)
	q.idx = idx
	return q
}

func (self *query_t) UnmarshalJSON(data []byte) error {
	var parsed interface{}

	if self.idx == nil {
		return errgo.New("need an index to unmarshal a query")
	}

	// parse the syntax tree
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return errgo.Mask(err)
	}

	// start parsing!
	q := new(query_generic_t)
	err = q.parse(self.idx, parsed)
	if err != nil {
		return errgo.Mask(err)
	}

	self.query = q
	return nil
}

func (self *query_generic_t) parse(idx Index, parsed interface{}) error {
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
		switch key {
		case "$or":
			q := new(query_or_t)
			queries = append(queries, q)
			err = q.parse(idx, subnode)
		case "$and":
			q := new(query_and_t)
			queries = append(queries, q)
			err = q.parse(idx, subnode)
		case "$order":
			q := new(query_order_t)
			order = q
			err = q.parse(idx, subnode)
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

var gOrderRE = regexp.MustCompile("^([+-]?)(.*)$")

func (self *query_order_t) parse(idx Index, parsed interface{}) error {
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

func (self *query_or_t) parse(idx Index, parsed interface{}) error {
	return errgo.New("not implemented")
}

func (self *query_and_t) parse(idx Index, parsed interface{}) error {
	return errgo.New("not implemented")
}

func (self *query_field_t) parse(idx Index, name string, parsed interface{}) error {
	var (
		err error = nil
	)

	field, err := idx.Field(name)
	if err != nil {
		return errgo.Mask(err)
	}

	// type coercion
	switch parsed.(type) {
	case map[string]interface{}:
	default:
		return errgo.Newf("bad node type '%T' expected object (%v)", parsed, parsed)
	}
	node := parsed.(map[string]interface{})

	// figure out the filter type
	for key, val := range node {
		switch key {
		case "$gt", "$lt":
			q := new(query_filter_between_t)
			err = q.parse(field, val)
			self.query = q
		case "$in", "$eq":
			q := new(query_filter_in_t)
			err = q.parse(field, val)
			self.query = q
		case "$ni", "$neq":
			q := new(query_filter_not_in_t)
			err = q.parse(field, val)
			self.query = q
		default:
			return errgo.Newf("unknown filter type '%s'", key)
		}
		break
	}

	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (self *query_filter_between_t) parse(field Field, parsed interface{}) error {
	return errgo.New("not implemented")
}

func (self *query_filter_in_t) parse(field Field, parsed interface{}) error {
	return errgo.New("not implemented")
}

func (self *query_filter_not_in_t) parse(field Field, parsed interface{}) error {
	return errgo.New("not implemented")
}

// func (self *query_t) cleanKey(key string) error {
// 	conn := self.idx.Conn()
// 	defer conn.Close()

// 	_, err := conn.Do("DEL", key)
// 	if err != nil {
// 		return errgo.Mask(err)
// 	}
// 	return nil
// }

// func (self *query_t) Run(offset int, ttl int) ([]Id, error) {
// 	// FIXME: generate random temp keys
// 	// FIXME: as a second step, general hashed result keys for caching
// 	resultKey := "temp:1"
// 	defer self.cleanKey(resultKey) // ignoring errors

// 	sourceKey := self.idx.RecordsKey()

// 	for _, fi := range self.Filters {
// 		err := fi.run(sourceKey, resultKey)
// 		if err != nil {
// 			return nil, errgo.Mask(err)
// 		}
// 		sourceKey = resultKey
// 	}

// 	conn := self.idx.Conn()
// 	defer conn.Close()

// 	reply, err := redis.Values(conn.Do("SMEMBERS", resultKey))
// 	if err != nil {
// 		return nil, errgo.Mask(err)
// 	}

// 	ids := make([]Id, len(reply))
// 	_, err = redis.Scan(reply, ids)
// 	if err != nil {
// 		return nil, errgo.Mask(err)
// 	}

// 	return ids, nil
// }
