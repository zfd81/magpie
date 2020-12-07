package sql

import (
	"strings"
)

const (
	FieldSeparator = "|"
)

type Row struct {
	capacity  int
	data      []string
	timestamp int64
}

func (r *Row) Append(data string) {
	r.data = append(r.data, data)
}

func (r *Row) Set(index int, val string) {
	r.data[index] = val
}

func (r *Row) Get(index int) string {
	return r.data[index]
}

func (r *Row) Data() []string {
	return r.data
}

func (r *Row) String() string {
	return strings.Join(r.data, FieldSeparator)
}

func (r *Row) Load(line, sep string) *Row {
	r.data = strings.SplitN(line, sep, r.capacity)
	return r
}
