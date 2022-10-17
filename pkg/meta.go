package glice

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

var chkSourceDomains = []string{
	"go.googlesource.com",
}
var gitHubLikeDomains = []string{
	"github.com",
	"bitbucket.com",
}

var errRequestPrefixInstead = errors.New("prefix is a subset of import; request prefix instead")

type Meta struct {
	Import string `meta:"go-import"`
	Source string `meta:"go-source"`
}

func GetMetaFromHTMLReader(doc io.Reader) (meta *Meta, err error) {
	var found int8
	meta = &Meta{}
	tokenizer := html.NewTokenizer(doc)
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			if tokenizer.Err() != io.EOF {
				err = fmt.Errorf("HTML parse error: %w", tokenizer.Err())
			}
			goto end
		}

		token := tokenizer.Token()

		if token.Type == html.EndTagToken && token.Data == htmlRegion {
			goto end
		}

		if token.Data == htmlTag {
			var property, content string
			for _, attr := range token.Attr {
				switch attr.Key {
				case "property", "name":
					property = strings.TrimSpace(strings.ToLower(attr.Val))
				case "content":
					content = strings.TrimSpace(attr.Val)
				}
			}
			switch property {
			case importName:
				meta.Import = ReplaceWhitespace(content, ' ')
				found++
			case sourceName:
				meta.Source = ReplaceWhitespace(content, ' ')
				found++
			}
		}
		if found == 2 {
			goto end
		}
	}
end:
	return meta, nil
}

func (m *Meta) ResolveGoSource() (repoURL string, err error) {
	// The 3rd field gives access to the repository
	// Examples of all fields:
	//    golang.org/x/crypto https://github.com/golang/crypto/ https://github.com/golang/crypto/tree/master{/dir} https://github.com/golang/crypto/blob/master{/dir}/{file}#L{line}
	//    gopkg.in/jarcoal/httpmock.v1 _ https://github.com/jarcoal/httpmock/tree/v1.2.0{/dir} https://github.com/jarcoal/httpmock/blob/v1.2.0{/dir}/{file}#L{line}
	//
	_url := ExtractField(m.Source, 2)
	if "_" == _url {
		_url = ExtractField(m.Source, 3)
	}
	domain := StripURLScheme(UpToN(_url, '/', 3))
	if !ContainsString(gitHubLikeDomains, domain) {
		err = fmt.Errorf("repo domain '%s' not handled yet [Full URL=%s]", domain, _url)
		goto end
	}
	// Get up to but not including the URL path
	repoURL = UpToN(_url, '/', 5)
end:
	return repoURL, err
}

func (m *Meta) ResolveGoImport(imp string) (repoURL string, err error) {
	var msg string
	var chkSource bool

	prefix := ExtractField(m.Import, 1)

	switch {
	case len(prefix) < len(imp):
		if imp[0:len(prefix)] == prefix {
			repoURL = fmt.Sprintf("https://%s", prefix)
			err = errRequestPrefixInstead
			goto end
		}
		msg = "does not match a subset of"

	case len(prefix) > len(imp):
		msg = "is longer than"

	default:
		repoURL = ExtractField(m.Import, 3)
		repo := StripURLScheme(repoURL)
		if repo == imp {
			chkSource = true
		}
		if ContainsString(chkSourceDomains, UpToN(repo, '/', 1)) {
			chkSource = true
		}
		if chkSource {
			// The import and repo are the same. We need to look in <meta name="go-source">
			// Example: https://gopkg.in/jarcoal/httpmock.v1?go-get=1
			repoURL, err = m.ResolveGoSource()
			if err != nil {
				err = fmt.Errorf("unable to resolve import '%s'; %w", imp, err)
				repoURL = ""
			}
			goto end
		}
		// We found the URL in <meta name="go-import">

		goto end
	}

	// Anything but switch-default generates this error:
	err = fmt.Errorf("remote data error in <meta name='%s'>; prefix '%s' %s import '%s'",
		importName,
		prefix,
		msg,
		imp)
	repoURL = ""

end:
	return repoURL, err
}
