package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yeferson59/db-migration-cli/internal/migration"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long:  `Displays the status of all migrations (applied or pending).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		migrationService := migration.NewService()
		status, err := migrationService.Status()
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}

		fmt.Println("Migration Status:")
		fmt.Println("================")
		for _, s := range status {
			fmt.Println(s)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
