package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

type Int struct {
	sql.NullInt64
}

func NewInt(i int64, valid bool) Int {
	return Int{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

func IntFrom(i int64) Int {
	return NewInt(i, true)
}

func IntFromPtr(i *int64) Int {
	if i == nil {
		return NewInt(0, false)
	}
	return NewInt(*i, true)
}

func (i Int) ValueOrZero() int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

func (i *Int) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		i.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &i.Int64); err != nil {
		return err
	}

	i.Valid = true
	return nil
}

func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	i.Valid = true
	return nil
}

func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

func (i Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

func (i *Int) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

func (i Int) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

func (i Int) IsZero() bool {
	return !i.Valid
}

func (i Int) Equal(other Int) bool {
	return i.Valid == other.Valid && (!i.Valid || i.Int64 == other.Int64)
}
