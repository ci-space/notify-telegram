package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownToHTMLConverter(t *testing.T) {
	cases := []struct {
		Title        string
		BodyMD       string
		ExpectedHTML string
	}{
		{
			Title: "heading with list",
			BodyMD: `## Heading 1
- item 1
- item 2
`,
			ExpectedHTML: `<b>Heading 1</b>
- item 1
- item 2
`,
		},
		{
			Title: "link in list",
			BodyMD: `## Heading 1
- [link](http://12.34)
- [link](http://12.34)
`,
			ExpectedHTML: `<b>Heading 1</b>
- <a href="http://12.34">link</a>
- <a href="http://12.34">link</a>
`,
		},
	}

	converter := NewMarkdownToHTMLConverter()

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			assert.Equal(t, c.ExpectedHTML, converter.Convert(c.BodyMD))
		})
	}
}
