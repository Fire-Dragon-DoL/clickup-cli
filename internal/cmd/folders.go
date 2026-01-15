package cmd

import (
	"fmt"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/api"
	"github.com/spf13/cobra"
)

var foldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "Manage folders",
}

var foldersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List folders in the current space",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig()
		if cfg.SpaceID == "" {
			PrintError(fmt.Errorf("space ID not configured"))
			return fmt.Errorf("space ID is required")
		}

		kr := GetKeyring()
		apiKey, err := kr.GetAPIKey()
		if err != nil {
			PrintError(err)
			return err
		}

		client := api.NewClient(apiKey, "")
		folders, err := api.GetFolders(client, cfg.SpaceID)
		if err != nil {
			PrintError(err)
			return err
		}

		return PrintOutput(folders)
	},
}

func init() {
	rootCmd.AddCommand(foldersCmd)
	foldersCmd.AddCommand(foldersListCmd)
}
