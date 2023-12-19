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

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Analyze images",
	RunE:  run.AnalyzeImages,
}

func init() {
	analyzeCmd.AddCommand(imagesCmd)
	f := imagesCmd.Flags()
	f.String(flags.Model, flags.ModelGeminiProVision, "Model name")
	f.Float32(flags.TopP, -1, "Model TopP value (-1 means do not configure)")
	f.Int32(flags.TopK, -1, "Model TopK value (-1 means do not configure)")
	f.Float32(flags.Temperature, -1, "Model temperature (-1 means do not configure)")
	f.Int32(flags.CandidateCount, -1, "Model candidate count (-1 means do not configure)")
	f.Int32(flags.MaxOutputTokens, -1, "Model max output tokens (-1 means do not configure)")
	f.StringSlice(flags.ImageFiles, nil, "Image filenames (input multiple names using comma or repeated use of the flag)")
	f.StringSlice(flags.Formats, nil, "Image formats (should be either empty, i.e., assumed jpeg, or of same length as image files)")
	_ = imagesCmd.RegisterFlagCompletionFunc(
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

	_ = imagesCmd.RegisterFlagCompletionFunc(
		flags.Formats,
		func(
			cmd *cobra.Command,
			args []string,
			toComplete string,
		) (
			[]string,
			cobra.ShellCompDirective,
		) {
			return []string{
					flags.FormatJpeg,
					flags.FormatPng,
					flags.FormatHeif,
					flags.FormatHeic,
					flags.FormatWebp,
				},
				cobra.ShellCompDirectiveDefault
		},
	)
}
