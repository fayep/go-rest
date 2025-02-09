package rest

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	tests := []struct {
		name          string
		s             string
		inargs        []any
		outputargs    []any
		outputnotargs []any
	}{
		{
			name:          "no args, header",
			s:             "no args",
			inargs:        []any{http.Header{"key": []string{"value"}}},
			outputargs:    nil,
			outputnotargs: []any{http.Header{"key": []string{"value"}}},
		},
		{
			name:          "one arg, body",
			s:             "one arg %s",
			inargs:        []any{"hello", NullCloser{bytes.NewBufferString("world")}},
			outputargs:    []any{"hello"},
			outputnotargs: []any{NullCloser{bytes.NewBufferString("world")}},
		},
		{
			name:          "one floating point arg, body",
			s:             "one arg %-0.2f",
			inargs:        []any{3.14159, NullCloser{bytes.NewBufferString("world")}},
			outputargs:    []any{3.14159},
			outputnotargs: []any{NullCloser{bytes.NewBufferString("world")}},
		},
		{
			name:          "escaped percent sign",
			s:             "100%%-smile %s",
			inargs:        []any{"confident", http.Header{"key": []string{"value"}}},
			outputargs:    []any{"confident"},
			outputnotargs: []any{http.Header{"key": []string{"value"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, notargs := splitArgs(tt.s, tt.inargs...)
			if got, want := args, tt.outputargs; !reflect.DeepEqual(got, want) {
				t.Errorf("got %#v, want %#v", got, want)
			}
			if got, want := notargs, tt.outputnotargs; !reflect.DeepEqual(got, want) {
				t.Errorf("got %#v, want %#v", got, want)
			}
		})
	}
}
