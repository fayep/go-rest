package rest

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var formatRegexp = regexp.MustCompile(`%[-+#0-9\.]*[a-z]`)

// splitArgs separates the arguments that are used in the format string from
// the arguments that are not used in the format string.
// The format string is expected to be a valid format string for fmt.Sprintf.
// The input arguments after any format arguments are expected to be of type io.Reader or http.Header.

func splitArgs(s string, input ...any) ([]any, []any) {
	var args []any
	var notargs []any

	matches := formatRegexp.FindAllStringIndex(s, -1)
	matchpos := 0
	for _, arg := range input {
		for matchpos < len(matches) &&
			s[matches[matchpos][0]-1] == '%' {
			matchpos++
		}
		if matchpos < len(matches) {
			args = append(args, arg)
			matchpos++
		} else {
			switch v := arg.(type) {
			case io.Reader, http.Header:
				notargs = append(notargs, arg)
			default:
				panic(fmt.Errorf("unexpected type after format string: %T", v))
			}
		}
	}
	return args, notargs
}
