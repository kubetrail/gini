package run

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/kubetrail/gini/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

func AnalyzeImages(cmd *cobra.Command, args []string) error {
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
	_ = viper.BindPFlag(flags.Render, cmd.Flag(flags.Render))
	_ = viper.BindPFlag(flags.AllowHarmProbability, cmd.Flag(flags.AllowHarmProbability))
	_ = viper.BindPFlag(flags.File, cmd.Flag(flags.File))
	_ = viper.BindPFlag(flags.Format, cmd.Flag(flags.Format))
	_ = viper.BindEnv(flags.ApiKey, flags.ApiKeyEnv)

	apiKey := viper.GetString(flags.ApiKey)
	modelName := viper.GetString(flags.Model)
	topP := float32(viper.GetFloat64(flags.TopP))
	topK := viper.GetInt32(flags.TopK)
	temperature := float32(viper.GetFloat64(flags.Temperature))
	candidateCount := viper.GetInt32(flags.CandidateCount)
	maxOutputTokens := viper.GetInt32(flags.MaxOutputTokens)
	autoSave := viper.GetBool(flags.AutoSave)
	render := viper.GetString(flags.Render)
	allowHarmProbability := viper.GetString(flags.AllowHarmProbability)
	files := viper.GetStringSlice(flags.File)
	formats := viper.GetStringSlice(flags.Format)

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

	parts := make([]genai.Part, len(files)+1)
	if formats == nil {
		formats = make([]string, len(files))
		for i := range formats {
			formats[i] = flags.FormatJpeg
		}
	} else {
		if len(formats) > len(files) {
			return fmt.Errorf("cannot provide more formats than number of images")
		}
		for i := len(formats); i < len(files); i++ {
			formats = append(formats, flags.FormatJpeg)
		}
	}

	for i, file := range files {
		b, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read image file: %w", err)
		}
		if len(b) > flags.MaxBlobBufferSizeBytes {
			return fmt.Errorf("%s file size needs to be less than %d bytes",
				file, flags.MaxBlobBufferSizeBytes)
		}

		if path.Base(formats[i]) == formats[i] {
			formats[i] = "image/" + formats[i]
		}
		parts[i] = genai.Blob{
			MIMEType: formats[i],
			Data:     b,
		}
	}

	var prompt string
	if len(args) > 0 {
		prompt = strings.Join(args, " ")
	} else {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), ">>> ")

		scanner := bufio.NewScanner(os.Stdin)
		var lines []string
		var hold bool
	InnerLoop:
		for {
			scanner.Scan()
			line := scanner.Text()

			select {
			case <-ctx.Done():
				break InnerLoop
			default:
				if len(line) == 0 && !hold {
					break InnerLoop
				}
				if strings.TrimSpace(line) == startHold {
					hold = true
					continue InnerLoop
				}
				if strings.TrimSpace(line) == endHold && hold {
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

		prompt = strings.Join(lines, "\n")
	}

	parts[len(files)] = genai.Text(prompt)

	send := func(msg string) (*genai.GenerateContentResponse, error) {
		res, err := model.GenerateContent(ctx, parts...)
		if err != nil {
			return nil, fmt.Errorf("failure at backend: %w", err)
		}
		return res, nil
	}

	s := "...sending prompt... please wait"
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\r", s)
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

		if res != nil && res.PromptFeedback != nil {
			for _, rating := range res.PromptFeedback.SafetyRatings {
				if rating.Probability > harmProbability {
					return fmt.Errorf("output harm probability threshold crossed")
				}
			}
		}
	}

	if err := printResponse(res, cmd.OutOrStdout(), render, autoSave, fileWriter); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
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
