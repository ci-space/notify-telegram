package internal

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"strings"
)

type MarkdownToHTMLConverter struct {
	renderer *MarkdownRenderer
}

func NewMarkdownToHTMLConverter() *MarkdownToHTMLConverter {
	return &MarkdownToHTMLConverter{renderer: NewMarkdownRenderer()}
}

func (c *MarkdownToHTMLConverter) Convert(text string) string {
	text = strings.ReplaceAll(text, "\\n", "\n")

	extensions := parser.CommonExtensions | parser.BackslashLineBreak
	p := parser.NewWithExtensions(extensions)

	renderer := NewMarkdownRenderer()

	return string(markdown.ToHTML([]byte(text), p, renderer))
}
