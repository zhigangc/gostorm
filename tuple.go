package gostorm

import "fmt"

type Tuple struct {
	Id string
	Component string
	Stream string
	Task int
	Values []interface{}
}

func NewTuple(vals Values, comp *Component) *Tuple {
	tup := &Tuple {}
	if val, ok := vals.GetString("id"); ok {
		tup.Id = val
	}
	if val, ok := vals.GetString("comp"); ok {
		tup.Component = val
	}
	if val, ok := vals.GetString("stream"); ok {
		tup.Stream = val
	}
	if val, ok := vals.GetInt("task"); ok {
		tup.Task = val
	}
	if val, ok := vals.GetList("tuple"); ok {
		tup.Values = val
	}
	return tup
}

func (t *Tuple) String() string {
	return fmt.Sprintf("<Tuple component=%s id=%s stream=%s task=%d values=%v>",
		t.Component,
		t.Id,
		t.Stream,
		t.Task,
		t.Values)
}