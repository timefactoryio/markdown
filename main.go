package markdown

import (
	"io"

	"github.com/alecthomas/chroma/v2"
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

type Style struct {
	Background  string
	Foreground  string
	Error       string
	Keyword     string
	Operator    string
	Punctuation string
	Name        string
	Literal     string
	Comment     string
	Generic     string
}

var Default = Style{
	Background:  "#1d2432",
	Foreground:  "#ffffff",
	Error:       "#ff636f",
	Keyword:     "#ff636f",
	Operator:    "#ff636f",
	Punctuation: "#ffffff",
	Name:        "#58a1dd",
	Literal:     "#a6be9d",
	Comment:     "italic #828b96",
	Generic:     "#ffffff",
}

func (s Style) style() *chroma.Style {
	return chroma.MustNewStyle("theme", chroma.StyleEntries{
		chroma.Background:  s.Foreground + " bg:" + s.Background,
		chroma.Other:       s.Foreground,
		chroma.Error:       s.Error,
		chroma.Keyword:     s.Keyword,
		chroma.Operator:    s.Operator,
		chroma.Punctuation: s.Punctuation,
		chroma.Name:        s.Name,
		chroma.Literal:     s.Literal,
		chroma.Comment:     s.Comment,
		chroma.Generic:     s.Generic,
	})
}

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

func New(style Style) *Markdown {
	return &Markdown{md: goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(highlighting.WithCustomStyle(style.style())),
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
