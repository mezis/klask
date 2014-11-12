package query

import (
	"github.com/juju/errgo"
	// "github.com/mezis/klask/index"
)

// A query that merges results.
// Represented by a JSON array of generic queries, as the value of an $or key.
type query_or_t struct {
	query_sequence_t
}

func (self *query_or_t) Run(records string, ctx Context) (string, error) {
	var err error = nil

	conn := ctx.Conn()
	defer conn.Close()

	partials := make([]interface{}, len(self.queries))

	for k, query := range self.queries {
		partial, err := query.Run(records, ctx)
		if err != nil {
			return "", errgo.Mask(err)
		}
		partials[k] = partial
	}

	argv := make([]interface{}, 0)
	result, err := ctx.Keys().Get()
	if err != nil {
		return "", errgo.Mask(err)
	}

	argv = append(argv, result, len(partials))
	argv = append(argv, partials...)
	argv = append(argv, "WEIGHTS")
	for _ = range partials {
		argv = append(argv, 0)
	}

	_, err = conn.Do("ZUNIONSTORE", argv...)
	if err != nil {
		return "", errgo.Mask(err)
	}

	// release intermediate sets
	err = ctx.Keys().Release(partials...)
	if err != nil {
		return "", errgo.Mask(err)
	}

	return result, nil
}
