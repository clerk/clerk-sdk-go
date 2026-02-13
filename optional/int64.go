// Credits
// https://www.calhoun.io/how-to-determine-if-a-json-key-has-been-set-to-null-or-not-provided/

package optional

import (
	"encoding/json"
	"reflect"
)

// Int64 represents an optional nullable int64 for PATCH semantics:
// IsSet false = key absent (do not change); IsSet true and Valid false = null (reset); IsSet true and Valid true = value.
type Int64 struct {
	Value int64 // The value of the field
	Valid bool  // Whether the value is present (i.e. not null)
	IsSet bool  // Whether the key was set in the payload or not
}

var (
	_ json.Marshaler   = (*Int64)(nil)
	_ json.Unmarshaler = (*Int64)(nil)
)

func (i *Int64) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.IsSet = true

	if string(data) == "null" {
		i.Valid = false
		return nil
	}

	var temp int64
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Valid = true
	return nil
}

func (i Int64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(i.Value)
}

// Int64FromPtr creates a new json.Int64
func Int64FromPtr(i *int64) Int64 {
	if i == nil {
		return Int64{0, false, true}
	}
	return Int64{*i, true, true}
}

// Int64Valuer extracts the value for this custom type so that go validator can perform its checks on it
func Int64Valuer(field reflect.Value) any {
	if i, ok := field.Interface().(Int64); ok {
		return i.Value
	}
	return nil
}
