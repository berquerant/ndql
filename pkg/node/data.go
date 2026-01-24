package node

import (
	"fmt"
	"time"
)

// Data is the smallest unit of data in ndql.
type Data interface {
	IsData()
	Display() string
	AsOp() *Op
	Any() any
}

type (
	Null     struct{}
	Float    float64
	Int      int64
	Bool     bool
	String   string
	Time     time.Time
	Duration time.Duration
)

func NewNull() Null { return Null{} }

func (Null) IsData()     {}
func (Float) IsData()    {}
func (Int) IsData()      {}
func (Bool) IsData()     {}
func (String) IsData()   {}
func (Time) IsData()     {}
func (Duration) IsData() {}

func (Null) Display() string       { return "Null" }
func (v Float) Display() string    { return fmt.Sprintf("Float(%f)", v.Raw()) }
func (v Int) Display() string      { return fmt.Sprintf("Int(%d)", v.Raw()) }
func (v Bool) Display() string     { return fmt.Sprintf("Bool(%v)", v.Raw()) }
func (v String) Display() string   { return fmt.Sprintf("String(%s)", v.Raw()) }
func (v Time) Display() string     { return fmt.Sprintf("Time(%v)", v.Raw()) }
func (v Duration) Display() string { return fmt.Sprintf("Duration(%v)", v.Raw()) }

func (v Null) AsOp() *Op     { return &Op{v} }
func (v Float) AsOp() *Op    { return &Op{v} }
func (v Int) AsOp() *Op      { return &Op{v} }
func (v Bool) AsOp() *Op     { return &Op{v} }
func (v String) AsOp() *Op   { return &Op{v} }
func (v Time) AsOp() *Op     { return &Op{v} }
func (v Duration) AsOp() *Op { return &Op{v} }

func (v Null) Raw() Null              { return v }
func (v Float) Raw() float64          { return float64(v) }
func (v Int) Raw() int64              { return int64(v) }
func (v Bool) Raw() bool              { return bool(v) }
func (v String) Raw() string          { return string(v) }
func (v Time) Raw() time.Time         { return time.Time(v) }
func (v Duration) Raw() time.Duration { return time.Duration(v) }

func (Null) Any() any       { return nil }
func (v Float) Any() any    { return v.Raw() }
func (v Int) Any() any      { return v.Raw() }
func (v Bool) Any() any     { return v.Raw() }
func (v String) Any() any   { return v.Raw() }
func (v Time) Any() any     { return v.Raw() }
func (v Duration) Any() any { return v.Raw() }

// DefaultData is the factory of the default value of Data.
type DefaultData struct{}

func (DefaultData) Null() Null     { return NewNull() }
func (DefaultData) Float() Float   { return Float(0) }
func (DefaultData) Int() Int       { return Int(0) }
func (DefaultData) Bool() Bool     { return Bool(false) }
func (DefaultData) String() String { return String("") }
func (DefaultData) Time() Time {
	var t time.Time
	return Time(t)
}
func (DefaultData) Duration() Duration { return Duration(0) }

func Default() *DefaultData { return &DefaultData{} }
