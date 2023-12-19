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

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/kubetrail/gem/pkg/flags"
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
	_ = viper.BindPFlag(flags.AutoSave, cmd.Flag(flags.AutoSave))
	_ = viper.BindEnv(flags.ApiKey, flags.ApiKeyEnv)

	apiKey := viper.GetString(flags.ApiKey)
	modelName := viper.GetString(flags.Model)
	autoSave := viper.GetBool(flags.AutoSave)

	fileName := fmt.Sprintf("history-%s.txt", uuid.New().String())
	var fileWriter *bufio.Writer

	if len(apiKey) == 0 || len(modelName) == 0 {
		return fmt.Errorf("api-key or model cannot be empty")
	}

	if autoSave {
		f, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("failed to create history file: %w", err)
		}
		defer f.Close()

		fileWriter = bufio.NewWriter(f)

		if _, err := fileWriter.WriteString(fmt.Sprintf("%s\n", time.Now().String())); err != nil {
			return fmt.Errorf("failed to write to history file: %w", err)
		}
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return fmt.Errorf("failed to create new genai client: %w", err)
	}
	defer client.Close()

	modelIterator := client.ListModels(ctx)
	for {
		m, err := modelIterator.Next()
		if err != nil {
			break
		}
		fmt.Println("***", m.Name, m.DisplayName)
	}

	model := client.GenerativeModel(modelName)
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
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
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

	return nil
}
