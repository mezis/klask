package query

import (
	"github.com/garyburd/redigo/redis"
	"github.com/mezis/klask/index"
	"github.com/mezis/klask/util/tempkey"
)

type Context interface {
	Conn() redis.Conn
	Keys() tempkey.Keys
	Idx() index.Index
}

type context_t struct {
	idx  index.Index
	keys tempkey.Keys
}

func newContext(idx index.Index) *context_t {
	rv := new(context_t)
	rv.idx = idx
	rv.keys = tempkey.New(idx.Conn)
	return rv
}

func (self *context_t) Conn() redis.Conn {
	return self.idx.Conn()
}

func (self *context_t) Idx() index.Index {
	return self.idx
}

func (self *context_t) Keys() tempkey.Keys {
	return self.keys
}
