package helpers

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// StdTime provides default SQL TIME format
const StdTime = "15:04:05"

// emptyTime allows default times to be considered
// null for insertion into the database.
var emptyTime = time.Time{}

// nullLiteral is helpful for checking
// for nulls, as they won't cause errors,
// yet we need the content of the file to change anyway
var nullLiteral = []byte("null")

/*-------------+
| Type aliases |
+-------------*/

// NullFloat64 aliases sql.NullFloat64
type NullFloat64 sql.NullFloat64

// NullString aliases sql.NullString
type NullString sql.NullString

// NullBool aliases sql.NullBool
type NullBool sql.NullBool

// NullInt64 aliases sql.NullInt64
type NullInt64 sql.NullInt64

// NullTime aliases sql.NullTime
type NullTime mysql.NullTime

// RawJSON aliases json.RawMessage
type RawJSON json.RawMessage

/*---------------------------+
| NullString implementations |
+---------------------------*/

// MarshalJSON for NullString
func (n NullString) MarshalJSON() ([]byte, error) {
	var a *string
	if n.Valid {
		a = &n.String
	}
	return json.Marshal(a)
}

// UnmarshalJSON for NullString
func (n *NullString) UnmarshalJSON(b []byte) error {
	if bytes.EqualFold(b, nullLiteral) {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.String)
	n.Valid = err == nil
	return err
}

// Value for NullString
func (n NullString) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String, nil
}

// Scan for NullString
func (n *NullString) Scan(src interface{}) error {
	var a sql.NullString
	if err := a.Scan(src); err != nil {
		return err
	}
	n.String = a.String
	if reflect.TypeOf(src) != nil {
		n.Valid = true
	}
	return nil
}

/*----------------------------+
| NullFloat64 implementations |
+----------------------------*/

// MarshalJSON for NullFloat64
func (n NullFloat64) MarshalJSON() ([]byte, error) {
	var a *float64
	if n.Valid {
		a = &n.Float64
	}
	return json.Marshal(a)
}

// Value for NullFloat64
func (n NullFloat64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Float64, nil
}

// UnmarshalJSON for NullFloat64
func (n *NullFloat64) UnmarshalJSON(b []byte) error {
	if bytes.EqualFold(b, nullLiteral) {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Float64)
	n.Valid = err == nil
	return err
}

// Scan for NullFloat64
func (n *NullFloat64) Scan(src interface{}) error {
	var a sql.NullFloat64
	if err := a.Scan(src); err != nil {
		return err
	}
	n.Float64 = a.Float64
	if reflect.TypeOf(src) != nil {
		n.Valid = true
	}
	return nil
}

/*--------------------------+
| NullInt64 implementations |
+--------------------------*/

// MarshalJSON for NullInt64
func (n NullInt64) MarshalJSON() ([]byte, error) {
	var a *int64
	if n.Valid {
		a = &n.Int64
	}
	return json.Marshal(a)
}

// Value for NullInt64
func (n NullInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

// UnmarshalJSON for NullInt64
func (n *NullInt64) UnmarshalJSON(b []byte) error {
	if bytes.EqualFold(b, nullLiteral) {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Int64)
	n.Valid = err == nil
	return err
}

// Scan for NullInt64
func (n *NullInt64) Scan(src interface{}) error {
	// Set initial state for subsequent scans.
	n.Valid = false

	var a sql.NullInt64
	if err := a.Scan(src); err != nil {
		return err
	}
	n.Int64 = a.Int64
	if reflect.TypeOf(src) != nil {
		n.Valid = true
	}
	return nil
}

/*-------------------------+
| NullBool implementations |
+-------------------------*/

// MarshalJSON for NullBool
func (n NullBool) MarshalJSON() ([]byte, error) {
	var a *bool
	if n.Valid {
		a = &n.Bool
	}
	return json.Marshal(a)
}

// Value for NullBool
func (n NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

// UnmarshalJSON for NullBool
func (n *NullBool) UnmarshalJSON(b []byte) error {
	var field *bool
	err := json.Unmarshal(b, &field)
	if field != nil {
		n.Valid = true
		n.Bool = *field
	}
	return err
}

// Scan for NullBool
func (n *NullBool) Scan(src interface{}) error {
	var a sql.NullBool
	if err := a.Scan(src); err != nil {
		return err
	}
	n.Bool = a.Bool
	if reflect.TypeOf(src) != nil {
		n.Valid = true
	}
	return nil
}

/*-------------------------+
| NullTime implementations |
+-------------------------*/

// MarshalJSON for NullTime
func (n NullTime) MarshalJSON() ([]byte, error) {
	var a *time.Time
	if n.Valid {
		a = &n.Time
	}
	return json.Marshal(a)
}

// Value for NullTime
func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

// UnmarshalJSON for NullTime
func (n *NullTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.Trim(s, `"`)

	var (
		zeroTime time.Time
		tim      time.Time
		err      error
	)

	if strings.EqualFold(s, "null") {
		return nil
	}

	if tim, err = time.Parse(time.RFC3339, s); err != nil {
		n.Valid = false
		return err
	}

	if tim == zeroTime {
		return nil
	}

	n.Time = tim
	n.Valid = true
	return nil
}

// Scan for NullTime
func (n *NullTime) Scan(src interface{}) error {
	// Set initial state for subsequent scans.
	n.Valid = false

	var a mysql.NullTime
	if err := a.Scan(src); err != nil {
		return err
	}
	n.Time = a.Time
	if reflect.TypeOf(src) != nil {
		n.Valid = true
	}
	return nil
}

/*------------------------+
| RawJSON implementations |
+------------------------*/

// MarshalJSON for NullString
func (n RawJSON) MarshalJSON() ([]byte, error) {
	if len(n) == 0 {
		return []byte("null"), nil
	}
	a := json.RawMessage(n)
	return a.MarshalJSON()
}

// Value for NullString
func (n RawJSON) Value() (driver.Value, error) {
	return string(n), nil
}

// UnmarshalJSON for NullString
func (n *RawJSON) UnmarshalJSON(b []byte) error {
	var a json.RawMessage
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	c := RawJSON(a)
	*n = c
	return nil
}

// Scan for NullString
func (n *RawJSON) Scan(src interface{}) error {
	var a sql.NullString
	if err := a.Scan(src); err != nil {
		return err
	}
	jsn := RawJSON([]byte(a.String))
	*n = jsn
	return nil
}

/*-----------------+
| Helper functions |
+-----------------*/

// ToNullString returns a new NullString
func ToNullString(s *string) NullString {
	if s == nil {
		return NullString(sql.NullString{Valid: false})
	}
	return NullString(sql.NullString{String: *s, Valid: true})
}

// ToNullInt64 returns a new NullInt64
func ToNullInt64(i *int64) NullInt64 {
	if i == nil {
		return NullInt64(sql.NullInt64{Valid: false})
	}
	return NullInt64(sql.NullInt64{Int64: *i, Valid: true})
}

// ToNullFloat64 returns a new NullFloat64
func ToNullFloat64(i *float64) NullFloat64 {
	if i == nil {
		return NullFloat64(sql.NullFloat64{Valid: false})
	}
	return NullFloat64(sql.NullFloat64{Float64: *i, Valid: true})
}

// ToNullBool creates a new NullBool
func ToNullBool(b *bool) NullBool {
	if b == nil {
		return NullBool(sql.NullBool{Valid: false})
	}
	return NullBool(sql.NullBool{Bool: *b, Valid: true})
}

// ToNullTime creates a new NullTime
func ToNullTime(t time.Time) NullTime {
	if t == emptyTime {
		return NullTime(mysql.NullTime{Valid: false})
	}
	return NullTime(mysql.NullTime{Time: t, Valid: true})
}
