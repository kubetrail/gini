package run

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/kubetrail/gini/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

func Chat(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_ = viper.BindPFlag(flags.ApiKey, cmd.Flag(flags.ApiKey))
	_ = viper.BindPFlag(flags.Model, cmd.Flag(flags.Model))
	_ = viper.BindPFlag(flags.TopP, cmd.Flag(flags.TopP))
	_ = viper.BindPFlag(flags.TopK, cmd.Flag(flags.TopK))
	_ = viper.BindPFlag(flags.Temperature, cmd.Flag(flags.Temperature))
	_ = viper.BindPFlag(flags.CandidateCount, cmd.Flag(flags.CandidateCount))
	_ = viper.BindPFlag(flags.MaxOutputTokens, cmd.Flag(flags.MaxOutputTokens))
	_ = viper.BindPFlag(flags.AutoSave, cmd.Flag(flags.AutoSave))
	_ = viper.BindPFlag(flags.AllowHarmProbability, cmd.Flag(flags.AllowHarmProbability))
	_ = viper.BindEnv(flags.ApiKey, flags.ApiKeyEnv)

	apiKey := viper.GetString(flags.ApiKey)
	modelName := viper.GetString(flags.Model)
	topP := float32(viper.GetFloat64(flags.TopP))
	topK := viper.GetInt32(flags.TopK)
	temperature := float32(viper.GetFloat64(flags.Temperature))
	candidateCount := viper.GetInt32(flags.CandidateCount)
	maxOutputTokens := viper.GetInt32(flags.MaxOutputTokens)
	autoSave := viper.GetBool(flags.AutoSave)
	allowHarmProbability := viper.GetString(flags.AllowHarmProbability)

	if len(apiKey) == 0 || len(modelName) == 0 {
		return fmt.Errorf("api-key or model cannot be empty")
	}

	fileName := fmt.Sprintf("history-%s.txt", uuid.New().String())
	var fileWriter *bufio.Writer
	if autoSave {
		f, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("failed to create history file: %w", err)
		}
		defer f.Close()

		fileWriter = bufio.NewWriter(f)

		if _, err := fileWriter.WriteString(
			fmt.Sprintf("Command: %s\nTimestamp: %s\n",
				strings.Join(os.Args, " "),
				time.Now().String(),
			),
		); err != nil {
			return fmt.Errorf("failed to write to history file: %w", err)
		}
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return fmt.Errorf("failed to create new genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	if topP >= 0 {
		model.SetTopP(topP)
	}
	if topK >= 0 {
		model.SetTopK(topK)
	}
	if temperature >= 0 {
		model.SetTemperature(temperature)
	}
	if candidateCount >= 0 {
		model.SetCandidateCount(candidateCount)
	}
	if maxOutputTokens >= 0 {
		model.SetMaxOutputTokens(maxOutputTokens)
	}
	cs := model.StartChat()

	send := func(msg string) (*genai.GenerateContentResponse, error) {
		res, err := cs.SendMessage(ctx, genai.Text(msg))
		if err != nil {
			return nil, fmt.Errorf("failed to send message: %w", err)
		}
		return res, nil
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "please type prompt below\npress enter twice to send prompt\njust enter to quit\n")

OuterLoop:
	for i := 0; ; i++ {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), fmt.Sprintf("[%d]>>> ", i+1))

		scanner := bufio.NewScanner(os.Stdin)
		var lines []string
	InnerLoop:
		for {
			scanner.Scan()
			line := scanner.Text()

			select {
			case <-ctx.Done():
				break InnerLoop
			default:
				if len(line) == 0 {
					break InnerLoop
				}
				lines = append(lines, line)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		if len(lines) == 0 {
			break OuterLoop
		}

		select {
		case <-ctx.Done():
			break OuterLoop
		default:
			s := "... sending prompt... please wait"
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\r", s)
			prompt := strings.Join(lines, "\n")
			res, err := send(prompt)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\r", strings.Repeat(" ", len(s)+2))

			if allowHarmProbability != flags.HarmProbabilityUnspecified {
				var harmProbability genai.HarmProbability
				switch allowHarmProbability {
				case flags.HarmProbabilityNegligible:
					harmProbability = genai.HarmProbabilityNegligible
				case flags.HarmProbabilityLow:
					harmProbability = genai.HarmProbabilityLow
				case flags.HarmProbabilityMedium:
					harmProbability = genai.HarmProbabilityMedium
				case flags.HarmProbabilityHigh:
					harmProbability = genai.HarmProbabilityHigh
				default:
					return fmt.Errorf("invalid harm probability:%s", allowHarmProbability)
				}

				for _, rating := range res.PromptFeedback.SafetyRatings {
					if rating.Probability > harmProbability {
						return fmt.Errorf("output harm probability threshold crossed")
					}
				}
			}

			if autoSave {
				if _, err := fileWriter.WriteString(fmt.Sprintf("[%d]>>> %s\n", i+1, prompt)); err != nil {
					return fmt.Errorf("failed to write to history file: %w", err)
				}
			}

			if err := printResponse(res, cmd.OutOrStdout(), autoSave, fileWriter); err != nil {
				return fmt.Errorf("failed to write response: %w", err)
			}
		}
	}

	if autoSave {
		if err := fileWriter.Flush(); err != nil {
			return fmt.Errorf("failed to write to history file: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("history saved to %s", fileName)); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	}

	return nil
}

func printResponse(resp *genai.GenerateContentResponse, w io.Writer, autoSave bool, fileWriter *bufio.Writer) error {
	if autoSave {
		if _, err := fileWriter.WriteString(fmt.Sprintf("%s\n", "[response]>>>")); err != nil {
			return fmt.Errorf("failed to write to history file: %w", err)
		}
	}
	var result []byte
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				text, ok := part.(genai.Text)
				if ok {
					result = markdown.Render(string(text), 80, 6)
					if _, err := fmt.Fprintln(w, string(result)); err != nil {
						return fmt.Errorf("failed to write to output: %w", err)
					}
				} else {
					if _, err := fmt.Fprintln(w, part); err != nil {
						return fmt.Errorf("failed to write to output: %w", err)
					}
				}
				if autoSave {
					if _, err := fileWriter.WriteString(fmt.Sprintf("%s\n", part)); err != nil {
						return fmt.Errorf("failed to write to history file: %w", err)
					}
				}
			}
		}
	}

	return nil
}
