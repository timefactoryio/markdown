# markdown

`github.com/timefactoryio/markdown`

A thin configuration layer over [goldmark](https://github.com/yuin/goldmark). This package pins goldmark's behavior at the configuration layer rather than in a fork. Goldmark is declared once, configured once, and the output is frozen until this package is deliberately changed. 

## Install

```
go get github.com/timefactoryio/markdown
```

## Usage

```go
import "github.com/timefactoryio/markdown"

md := markdown.New()

// Render without table of contents
if err := md.Convert(src, w, false); err != nil {
    // handle
}

// Render with table of contents prepended
if err := md.Convert(src, w, true); err != nil {
    // handle
}
```

`src` is `[]byte` Markdown source. `w` is any `io.Writer`. The `bool` controls whether a compact table of contents is prepended to the output for that document.

## Output

| Feature             | Detail                                                                    |
| ------------------- | ------------------------------------------------------------------------- |
| GFM                 | Tables, task lists, strikethrough, auto-linked URLs                       |
| Syntax highlighting | Fenced code blocks highlighted at convert time via Chroma, `hrdark` style |
| Image unwrap        | Standalone images rendered as bare `<img>` — no `<p>` wrapper             |
| Raw HTML            | `WithUnsafe` enabled — HTML in source passes through as-is                |
| TOC                 | Optional per call — compact, heading-linked, auto-generated IDs           |
| XHTML               | Self-closing void elements (`<br />`, `<img />`)                          |

## Dependencies

| Package                                    | Version | Purpose                    |
| ------------------------------------------ | ------- | -------------------------- |
| `github.com/yuin/goldmark`                 | v1.8.2  | Parser and renderer        |
| `github.com/yuin/goldmark-highlighting/v2` | v2.0.0  | Chroma syntax highlighting |
| `go.abhg.dev/goldmark/toc`                 | v0.12.0 | Table of contents          |
