package markdown

import (
	"bytes"
	"fmt"
	chroma "github.com/alecthomas/chroma/formatters/html"
	"github.com/worldiety/wikigrabber/internal/markdown/hugo"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TransformedProject struct {
	Project Project
	Pages   []Page
}

type Page struct {
	Markdown string
	Html     string
	Plain    string
	HtmlFile string
	SrcFile  string
	Config   *PageConfig
}

func Transform(cfg Config) ([]TransformedProject, error) {
	var res []TransformedProject
	projects, err := cfg.Collect()
	if err != nil {
		return nil, err
	}

	r := goldmark.DefaultRenderer()

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.New(meta.WithTable()),
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chroma.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(util.Prioritized(&hugo.TocTransformer{
				R: r,
			}, 10)),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			renderer.WithNodeRenderers(
				util.Prioritized(extension.NewTableHTMLRenderer(), 500),
			),
		),
	)

	for _, project := range projects {
		tp := TransformedProject{
			Project: project,
			Pages:   nil,
		}

		for _, mdFile := range project.Markdown {

			buf, err := ioutil.ReadFile(filepath.Join(project.Path, mdFile))
			if err != nil {
				return nil, err
			}

			relFname := filepath.Join(project.Name, mdFile+".html")
			fname := filepath.Join(cfg.OutDir, relFname)

			log.Printf("%s ->\n%s\n", mdFile, fname)

			if err := os.MkdirAll(filepath.Dir(fname), os.ModePerm); err != nil {
				return nil, err
			}

			buf, pageCfg := stripAndParseFrontmatter(buf)
			if pageCfg != nil {
				log.Printf("%+v\n", pageCfg)
			}

			out := &bytes.Buffer{}

			ctx := parser.NewContext()
			md.Parser().Parse(text.NewReader(buf), parser.WithContext(ctx))

			if err := md.Convert(buf, out); err != nil {
				return nil, fmt.Errorf("failed to parse markdown from '%s': %w", mdFile, err)
			}

			tmpStr := out.String()
			myToc := hugo.TableOfContents(ctx)
			if len(myToc.Headers) > 0 {
				myTocHtml := myToc.ToHTML(0, 10, false)
				tmpStr = strings.ReplaceAll(tmpStr, "[[<em>TOC</em>]]", myTocHtml)

			}

			if err = ioutil.WriteFile(fname, []byte(tmpStr), os.ModePerm); err != nil {
				return nil, err
			}

			tp.Pages = append(tp.Pages, Page{
				Markdown: string(buf),
				Html:     tmpStr,
				Plain:    string(buf),
				HtmlFile: fname,
				SrcFile:  relFname,
				Config:   pageCfg,
			})

		}

		res = append(res, tp)
	}

	return res, nil
}

func stripAndParseFrontmatter(in []byte) ([]byte, *PageConfig) {
	dst := &PageConfig{}

	start, end := Frontmatter(in)
	if start == -1 {
		return in, nil
	}

	if err := yaml.Unmarshal(in[start:end], dst); err != nil {
		fmt.Println(string(in))
		log.Printf("failed to parse frontmatter: %v", err)
	}

	return in[end:], dst
}
