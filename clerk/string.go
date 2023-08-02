package clerk

import (
	"encoding/json"
)

// String represents all the different states that a value can take
// in json. It is either an unset attribute, an attribute with a
// null value or a string.
type String struct {
	Value  string
	IsNull bool
	IsSet  bool
}

func NewString(value string) String {
	return String{
		Value:  value,
		IsNull: false,
		IsSet:  true,
	}
}

func NewNullString() String {
	return String{
		Value:  "",
		IsNull: true,
		IsSet:  true,
	}
}

func (s *String) Ptr() *String {
	if s.IsSet {
		return s
	}

	return nil
}

func (s *String) UnmarshalJSON(data []byte) error {
	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	s.IsNull = value == nil
	if s.IsNull {
		s.Value = ""
	} else {
		s.Value = value.(string)
	}
	s.IsSet = true

	return nil
}

func (s *String) MarshalJSON() ([]byte, error) {
	// This is called only if attribute is set.
	if s.IsNull {
		return json.Marshal(nil)
	}

	return json.Marshal(&s.Value)
}
