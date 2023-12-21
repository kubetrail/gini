package run

import (
	"bufio"
	"fmt"
	"io"

	termmarkdown "github.com/MichaelMure/go-term-markdown"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/google/generative-ai-go/genai"
	"github.com/kubetrail/gini/pkg/flags"
)

const (
	startHold = "{{"
	endHold   = "}}"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func mdToPretty(md []byte) []byte {
	return termmarkdown.Render(string(md), 80, 6)
}

func mdToMd(md []byte) []byte {
	return md
}

func printResponse(resp *genai.GenerateContentResponse, w io.Writer, render string, autoSave bool, fileWriter *bufio.Writer) error {
	if autoSave {
		if _, err := fileWriter.WriteString(fmt.Sprintf("%s\n", "[response]>>>")); err != nil {
			return fmt.Errorf("failed to write to history file: %w", err)
		}
	}

	var renderFunc func([]byte) []byte
	switch render {
	case flags.RenderFormatHtml:
		renderFunc = mdToHTML
	case flags.RenderFormatMarkdown:
		renderFunc = mdToMd
	case flags.RenderFormatPretty:
		renderFunc = mdToPretty
	default:
		return fmt.Errorf("invalid render format: %s", render)
	}

	var result []byte
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				text, ok := part.(genai.Text)
				if ok {
					result = mdToPretty([]byte(text))
					if _, err := fmt.Fprintln(w, string(result)); err != nil {
						return fmt.Errorf("failed to write to output: %w", err)
					}

					result = renderFunc([]byte(text))
					if autoSave {
						if _, err := fileWriter.WriteString(fmt.Sprintf("%s\n", result)); err != nil {
							return fmt.Errorf("failed to write to history file: %w", err)
						}
					}
				} else {
					if _, err := fmt.Fprintln(w, part); err != nil {
						return fmt.Errorf("failed to write to output: %w", err)
					}
					if autoSave {
						if _, err := fileWriter.WriteString(fmt.Sprintf("%s\n", part)); err != nil {
							return fmt.Errorf("failed to write to history file: %w", err)
						}
					}
				}
			}
		}
	}

	return nil
}
