package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yeferson59/db-migration-cli/internal/migration"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply pending migrations",
	Long:  `Applies all pending migrations to the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Applying migrations...")

		migrationService := migration.NewService()
		count, err := migrationService.Up()
		if err != nil {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}

		if count == 0 {
			fmt.Println("No pending migrations")
		} else {
			fmt.Printf("✓ Applied %d migration(s)\n", count)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
