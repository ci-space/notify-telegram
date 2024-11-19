package internal

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

type MarkdownRenderer struct {
	inner *html.Renderer
}

func NewMarkdownRenderer() *MarkdownRenderer {
	return &MarkdownRenderer{inner: html.NewRenderer(html.RendererOptions{})}
}

func (r *MarkdownRenderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch n := node.(type) {
	case *ast.Paragraph:
		if html.SkipParagraphTags(n) {
			return ast.GoToNext
		}

		r.inner.Outs(w, "\n")
	case *ast.Strong:
		r.inner.OutOneOf(w, entering, "<b>", "</b>")
	case *ast.Emph:
		r.inner.OutOneOf(w, entering, "<i>", "</i>")
	case *ast.Heading:
		r.inner.OutOneOf(w, entering, "<b>", "</b>")
	case *ast.List:
		if entering {
			r.inner.Outs(w, "\n")
		}
	case *ast.ListItem:
		if entering {
			r.inner.Outs(w, "- ")
		} else {
			r.inner.Outs(w, "\n")
		}
	default:
		return r.inner.RenderNode(w, node, entering)
	}

	return ast.GoToNext
}

func (r *MarkdownRenderer) RenderHeader(_ io.Writer, _ ast.Node) {}

func (r *MarkdownRenderer) RenderFooter(_ io.Writer, _ ast.Node) {}
