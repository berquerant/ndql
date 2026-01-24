package node

import (
	"encoding/json"
	"fmt"
	"time"
)

func (*Null) MarshalJSON() ([]byte, error) { return []byte("null"), nil }
func (*Null) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	return fmt.Errorf("cannot unmarshal Null from %s", b)
}

func (v *Float) MarshalJSON() ([]byte, error) { return json.Marshal(v.Raw()) }
func (v *Float) UnmarshalJSON(b []byte) error {
	var x float64
	if err := json.Unmarshal(b, &x); err != nil {
		return fmt.Errorf("%w: cannot unmarshal Float from %s", err, b)
	}
	*v = Float(x)
	return nil
}

func (v *Int) MarshalJSON() ([]byte, error) { return json.Marshal(v.Raw()) }
func (v *Int) UnmarshalJSON(b []byte) error {
	var x int64
	if err := json.Unmarshal(b, &x); err != nil {
		return fmt.Errorf("%w: cannot unmarshal Int from %s", err, b)
	}
	*v = Int(x)
	return nil
}

func (v *Bool) MarshalJSON() ([]byte, error) { return json.Marshal(v.Raw()) }
func (v *Bool) UnmarshalJSON(b []byte) error {
	var x bool
	if err := json.Unmarshal(b, &x); err != nil {
		return fmt.Errorf("%w: cannot unmarshal Bool from %s", err, b)
	}
	*v = Bool(x)
	return nil
}

func (v *String) MarshalJSON() ([]byte, error) { return json.Marshal(v.Raw()) }
func (v *String) UnmarshalJSON(b []byte) error {
	var x string
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	*v = String(x)
	return nil
}

func (v *Time) MarshalJSON() ([]byte, error) { return json.Marshal(v.Raw().Format(time.DateTime)) }
func (v *Time) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	x, err := time.Parse(time.DateTime, s)
	if err != nil {
		return fmt.Errorf("%w: cannot unmarshal Time from %s", err, b)
	}
	*v = Time(x)
	return nil
}

func (v *Duration) MarshalJSON() ([]byte, error) { return json.Marshal(v.Raw().String()) }
func (v *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	x, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("%w: cannot unmarshal Duration from %s", err, b)
	}
	*v = Duration(x)
	return nil
}
