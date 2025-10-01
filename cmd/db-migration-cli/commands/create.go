package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yeferson59/db-migration-cli/internal/migration"
)

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration file",
	Long:  `Creates a new migration file with timestamp prefix in the migrations directory.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		migrationService := migration.NewService()
		filename, err := migrationService.Create(name)
		if err != nil {
			return fmt.Errorf("failed to create migration: %w", err)
		}

		fmt.Printf("✓ Created migration: %s\n", filename)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
