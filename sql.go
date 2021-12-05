//go:build sql_scanner
// +build sql_scanner

package fixed

import (
	"database/sql/driver"
	"fmt"
)

// Scan implements the sql.Scanner interface for database deserialization.
func (f *Fixed) Scan(value interface{}) error {
	// first try to see if the data is stored in database as a Numeric datatype
	switch v := value.(type) {

	case float32:
		*f = NewF(float64(v))
		return nil

	case float64:
		// numeric in sqlite3 sends us float64
		*f = NewF(v)
		return nil

	case int64:
		*f = NewI(v, nPlaces)
		return nil

	default:
		// default is trying to interpret value stored as string
		str, err := unquoteIfQuoted(v)
		if err != nil {
			return err
		}
		val, err := NewSErr(str)
		if err != nil {
			return err
		}
		*f = val
		return nil
	}
}

func unquoteIfQuoted(value interface{}) (string, error) {
	var bytes []byte

	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return "", fmt.Errorf("could not convert value '%+v' to byte array of type '%T'",
			value, value)
	}

	// If the amount is quoted, strip the quotes
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes), nil
}

// Value implements the driver.Valuer interface for database serialization.
func (f Fixed) Value() (driver.Value, error) {
	return f.String(), nil
}
