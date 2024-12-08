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
- item 2`,
			ExpectedHTML: `<b>Heading 1</b>
- item 1
- item 2`,
		},
		{
			Title: "link in list",
			BodyMD: `## Heading 1
- [link](http://12.34)
- [link](http://12.34)`,
			ExpectedHTML: `<b>Heading 1</b>
- <a href="http://12.34">link</a>
- <a href="http://12.34">link</a>`,
		},
		{
			Title: "Two headings",
			BodyMD: `## Heading 1
- item 1
- item 2

## Heading 2
- item 1
- item 2`,
			ExpectedHTML: `<b>Heading 1</b>
- item 1
- item 2

<b>Heading 2</b>
- item 1
- item 2`,
		},
		{
			Title: "Two headings with under text",
			BodyMD: `## Heading 1
- item 1
- item 2

## Heading 2
- item 1
- item 2

some text`,
			ExpectedHTML: `<b>Heading 1</b>
- item 1
- item 2

<b>Heading 2</b>
- item 1
- item 2

some text`,
		},
		{
			Title: "Two headings with under texts",
			BodyMD: `## Heading 1
- item 1
- item 2

## Heading 2
- item 1
- item 2

some text1
some text2

some text3`,
			ExpectedHTML: `<b>Heading 1</b>
- item 1
- item 2

<b>Heading 2</b>
- item 1
- item 2

some text1
some text2

some text3`,
		},
		{
			Title: "text and heading with list",
			BodyMD: `text

## Heading 1
- item 1
- item 2

## Heading 2
- item 1
- item 2`,

			ExpectedHTML: `text

<b>Heading 1</b>
- item 1
- item 2

<b>Heading 2</b>
- item 1
- item 2`,
		},
	}

	converter := NewMarkdownToHTMLConverter()

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			assert.Equal(t, c.ExpectedHTML, converter.Convert(c.BodyMD))
		})
	}
}
