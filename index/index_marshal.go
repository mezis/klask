package index

import (
	"encoding/json"
	// "fmt"
)

type field_marshal_t struct {
	Name string    `json:"name"`
	Type FieldType `json:"type"`
}

type index_marshal_t struct {
	ID     string            `json:"id"`
	Fields []field_marshal_t `json:"fields"`
}

// Satisfiy the json.Marhshaler interface
func (self *index_t) MarshalJSON() ([]byte, error) {
	var record index_marshal_t

	// fmt.Println("converting index to JSON")
	record.ID = self.name
	record.Fields = make([]field_marshal_t, len(self.fields))
	k := 0
	for _, field := range self.fields {
		// fmt.Println("  field: ", field.Name(), field.Type())
		record.Fields[k] = field_marshal_t{field.Name(), field.Type()}
		k++
	}

	// data, err := json.MarshalIndent(record, "", "  ")
	// fmt.Println(string(data))
	data, err := json.Marshal(record)
	return data, err
}

func (self *index_t) UnmarshalJSON(data []byte) error {
	var record index_marshal_t

	// fmt.Println("parsing index: ", string(data))
	err := json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	self.name = record.ID
	for _, val := range record.Fields {
		err := self.AddField(val.Name, val.Type)
		if err != nil {
			return err
		}
	}

	return nil
}
