package index

import (
	"encoding/json"
	"errors"
	"fmt"
)

type FieldType string

var (
	FIntEq  = FieldType("integer/equality")
	FIntNeq = FieldType("integer/inequality")
)

func (self FieldType) IsValid() bool {
	switch self {
	case FIntEq:
		return true
	case FIntNeq:
		return true
	default:
		return false
	}
}

func (self *FieldType) UnmarshalJSON(data []byte) error {
	var name string
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}

	result := FieldType(name)
	if result.IsValid() {
		*self = result
		return nil
	} else {
		fmt.Println("offending type", name)
		return errors.New("ha!")
	}
}
