package nullable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Time struct {
	sql.NullTime
}

func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

func NewTime(t time.Time, valid bool) Time {
	return Time{
		NullTime: sql.NullTime{
			Time:  t,
			Valid: valid,
		},
	}
}

func TimeFrom(t time.Time) Time {
	return NewTime(t, true)
}

func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Time{}, false)
	}
	return NewTime(*t, true)
}

func (t Time) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalJSON()
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		t.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &t.Time); err != nil {
		return err
	}

	t.Valid = true
	return nil
}

func (t Time) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return t.Time.MarshalText()
}

func (t *Time) UnmarshalText(text []byte) error {
	str := string(text)
	// allowing "null" is for backwards compatibility with v3
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	t.Valid = true
	return nil
}

func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func (t Time) IsZero() bool {
	return !t.Valid
}

func (t Time) Equal(other Time) bool {
	return t.Valid == other.Valid && (!t.Valid || t.Time.Equal(other.Time))
}

func (t Time) ExactEqual(other Time) bool {
	return t.Valid == other.Valid && (!t.Valid || t.Time == other.Time)
}
