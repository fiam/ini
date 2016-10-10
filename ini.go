package ini

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"github.com/fiam/stringutil"
)

// Options specify the options for ParseOptions.
type Options struct {
	// Separator indicates the characters used as key-value separator.
	// If empty, "=" is used.
	Separator string
	// Comment indicates the characters used to check if a line is a comment.
	// Lines starting with any character in this string are ignored.
	// If empty, all lines are parsed.
	Comment string
}

// Parse parses a .ini style file in the form:
//
//  key1 = value
//  key2 = value
//
// Values that contain newlines ("\n" or "\r\n") need to
// be escaped by ending the previous line with a '\'
// character. Lines starting with ';' or '#' are
// considered comments and ignored. Empty lines are ignored
// too. If a non-empty, non-comment line does not contain
// a '=' an error is returned.
func Parse(r io.Reader) (map[string]string, error) {
	return ParseOptions(r, nil)
}

// ParseOptions works like Parse, but allows the caller to specify
// the strings which represent separators and comments. If opts is nil, this
// function acts like Parse. If Separator is empty, it defaults to '='. If
// Comment is empty, no lines are considered comments.
func ParseOptions(r io.Reader, opts *Options) (map[string]string, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading ini input: %s", err)
	}
	var separator string
	var comment string
	if opts != nil {
		separator = opts.Separator
		comment = opts.Comment
	} else {
		comment = ";#"
	}
	isComment := stringutil.RuneChecker(comment)
	if separator == "" {
		separator = "="
	}
	isSeparator := stringutil.RuneChecker(separator)
	lines := stringutil.SplitLines(string(data))
	values := make(map[string]string, len(lines))
	for ii, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		first, _ := utf8.DecodeRuneInString(line)
		if isComment(first) {
			continue
		}
		sep := -1
		for jj := 0; jj < len(line); jj++ {
			if isSeparator(rune(line[jj])) {
				sep = jj
				break
			}
		}
		if sep < 0 {
			return nil, fmt.Errorf("invalid line %d %q - missing separator %q", ii+1, line, separator)
		}
		key := strings.TrimSpace(line[:sep])
		value := strings.TrimSpace(line[sep+1:])
		values[key] = value
	}
	return values, nil
}
