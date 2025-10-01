package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "db-migration-cli",
		Short: "A PostgreSQL database migration CLI tool",
		Long: `db-migration-cli is a command-line interface tool for managing 
PostgreSQL database migrations. It provides functionality to create, 
apply, and rollback database schema changes in a controlled manner.`,
	}
)

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().String("db-host", "localhost", "Database host")
	rootCmd.PersistentFlags().Int("db-port", 5432, "Database port")
	rootCmd.PersistentFlags().String("db-user", "", "Database user")
	rootCmd.PersistentFlags().String("db-password", "", "Database password")
	rootCmd.PersistentFlags().String("db-name", "", "Database name")
	rootCmd.PersistentFlags().String("migrations-dir", "./migrations", "Migrations directory")

	// Bind flags to viper
	viper.BindPFlag("db.host", rootCmd.PersistentFlags().Lookup("db-host"))
	viper.BindPFlag("db.port", rootCmd.PersistentFlags().Lookup("db-port"))
	viper.BindPFlag("db.user", rootCmd.PersistentFlags().Lookup("db-user"))
	viper.BindPFlag("db.password", rootCmd.PersistentFlags().Lookup("db-password"))
	viper.BindPFlag("db.name", rootCmd.PersistentFlags().Lookup("db-name"))
	viper.BindPFlag("migrations.dir", rootCmd.PersistentFlags().Lookup("migrations-dir"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in current directory
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Read environment variables
	viper.SetEnvPrefix("DB_MIGRATION")
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
