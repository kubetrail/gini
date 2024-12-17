/*
Copyright © 2023 kubetrail.io authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/kubetrail/gini/pkg/flags"
	"github.com/kubetrail/gini/pkg/run"
	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start chat",
	Long: `
Start an interactive chat with Google Gemini model using this command.

Hit enter twice to send your prompt.

Alternatively, if you have blank lines in your text prompt then
enclose your prompt within {{ and }} as shown below
{{
This is a prompt with blank lines


The prompt ends here
}}
`,
	RunE: run.Chat,
}

func init() {
	rootCmd.AddCommand(chatCmd)
	f := chatCmd.Flags()
	f.String(flags.Model, flags.M09, "Model name")
	f.StringSlice(flags.File, nil, "Image filenames")
	f.StringSlice(flags.Format, nil, "Image formats (assumes application/pdf when unspecified)")
	_ = chatCmd.RegisterFlagCompletionFunc(
		flags.Model,
		func(
			cmd *cobra.Command,
			args []string,
			toComplete string,
		) (
			[]string,
			cobra.ShellCompDirective,
		) {
			return flags.Models, cobra.ShellCompDirectiveDefault
		},
	)

	_ = chatCmd.RegisterFlagCompletionFunc(
		flags.Format,
		func(
			cmd *cobra.Command,
			args []string,
			toComplete string,
		) (
			[]string,
			cobra.ShellCompDirective,
		) {
			return []string{
					flags.FormatPdf,
					flags.FormatText,
					flags.FormatJpeg,
				},
				cobra.ShellCompDirectiveDefault
		},
	)
}
