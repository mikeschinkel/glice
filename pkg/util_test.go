package glice_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ribice/glice/v3/pkg"
)

const (
	ByteType   = reflect.Uint8
	StringType = reflect.String
	ByteSlice  = reflect.Slice
)

var replaceWhitespaceTests = []struct {
	String  string
	Replace string
	Result  string
	Type    reflect.Kind
}{
	{
		String:  "   foo   bar    ",
		Replace: "___",
		Result:  "___foo___bar___",
		Type:    StringType,
	},
	{
		String:  "         ",
		Replace: "x",
		Result:  "x",
		Type:    StringType,
	},
	{
		String:  "",
		Replace: "xxx",
		Result:  "",
		Type:    StringType,
	},
	{
		String:  "         ",
		Replace: "xxx",
		Result:  "xxx",
		Type:    StringType,
	},
	{
		String:  "a\nb\tc",
		Replace: "xxx",
		Result:  "axxxbxxxc",
		Type:    ByteSlice,
	},
	{
		String:  "a\nb\tc",
		Replace: " ",
		Result:  "a b c",
		Type:    ByteType,
	},
	{
		String:  "abc",
		Replace: "xxx",
		Result:  "abc",
		Type:    StringType,
	},
	{
		String:  "",
		Replace: "shazam",
		Result:  "",
		Type:    StringType,
	},
	{
		String:  "a\t\n\v\f\rb",
		Replace: " ",
		Result:  "a b",
		Type:    StringType,
	},
	{
		String:  "foo    bar     baz",
		Replace: " ",
		Result:  "foo bar baz",
		Type:    StringType,
	},
	{
		String:  "foo bar baz",
		Replace: " ",
		Result:  "foo bar baz",
		Type:    StringType,
	},
}

func TestReplaceWhitespace(t *testing.T) {
	var result string
	for _, _t := range replaceWhitespaceTests {
		t.Run(_t.String, func(t *testing.T) {
			switch _t.Type {
			case ByteType:
				result = glice.ReplaceWhitespace(_t.String, _t.Replace[0])
			case ByteSlice:
				result = glice.ReplaceWhitespace(_t.String, []byte(_t.Replace))
			case StringType:
				fallthrough
			default:
				result = glice.ReplaceWhitespace(_t.String, _t.Replace)
			}
			if result != _t.Result {
				t.Errorf("failed: want='%s', got='%s'", _t.Result, result)
			}
		})
	}
}

var getContainsStringTests = []struct {
	Strings   []string
	Substring string
	Result    bool
}{
	{
		Strings:   []string{},
		Substring: "foo",
		Result:    false,
	},
	{
		Strings:   []string{"foo", "bar", "baz"},
		Substring: "bar",
		Result:    true,
	},
	{
		Strings:   []string{"foo", "bar", "baz"},
		Substring: "bazoom",
		Result:    false,
	},
}

func TestContainsString(t *testing.T) {
	for _, _t := range getContainsStringTests {
		t.Run(_t.Substring, func(t *testing.T) {
			result := glice.ContainsString(_t.Strings, _t.Substring)
			if result != _t.Result {
				t.Errorf("failed: want='%t', got='%t'", _t.Result, result)
			}
		})
	}
}

var getStripURLSchemeTests = []struct {
	URL    string
	Result string
}{
	{
		URL:    "https://github.com/jarcoal/httpmock",
		Result: "github.com/jarcoal/httpmock",
	},
	{
		URL:    "http" + "://github.com/jarcoal/httpmock",
		Result: "github.com/jarcoal/httpmock",
	},
	{
		URL:    "github.com/jarcoal/httpmock",
		Result: "github.com/jarcoal/httpmock",
	},
}

func TestStripURLScheme(t *testing.T) {
	var result string
	for _, _t := range getStripURLSchemeTests {
		t.Run(_t.URL, func(t *testing.T) {
			result = glice.StripURLScheme(_t.URL)
			if result != _t.Result {
				t.Errorf("failed: want='%s', got='%s'", _t.Result, result)
			}
		})
	}
}

var upToNTests = []struct {
	Name       string
	String     string
	Char       byte
	Occurrence int
	Result     string
}{
	{String: "https://github.com/jarcoal/httpmock/tree/v1.2.0{/dir}", Char: '/', Occurrence: 5, Result: "https://github.com/jarcoal/httpmock"},
	{String: "foo bar baz bazoom", Char: ' ', Occurrence: 3, Result: "foo bar baz"},
	{String: "", Char: 'x', Occurrence: 3, Result: ""},
	{String: "abcdefghi", Char: 'x', Occurrence: 1, Result: "abcdefghi"},
}

func TestUpToN(t *testing.T) {
	var result string
	for _, _t := range upToNTests {
		name := fmt.Sprintf("%s[%d%c]", _t.String, _t.Occurrence, _t.Char)
		t.Run(name, func(t *testing.T) {
			result = glice.UpToN(_t.String, _t.Char, _t.Occurrence)
			if result != _t.Result {
				t.Errorf("failed: want='%s', got='%s'", _t.Result, result)
			}
		})
	}
}

var extractFieldTests = []struct {
	Content string
	Field   int
	Result  string
}{
	{Content: "foo", Field: 1, Result: "foo"},
	{Content: "foo bar baz", Field: 1, Result: "foo"},
	{Content: "", Field: 2, Result: ""},
	{Content: "", Field: 1, Result: ""},
	{Content: "foo bar baz", Field: 4, Result: ""},
	{Content: "foo bar baz", Field: 3, Result: "baz"},
	{Content: "foo bar baz bazoom", Field: 3, Result: "baz"},
}

func TestExtractField(t *testing.T) {
	var result string
	for _, _t := range extractFieldTests {
		name := fmt.Sprintf("%s#%d==%s", _t.Content, _t.Field, _t.Result)
		t.Run(name, func(t *testing.T) {
			result = glice.ExtractField(_t.Content, _t.Field)
			if result != _t.Result {
				t.Errorf("failed: want='%s', got='%s'", _t.Result, result)
			}
		})
	}
}
