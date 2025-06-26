package md2html

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"strings"
)

func Convert(text string) string {
	text = strings.ReplaceAll(text, "\\n", "\n")

	extensions := parser.CommonExtensions | parser.BackslashLineBreak
	p := parser.NewWithExtensions(extensions)

	renderer := NewMarkdownRenderer()

	return string(markdown.ToHTML([]byte(text), p, renderer))
}
