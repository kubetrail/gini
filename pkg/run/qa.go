package run

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/generative-ai-go/genai"
	"github.com/kubetrail/gem/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

func Qa(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		select {
		case <-ctx.Done():
			os.Exit(0)
		}
	}()

	_ = viper.BindPFlag(flags.ApiKey, cmd.Flag(flags.ApiKey))
	_ = viper.BindPFlag(flags.Model, cmd.Flag(flags.Model))
	_ = viper.BindEnv(flags.ApiKey, flags.ApiKeyEnv)

	apiKey := viper.GetString(flags.ApiKey)
	modelName := viper.GetString(flags.Model)

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return fmt.Errorf("failed to create new genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	cs := model.StartChat()

	send := func(msg string) (*genai.GenerateContentResponse, error) {
		res, err := cs.SendMessage(ctx, genai.Text(msg))
		if err != nil {
			return nil, fmt.Errorf("failed to send message: %w", err)
		}
		return res, nil
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "please type query below\npress enter twice to send query\nquit or exit to end session\n")

OuterLoop:
	for {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "> ")

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
				if len(line) == 0 ||
					strings.ToLower(line) == "quit" ||
					strings.ToLower(line) == "exit" {
					break InnerLoop
				}
				lines = append(lines, line)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		if len(lines) == 0 {
			return nil
		}

		select {
		case <-ctx.Done():
			break OuterLoop
		default:
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "sending...")
			res, err := send(strings.Join(lines, "\n"))
			if err != nil {
				return err
			}
			printResponse(res)
		}
	}

	return nil
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("")
}
