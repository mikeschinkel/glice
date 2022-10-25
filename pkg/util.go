package glice

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"
	"unicode"
)

// StripURLScheme takes a URL and strips the HTTP(S) scheme from it.
func StripURLScheme(u string) string {
	return regexStripScheme.ReplaceAllString(u, "$1")
}

// UpToN accepts a string, a character (byte), and a field number, and
// then returns all character up to but excluding the nth occurrence
// of the character, e.g.:
//
//		UpToN("https://github.com/golang/go/tree/master/src/net", '/', 5)
//			=> "https://github.com/golang/go"
//
func UpToN(str string, ch byte, n int) (s string) {
	var pos int
	if len(str) == 0 {
		goto end
	}
	for pos = 0; pos < len(str); pos++ {
		if str[pos] == ch {
			n--
		}
		if n == 0 {
			break
		}
	}
	s = str[0:pos]
end:
	return s
}

// ExtractField accepts a string and a field number and then returns
// the nth field delimiter by a space, e.g.
//
//		ExtractField("foo bar baz", 2) => "bar"
//
func ExtractField(str string, n int) (f string) {
	return ExtractFieldWithDelimiter(str, n, ' ')
}

// ExtractFieldWithDelimiter accepts a string, a field number, and a field
// delimiter and returns the nth field based on that delimiter, e.g.
//
//		ExtractFieldWithDelimiter("foo|bar|baz", 2, '|') => "bar"
//
func ExtractFieldWithDelimiter(str string, n int, delim byte) (f string) {
	var spaces = n
	var start = -1
	var pos int

	switch n {
	case 1:
		start = 0
		for pos = 0; pos < len(str); pos++ {
			if str[pos] == delim {
				break
			}
		}
	default:
		for pos = 0; pos < len(str); pos++ {
			if str[pos] != delim {
				continue
			}
			spaces--
			if spaces == 1 {
				start = pos + 1
			}
			if spaces == 0 {
				break
			}
		}
		if start == -1 {
			f = ""
			goto end
		}
	}
	f = str[start:pos]
end:
	return f
}

// ContainsString accepts a string slice and potentially contained string
// and returns true if the string was contained in the slice.
func ContainsString(s []string, c string) bool {
	for _, v := range s {
		if v != c {
			continue
		}
		return true
	}
	return false
}

// HTTPGetWithContext makes a GET request using the protocol of the URL scheme and
// does so with a content.Context passed in.
func HTTPGetWithContext(ctx context.Context, _url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", _url, nil)
	if err != nil {
		err = fmt.Errorf("unable to create new HTTP request object for %s; %w", _url, err)
		goto end
	}
	resp, err = http.DefaultClient.Do(req)
end:
	return resp, err
}

// CallerName returns the name of the function which calls this function.
func CallerName() (name string) {
	pc := make([]uintptr, 2)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}

// ReplaceWhitespace accepts a string and a sequence of replace characters — as
// a string, a byte slice, or a byte — and replaces all segments of which space
// with the sequence of replacements characters, e.g.:
//
//		ReplaceWhitespace("a   b   c", '|') => "a|b|c"
//		ReplaceWhitespace("a\t\t\tb\n\n\nc", '---') => "a---b---c"
//
func ReplaceWhitespace[C Chars](inputString string, replaceChars C) string {
	var lastNonWhitespacePos int
	inputLength := len(inputString)
	inputAsByteSlice := []byte(inputString)
	replaceString := string(replaceChars)
	replaceLength := len(replaceString)

	for byteIndex := inputLength - 1; byteIndex >= 0; byteIndex-- {
		if !unicode.IsSpace(rune(inputAsByteSlice[byteIndex])) {
			continue
		}
		lastNonWhitespacePos = byteIndex + 1
		for {
			if !unicode.IsSpace(rune(inputAsByteSlice[byteIndex])) {
				break
			}
			byteIndex--
			if byteIndex == -1 {
				break
			}
		}
		inputLength += replaceLength - lastNonWhitespacePos + byteIndex + 1
		if replaceLength > 1 && len(inputAsByteSlice) < inputLength {
			grownByteSlice := make([]byte, len(inputAsByteSlice)*2)
			copy(grownByteSlice, inputAsByteSlice)
			inputAsByteSlice = grownByteSlice
		}
		copy(inputAsByteSlice[byteIndex+replaceLength+1:], inputAsByteSlice[lastNonWhitespacePos:])
		copy(inputAsByteSlice[byteIndex+1:], replaceString)
	}
	return string(inputAsByteSlice[:inputLength])
}

// Timestamp returns current date/time at UTC as a string in RFC 3339 format.
func Timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
