package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/config"
	"github.com/Fire-Dragon-DoL/clickup-cli/internal/keyring"
	"github.com/Fire-Dragon-DoL/clickup-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	cfgFile      string
	spaceID      string
	outputFormat string
	strictResolve bool

	cfg       *config.Config
	kr        *keyring.Keyring
	formatter *output.Formatter
)

var rootCmd = &cobra.Command{
	Use:   "clickup",
	Short: "CLI for ClickUp",
	Long: `A command-line interface for interacting with ClickUp tasks and projects.

Configure your ClickUp space and API key to get started.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cfgFile != "" {
			cfg = config.LoadFromFile(cfgFile)
		} else {
			defaultPath := defaultConfigPath()
			if _, err := os.Stat(defaultPath); err == nil {
				cfg = config.LoadFromFile(defaultPath)
			} else {
				cfg = config.Load()
			}
		}

		cfg.ApplyCLIOverrides(spaceID, outputFormat, strictResolve)
		formatter = output.NewFormatter(cfg.OutputFormat)
		kr = keyring.New(keyring.NewSystemProvider())

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	rootCmd.PersistentFlags().StringVar(&spaceID, "space", "", "ClickUp space ID")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format (text|json)")
	rootCmd.PersistentFlags().BoolVar(&strictResolve, "strict", false, "fail on ambiguous name resolution")
}

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "clickup", "config.json")
}

func GetConfig() *config.Config {
	return cfg
}

func GetKeyring() *keyring.Keyring {
	return kr
}

func GetFormatter() *output.Formatter {
	return formatter
}

func PrintOutput(data any) error {
	out, err := formatter.Format(data)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}
