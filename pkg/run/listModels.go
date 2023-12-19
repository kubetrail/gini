package run

import (
	"errors"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/google/generative-ai-go/genai"
	"github.com/kubetrail/gem/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"
)

func ListModels(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_ = viper.BindPFlag(flags.ApiKey, cmd.Flag(flags.ApiKey))
	_ = viper.BindEnv(flags.ApiKey, flags.ApiKeyEnv)

	apiKey := viper.GetString(flags.ApiKey)

	if len(apiKey) == 0 {
		return fmt.Errorf("api-key cannot be empty")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return fmt.Errorf("failed to create new genai client: %w", err)
	}
	defer client.Close()

	modelIterator := client.ListModels(ctx)
	for {
		model, err := modelIterator.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterator over model list: %w", err)
		}

		b, err := yaml.Marshal(model)
		if err != nil {
			return fmt.Errorf("failed to serialize model info: %w", err)
		}

		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", b); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	}

	return nil
}
