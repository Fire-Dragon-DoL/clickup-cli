package cmd

import (
	"fmt"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/api"
	"github.com/Fire-Dragon-DoL/clickup-cli/internal/resolver"
	"github.com/spf13/cobra"
)

var listsCmd = &cobra.Command{
	Use:   "lists",
	Short: "Manage lists",
}

var listsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List lists in a folder",
	RunE: func(cmd *cobra.Command, args []string) error {
		folderArg, err := cmd.Flags().GetString("folder")
		if folderArg == "" {
			return fmt.Errorf("--folder flag is required")
		}

		kr := GetKeyring()
		apiKey, err := kr.GetAPIKey()
		if err != nil {
			return err
		}

		cfg := GetConfig()
		client := api.NewClient(apiKey, cfg.BaseURL, cfg.SpaceID)
		res := resolver.New(client, cfg.StrictResolve)

		folderID, err := res.ResolveFolder(folderArg)
		if err != nil {
			return err
		}

		lists, err := api.GetLists(client, folderID)
		if err != nil {
			return err
		}

		return PrintOutput(lists)
	},
}

func init() {
	rootCmd.AddCommand(listsCmd)
	listsCmd.AddCommand(listsListCmd)
	listsListCmd.Flags().StringP("folder", "f", "", "folder name, ID, or URL")
}
