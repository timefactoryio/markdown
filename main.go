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

// Browse available themes: https://xyproto.github.io/splash/docs/all.html
const DefaultTheme = "github-dark"

// New creates a goldmark instance with the given Chroma theme name.
// Browse available themes: https://xyproto.github.io/splash/docs/all.html
func New(theme string) *Markdown {
	if theme == "" {
		theme = DefaultTheme
	}
	return &Markdown{md: goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(highlighting.WithStyle(theme)),
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
