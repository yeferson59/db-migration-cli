package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yeferson59/db-migration-cli/internal/migration"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize migration environment",
	Long:  `Creates the migrations directory and sets up the schema_migrations table in the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Initializing migration environment...")

		migrationService := migration.NewService()
		if err := migrationService.Init(); err != nil {
			return fmt.Errorf("failed to initialize: %w", err)
		}

		fmt.Println("✓ Migration environment initialized successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
