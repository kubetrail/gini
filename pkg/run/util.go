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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	startHold = "{{"
	endHold   = "}}"
)

func mdToHTML(md []byte) []byte {
	// create Markdown parser with extensions
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

type persistentFlagValues struct {
	ApiKey               string
	TopP                 float32
	TopK                 int32
	Temperature          float32
	CandidateCount       int32
	MaxOutputTokens      int32
	AutoSave             bool
	Render               string
	AllowHarmProbability string
}

func getPersistentFlags(cmd *cobra.Command) persistentFlagValues {
	pFlags := cmd.Root().PersistentFlags()

	_ = viper.BindPFlag(flags.ApiKey, pFlags.Lookup(flags.ApiKey))
	_ = viper.BindPFlag(flags.TopP, pFlags.Lookup(flags.TopP))
	_ = viper.BindPFlag(flags.TopK, pFlags.Lookup(flags.TopK))
	_ = viper.BindPFlag(flags.Temperature, pFlags.Lookup(flags.Temperature))
	_ = viper.BindPFlag(flags.CandidateCount, pFlags.Lookup(flags.CandidateCount))
	_ = viper.BindPFlag(flags.MaxOutputTokens, pFlags.Lookup(flags.MaxOutputTokens))
	_ = viper.BindPFlag(flags.AutoSave, pFlags.Lookup(flags.AutoSave))
	_ = viper.BindPFlag(flags.Render, pFlags.Lookup(flags.Render))
	_ = viper.BindPFlag(flags.AllowHarmProbability, pFlags.Lookup(flags.AllowHarmProbability))

	_ = viper.BindEnv(flags.ApiKey, flags.ApiKeyEnv)

	apiKey := viper.GetString(flags.ApiKey)
	topP := float32(viper.GetFloat64(flags.TopP))
	topK := viper.GetInt32(flags.TopK)
	temperature := float32(viper.GetFloat64(flags.Temperature))
	candidateCount := viper.GetInt32(flags.CandidateCount)
	maxOutputTokens := viper.GetInt32(flags.MaxOutputTokens)
	autoSave := viper.GetBool(flags.AutoSave)
	render := viper.GetString(flags.Render)
	allowHarmProbability := viper.GetString(flags.AllowHarmProbability)

	return persistentFlagValues{
		ApiKey:               apiKey,
		TopP:                 topP,
		TopK:                 topK,
		Temperature:          temperature,
		CandidateCount:       candidateCount,
		MaxOutputTokens:      maxOutputTokens,
		AutoSave:             autoSave,
		Render:               render,
		AllowHarmProbability: allowHarmProbability,
	}
}
