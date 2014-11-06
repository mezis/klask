package query

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errgo"
	"github.com/mezis/klask/index"
)

type Query interface {
	Run(idx index.Index, targetKey string) error
}

// the toplevel type of query, wraps a single query
type query_t struct {
	idx   index.Index
	query Query
}

func NewQuery(idx index.Index) Query {
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

	fmt.Printf("parsed JSON query:\n%+v\n\n%#v\n", parsed, parsed)

	// start parsing!
	q := new(query_generic_t)
	err = q.parse(self.idx, parsed)
	if err != nil {
		return errgo.Mask(err)
	}

	self.query = q
	return nil
}

func (self *query_t) Run(idx index.Index, targetKey string) error {
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
