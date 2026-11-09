#!/usr/bin/env bash
#
# Render the ABP documentation pages into a single browsable HTML file that
# approximates how they will look on the Terraform Registry, for text review.
#
# There is no official CLI for this: `tfplugindocs` only generates Markdown
# (see scripts/gen-abp-docs.sh), and the Registry's own doc preview
# (https://registry.terraform.io/tools/doc-preview) is a web form that takes
# one page at a time. This script renders every `website/docs/{r,d}/abp_*`
# page locally with goldmark (GitHub-flavored Markdown, same dialect the
# Registry uses) so all of them can be read in one pass.
#
# Usage:
#   scripts/preview-abp-docs.sh          # render and print the output path
#   OPEN=1 scripts/preview-abp-docs.sh   # also open it in the default browser
#
# Output: abp-extras/abp-docs-preview.html (one self-contained file, no assets)

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

OUT_FILE="abp-extras/abp-docs-preview.html"
TOOL_DIR=".tmp/abp-docs-preview-tool"   # .tmp/ is gitignored

mkdir -p "$TOOL_DIR"

cat > "$TOOL_DIR/go.mod" <<'EOF'
module abpdocspreview

go 1.22

require github.com/yuin/goldmark v1.7.7
EOF

cat > "$TOOL_DIR/main.go" <<'EOF'
// Renders website/docs/{r,d}/abp_* into one HTML file for review.
package main

import (
	"bytes"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	ghtml "github.com/yuin/goldmark/renderer/html"
)

type page struct {
	Slug        string
	Title       string
	Subcategory string
	Description string
	Source      string
	Body        string
}

var (
	idRe   = regexp.MustCompile(`id="([^"]+)"`)
	hrefRe = regexp.MustCompile(`href="#([^"]+)"`)
)

// frontMatter splits the YAML front matter from the Markdown body and pulls
// out the few keys the Registry actually consumes.
func frontMatter(src string) (title, subcat, desc, body string) {
	lines := strings.Split(src, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return "", "", "", src
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end == -1 {
		return "", "", "", src
	}

	inDesc := false
	var descLines []string
	for _, l := range lines[1:end] {
		switch {
		case strings.HasPrefix(l, "page_title:"):
			inDesc = false
			title = strings.Trim(strings.TrimSpace(strings.TrimPrefix(l, "page_title:")), `"`)
		case strings.HasPrefix(l, "subcategory:"):
			inDesc = false
			subcat = strings.Trim(strings.TrimSpace(strings.TrimPrefix(l, "subcategory:")), `"`)
		case strings.HasPrefix(l, "description:"):
			inDesc = true
			rest := strings.TrimSpace(strings.TrimPrefix(l, "description:"))
			if rest != "" && rest != "|-" && rest != "|" {
				descLines = append(descLines, strings.Trim(rest, `"`))
				inDesc = false
			}
		case inDesc:
			descLines = append(descLines, strings.TrimSpace(l))
		default:
			inDesc = false
		}
	}
	return title, subcat, strings.TrimSpace(strings.Join(descLines, " ")), strings.Join(lines[end+1:], "\n")
}

func collect(dir, kind string, md goldmark.Markdown) ([]page, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "abp_*"))
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)

	var pages []page
	for _, path := range matches {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		base := filepath.Base(path)
		base = strings.TrimSuffix(base, ".markdown")
		base = strings.TrimSuffix(base, ".html")

		title, subcat, desc, body := frontMatter(string(raw))
		if title == "" {
			title = "incapsula_" + base
		}

		var buf bytes.Buffer
		if err := md.Convert([]byte(body), &buf); err != nil {
			return nil, err
		}

		// Anchors inside a page (`#nestedblock--kind`, `#required`, ...) repeat
		// across pages, so namespace them to keep every link pointing at its
		// own page in the combined document.
		slug := kind + "-" + base
		rendered := idRe.ReplaceAllString(buf.String(), `id="`+slug+`--$1"`)
		rendered = hrefRe.ReplaceAllString(rendered, `href="#`+slug+`--$1"`)

		pages = append(pages, page{
			Slug:        slug,
			Title:       title,
			Subcategory: subcat,
			Description: desc,
			Source:      path,
			Body:        rendered,
		})
	}
	return pages, nil
}

func main() {
	// The tool lives in its own throwaway module, so it is run from there and
	// told where the repo root is.
	if len(os.Args) > 1 {
		if err := os.Chdir(os.Args[1]); err != nil {
			panic(err)
		}
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(ghtml.WithUnsafe()),
	)

	resources, err := collect("website/docs/r", "r", md)
	if err != nil {
		panic(err)
	}
	dataSources, err := collect("website/docs/d", "d", md)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	out.WriteString(head)

	writeNav := func(label string, pages []page) {
		fmt.Fprintf(&out, `<div class="nav-group">%s <span class="count">%d</span></div><ul>`, label, len(pages))
		for _, p := range pages {
			fmt.Fprintf(&out, `<li><a href="#%s">%s</a></li>`, p.Slug, html.EscapeString(p.Title))
		}
		out.WriteString("</ul>")
	}

	out.WriteString(`<nav>`)
	out.WriteString(`<div class="nav-title">Advanced Bot Protection</div>`)
	writeNav("Resources", resources)
	writeNav("Data Sources", dataSources)
	out.WriteString(`</nav><main>`)

	writePages := func(label string, pages []page) {
		fmt.Fprintf(&out, `<h1 class="section">%s</h1>`, label)
		for _, p := range pages {
			fmt.Fprintf(&out, `<article id="%s">`, p.Slug)
			fmt.Fprintf(&out,
				`<div class="meta"><div><span>page_title</span>%s</div><div><span>subcategory</span>%s</div><div><span>description</span>%s</div><div><span>source</span><code>%s</code></div></div>`,
				html.EscapeString(p.Title),
				html.EscapeString(p.Subcategory),
				html.EscapeString(p.Description),
				html.EscapeString(p.Source),
			)
			out.WriteString(p.Body)
			out.WriteString(`</article>`)
		}
	}

	writePages("Resources", resources)
	writePages("Data Sources", dataSources)

	out.WriteString(`</main></body></html>`)

	target := "abp-extras/abp-docs-preview.html"
	if err := os.WriteFile(target, out.Bytes(), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("%d resource pages, %d data source pages -> %s\n", len(resources), len(dataSources), target)
}

const head = `<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8">
<title>ABP documentation preview</title>
<style>
:root { --fg:#13141a; --muted:#656a76; --line:#dedfe3; --accent:#7b42bc; --bg-code:#f4f4f6; }
* { box-sizing:border-box }
body { margin:0; display:flex; color:var(--fg); background:#fff;
  font:16px/1.65 -apple-system,BlinkMacSystemFont,"Segoe UI",Helvetica,Arial,sans-serif; }
nav { position:sticky; top:0; align-self:flex-start; width:300px; height:100vh; overflow-y:auto;
  border-right:1px solid var(--line); padding:24px 16px 48px; background:#fafafa; font-size:14px; }
.nav-title { font-weight:600; font-size:15px; margin-bottom:16px; color:var(--accent); }
.nav-group { font-weight:600; margin:18px 0 6px; text-transform:uppercase; font-size:11px;
  letter-spacing:.06em; color:var(--muted); }
.nav-group .count { font-weight:400; }
nav ul { list-style:none; margin:0; padding:0; }
nav li a { display:block; padding:3px 8px; border-radius:4px; color:var(--fg); text-decoration:none;
  font-family:ui-monospace,SFMono-Regular,Menlo,monospace; font-size:12.5px; }
nav li a:hover { background:#ece7f4; color:var(--accent); }
main { flex:1; min-width:0; padding:32px 48px 25vh; max-width:920px; }
h1.section { font-size:13px; text-transform:uppercase; letter-spacing:.1em; color:var(--muted);
  border-bottom:1px solid var(--line); padding-bottom:8px; margin:48px 0 0; }
article { padding-top:40px; margin-top:40px; border-top:1px dashed var(--line); scroll-margin-top:16px; }
article:first-of-type { border-top:none; }
article h1 { font-size:28px; margin:0 0 20px; font-family:ui-monospace,SFMono-Regular,Menlo,monospace; }
article h2 { font-size:21px; margin:36px 0 12px; padding-bottom:6px; border-bottom:1px solid var(--line); }
article h3 { font-size:17px; margin:28px 0 10px; }
.meta { background:#fafafa; border:1px solid var(--line); border-radius:6px; padding:12px 14px;
  margin:0 0 24px; font-size:13px; color:var(--muted); }
.meta div { display:flex; gap:10px; padding:2px 0; }
.meta span { flex:0 0 96px; font-weight:600; text-transform:uppercase; font-size:10.5px;
  letter-spacing:.05em; padding-top:3px; }
code { background:var(--bg-code); padding:.15em .4em; border-radius:3px; font-size:13.5px;
  font-family:ui-monospace,SFMono-Regular,Menlo,monospace; }
pre { background:var(--bg-code); border:1px solid var(--line); border-radius:6px; padding:14px 16px;
  overflow-x:auto; line-height:1.5; }
pre code { background:none; padding:0; font-size:13px; }
ul, ol { padding-left:24px; }
li { margin:5px 0; }
a { color:var(--accent); }
table { border-collapse:collapse; width:100%; margin:16px 0; font-size:14px; }
th, td { border:1px solid var(--line); padding:7px 10px; text-align:left; }
th { background:#fafafa; }
blockquote { margin:16px 0; padding:4px 16px; border-left:3px solid var(--accent); color:var(--muted); }
</style></head><body>
`
EOF

echo "==> Resolving goldmark and rendering..."
(
  cd "$TOOL_DIR"
  GOFLAGS=-mod=mod GOPROXY=off GOSUMDB=off go mod tidy 2>/dev/null ||
    GOFLAGS=-mod=mod GOSUMDB=off go mod tidy
  go run . "$REPO_ROOT"
)

if [ "${OPEN:-0}" = "1" ]; then
  open "$OUT_FILE"
fi
