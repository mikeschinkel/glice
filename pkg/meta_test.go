package glice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ribice/glice/v3/pkg"
)

type mt struct {
	Name     string
	HTML     string
	Import   string
	Expected string
}

var metaTests = []mt{
	{Name: "Has match in 'go-import'", HTML: httpmock, Import: "gopkg.in/jarcoal/httpmock.v1", Expected: "https://github.com/jarcoal/httpmock"},
	{Name: "Has googlesource in 'go-import'", HTML: crypto, Import: "golang.org/x/crypto", Expected: "https://github.com/golang/crypto"},
}

func TestMeta(t *testing.T) {
	var result string
	for _, test := range metaTests {
		t.Run(fmt.Sprintf("%s", test.Name), func(t *testing.T) {
			r := strings.NewReader(test.HTML)
			m, err := glice.GetMetaFromHTMLReader(r)
			if err != nil {
				t.Errorf("failed to get Meta from HTML '%s'; %s", test.Name, err.Error())
			}
			result, err = m.ResolveGoImport(test.Import)
			if result != test.Expected {
				t.Errorf("failed: want='%s', got='%s'", test.Expected, result)
			}
		})
	}
}

// httpmock from https://gopkg.in/jarcoal/httpmock.v1?go-get=1
// Has go-import repo matching prefix
var httpmock = `<html>
<head>
<meta name="go-import" content="gopkg.in/jarcoal/httpmock.v1 git https://gopkg.in/jarcoal/httpmock.v1">
<meta name="go-source" content="gopkg.in/jarcoal/httpmock.v1 _ https://github.com/jarcoal/httpmock/tree/v1.2.0{/dir} https://github.com/jarcoal/httpmock/blob/v1.2.0{/dir}/{file}#L{line}">
</head>
<body>
go get gopkg.in/jarcoal/httpmock.v1
</body>
</html>
`

// crypto from ttps://golang.org/x/crypto?go-get=1
// Has go.googlesource.com in go-import
var crypto = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="golang.org/x/crypto git https://go.googlesource.com/crypto">
<meta name="go-source" content="golang.org/x/crypto https://github.com/golang/crypto/ https://github.com/golang/crypto/tree/master{/dir} https://github.com/golang/crypto/blob/master{/dir}/{file}#L{line}">
<meta http-equiv="refresh" content="0; url=https://pkg.go.dev/golang.org/x/crypto">
</head>
<body>
<a href="https://pkg.go.dev/golang.org/x/crypto">Redirecting to documentation...</a>
</body>
</html>
`
