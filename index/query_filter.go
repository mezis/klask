package index

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errgo"
	"regexp"
)

const (
	OperatorEq  = "eq"
	OperatorNeq = "neq"
	OperatorIn  = "in"
	OperatorNi  = "ni"
)

var OperatorRE = regexp.MustCompile(`^eq|neq|in|ni$`)

type filter_t struct {
	idx       Index
	fieldName string
	field     Field
	operator  string
	operand   interface{}
}

type raw_filter_t [3]interface{}

func (self *filter_t) UnmarshalJSON(data []byte) error {
	var raw []interface{}

	fmt.Printf("%T: %v\n", string(data), string(data))

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Println("unmarshalled!")
	fmt.Printf("%T: %v\n", raw, raw)

	switch x := raw[0].(type) {
	case string:
		self.fieldName = x
	default:
		return errgo.Newf("bad field name '%v'", x)
	}

	switch x := raw[1].(type) {
	case string:
		self.operator = x
	default:
		return errgo.Newf("bad operator '%v'", x)
	}
	self.operand = raw[2]

	return nil
}

func (self *filter_t) Validate(idx Index) error {
	self.idx = idx

	// validate field
	field, err := self.idx.Field(self.fieldName)
	if err != nil {
		return errgo.Mask(err)
	}
	self.field = field

	// validate operator
	if !OperatorRE.MatchString(self.operator) {
		return errgo.Newf("bad filter: operator '%s' unknown", self.operator)
	}
	if !self.isValidOp(self.operator) {
		return errgo.Newf("bad filter: operator '%s' invalid for type '%s'",
			self.operator, self.field.Type())
	}

	// TODO: check operand

	return nil
}

func (self *filter_t) isValidOp(op string) bool {
	switch self.field.Type() {
	case FIntEq:
		switch self.operator {
		case OperatorEq:
		case OperatorNeq:
		default:
			return false
		}
	case FIntNeq:
		switch self.operator {
		case OperatorIn:
		case OperatorNi:
		default:
			return false
		}
	default:
		return false
	}
	return true
}

func (self *filter_t) run(sourceKey string, targetKey string) error {
	err := self.field.Filter(self.operator, self.operand, sourceKey, targetKey)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}
