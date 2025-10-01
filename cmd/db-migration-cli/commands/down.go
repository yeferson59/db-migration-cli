package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yeferson59/db-migration-cli/internal/migration"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback last migration",
	Long:  `Rolls back the most recently applied migration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Rolling back migration...")

		migrationService := migration.NewService()
		if err := migrationService.Down(); err != nil {
			return fmt.Errorf("failed to rollback migration: %w", err)
		}

		fmt.Println("✓ Rolled back 1 migration")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
