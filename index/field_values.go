package index

import (
	"fmt"
	"github.com/juju/errgo"
)

type field_with_int_value_t struct {
}

func (self *field_with_int_value_t) CheckValidValue(value interface{}) (interface{}, error) {
	switch x := value.(type) {
	case int64:
		return x, nil
	case float64:
		return int64(x), nil
	case string:
		var result int
		_, err := fmt.Sscanf(x, "%d", &result)
		if err != nil {
			return nil, errgo.Mask(err)
		}
	}
	return nil, errgo.Newf("bad value of type '%T'; expected 'int'", value)
}
