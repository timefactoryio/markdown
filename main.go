package markdown

import (
	"io"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	h "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/toc"
)

type imageUnwrap struct{}

func (t *imageUnwrap) Transform(node *ast.Paragraph, _ text.Reader, _ parser.Context) {
	if node.ChildCount() == 1 {
		if _, ok := node.FirstChild().(*ast.Image); ok {
			if parent := node.Parent(); parent != nil {
				parent.ReplaceChild(parent, node, node.FirstChild())
			}
		}
	}
}

type Markdown struct {
	md goldmark.Markdown
}

func New() *Markdown {
	return &Markdown{md: goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(highlighting.WithStyle("hrdark")),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
			parser.WithParagraphTransformers(util.Prioritized(&imageUnwrap{}, 100)),
		),
		goldmark.WithRendererOptions(
			h.WithXHTML(),
			h.WithUnsafe(),
		),
	)}
}

func (m *Markdown) Convert(src []byte, w io.Writer, withTOC bool) error {
	if !withTOC {
		return m.md.Convert(src, w)
	}
	reader := text.NewReader(src)
	doc := m.md.Parser().Parse(reader)
	tree, _ := toc.Inspect(doc, src, toc.Compact(true))
	if list := toc.RenderList(tree); list != nil {
		doc.InsertBefore(doc, doc.FirstChild(), list)
	}
	return m.md.Renderer().Render(w, src, doc)
}
