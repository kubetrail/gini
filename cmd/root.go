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
	"fmt"
	"os"

	"github.com/kubetrail/gem/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gem",
	Short: "CLI to interact with Google Gemini AI model",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gem.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	f := rootCmd.PersistentFlags()
	f.String(flags.ApiKey, "", fmt.Sprintf("API Key (Env. %s)", flags.ApiKeyEnv))
	f.String(flags.Model, flags.ModelGeminiPro, "Model name")
	f.Bool(flags.AutoSave, false, "Auto save chat history")
	f.String(flags.AllowHarmProbability, flags.HarmProbabilityNegligible, "Allow harm probability at or below specified level")

	_ = rootCmd.RegisterFlagCompletionFunc(
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

	_ = rootCmd.RegisterFlagCompletionFunc(
		flags.AllowHarmProbability,
		func(
			cmd *cobra.Command,
			args []string,
			toComplete string,
		) (
			[]string,
			cobra.ShellCompDirective,
		) {
			return []string{
					flags.HarmProbabilityUnspecified,
					flags.HarmProbabilityNegligible,
					flags.HarmProbabilityLow,
					flags.HarmProbabilityMedium,
					flags.HarmProbabilityHigh,
				},
				cobra.ShellCompDirectiveDefault
		},
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gem" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gem")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
