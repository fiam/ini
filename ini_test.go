package ini

import (
	"reflect"
	"strings"
	"testing"
)

type iniTest struct {
	text   string
	expect map[string]string
	err    string
}

func TestIni(t *testing.T) {
	iniTests := []*iniTest{
		{"a = b  \n 3 = 7", map[string]string{"a": "b", "3": "7"}, ""},
		{"a = b  \r\n 3 = 7", map[string]string{"a": "b", "3": "7"}, ""},
		{"a = b  \r\n 3 = 7=7", map[string]string{"a": "b", "3": "7=7"}, ""},
		{"a = multiline\\\n value  \n 3 = 7", map[string]string{"a": "multiline value", "3": "7"}, ""},
		{"3 = 7\ninvalid", map[string]string{"a": "multiline value", "3": "7"}, "invalid line 2 \"invalid\" - missing separator \"=\""},
	}
	for _, v := range iniTests {
		res, err := Parse(strings.NewReader(v.text))
		if err != nil {
			if v.err != err.Error() {
				if v.err == "" {
					t.Errorf("unexpected error parsing %q: %s", v.text, err)
				} else {
					t.Errorf("expecting error %s parsing %q, got %s instead", v.err, v.text, err)
				}
			}
		} else {
			if v.err != "" {
				t.Errorf("expecting error %s parsing %q, got no error instead", v.err, v.text)
			} else {
				if !reflect.DeepEqual(v.expect, res) {
					t.Errorf("expecting %v parsing %q, got %v instead", v.expect, v.text, res)
				}
			}
		}
	}
}
