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

// pdfCmd represents the pdfs command
var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "Analyze PDF files",
	RunE:  run.AnalyzePdfs,
}

func init() {
	analyzeCmd.AddCommand(pdfCmd)
	f := pdfCmd.Flags()
	f.String(flags.Model, flags.ModelGeminiProVision, "Model name")
	f.StringSlice(flags.File, nil, "PDF filenames")
	_ = pdfCmd.RegisterFlagCompletionFunc(
		flags.Model,
		func(
			cmd *cobra.Command,
			args []string,
			toComplete string,
		) (
			[]string,
			cobra.ShellCompDirective,
		) {
			return []string{
					flags.ModelGeminiPro,
					flags.ModelGeminiProVision,
					flags.ModelEmbedding001,
				},
				cobra.ShellCompDirectiveDefault
		},
	)
}
