package query

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errgo"
	"github.com/mezis/klask/index"
)

type Query interface {
	// `records` is a ZSET key, containing subset of records IDs.
	// The key returned should contain a subset of `sourceKey`.
	// Pass `nil` as a context to the top-level query.
	Run(records string, context Context) (string, error)
}

// the toplevel type of query, wraps a single query
type query_t struct {
	query Query
}

func New(idx index.Index) Query {
	q := new(query_t)
	return q
}

func (self *query_t) UnmarshalJSON(data []byte) error {
	var parsed interface{}

	// lexical parsing: get a structure tree form JSON
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("parsed JSON query:\n%+v\n\n%#v\n", parsed, parsed)

	// syntactic parsing: build a tree of queries
	q := new(query_generic_t)
	err = q.parse(parsed)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("query AST:\n%+v\n\n%#v\n", q, q)

	self.query = q
	return nil
}

func (self *query_t) Run(records string, context Context) (string, error) {
	results, err := self.query.Run(records, context)
	if err != nil {
		return "", errgo.Mask(err)
	}
	return results, nil
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
