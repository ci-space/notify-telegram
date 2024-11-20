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

		if entering {
			prev := ast.GetPrevNode(n)
			if prev != nil {
				switch prev.(type) {
				case *ast.HTMLBlock, *ast.List, *ast.Paragraph, *ast.Heading, *ast.CaptionFigure,
					*ast.CodeBlock, *ast.BlockQuote, *ast.Aside, *ast.HorizontalRule:
					r.inner.CR(w)
				}
			}

			if prev == nil {
				_, isParentBlockQuote := n.Parent.(*ast.BlockQuote)
				if isParentBlockQuote {
					r.inner.CR(w)
				}
				_, isParentAside := n.Parent.(*ast.Aside)
				if isParentAside {
					r.inner.CR(w)
				}
			}
		} else if ast.GetNextNode(n) != nil {
			r.inner.CR(w)
		}
	case *ast.Strong:
		r.inner.OutOneOf(w, entering, "<b>", "</b>")
	case *ast.Emph:
		r.inner.OutOneOf(w, entering, "<i>", "</i>")
	case *ast.Heading:
		if entering {
			r.inner.CR(w)
		}
		r.inner.OutOneOf(w, entering, "<b>", "</b>")
	case *ast.List:
	case *ast.ListItem:
		if entering {
			r.inner.CR(w)
			r.inner.Outs(w, "- ")
		} else if n.ListFlags&ast.ListItemEndOfList != 0 {
			r.inner.CR(w)
		}
	default:
		return r.inner.RenderNode(w, node, entering)
	}

	return ast.GoToNext
}

func (r *MarkdownRenderer) RenderHeader(_ io.Writer, _ ast.Node) {}

func (r *MarkdownRenderer) RenderFooter(_ io.Writer, _ ast.Node) {}
