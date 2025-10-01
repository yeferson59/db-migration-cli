# db-migration-cli

A command-line interface tool for managing PostgreSQL database migrations. Built with Go and following clean architecture principles.

## Features

- ✨ Create timestamped migration files
- 🚀 Apply migrations (up)
- ⏪ Rollback migrations (down)
- 📊 View migration status
- 🔧 Configurable via YAML or environment variables
- 📁 Clean project structure with separation of concerns

## Project Structure

```
db-migration-cli/
├── cmd/
│   └── db-migration-cli/
│       └── commands/          # CLI commands (cobra)
│           ├── root.go        # Root command
│           ├── init.go        # Initialize migration environment
│           ├── create.go      # Create new migration
│           ├── up.go          # Apply migrations
│           ├── down.go        # Rollback migrations
│           └── status.go      # Show migration status
├── internal/
│   ├── config/                # Configuration management
│   │   └── config.go
│   ├── database/              # Database connection layer
│   │   └── database.go
│   └── migration/             # Core migration logic
│       ├── service.go
│       └── parser.go
├── pkg/
│   └── utils/                 # Utility functions
│       └── fileutil.go
├── migrations/                # Migration files directory
├── main.go                    # Application entry point
├── Makefile                   # Build and development tasks
├── config.example.yaml        # Example configuration file
└── README.md
```

## Architecture

The project follows clean architecture principles:

- **cmd/**: Entry point and CLI command definitions
- **internal/**: Private application code
  - **config/**: Configuration loading and management
  - **database/**: Database connection abstraction
  - **migration/**: Business logic for migrations
- **pkg/**: Public, reusable utility packages

## Installation

### Prerequisites

- Go 1.24+ installed
- PostgreSQL database

### Build from source

```bash
# Clone the repository
git clone https://github.com/yeferson59/db-migration-cli.git
cd db-migration-cli

# Install dependencies
make install

# Build the binary
make build

# The binary will be in ./bin/db-migration-cli
```

## Configuration

### Option 1: Configuration File

Create a `config.yaml` file in your project root:

```yaml
db:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: mydb
  sslmode: disable

migrations:
  dir: ./migrations
```

### Option 2: Command-line Flags

```bash
./bin/db-migration-cli --db-host=localhost --db-port=5432 --db-user=postgres --db-password=postgres --db-name=mydb
```

### Option 3: Environment Variables

```bash
export DB_MIGRATION_DB_HOST=localhost
export DB_MIGRATION_DB_PORT=5432
export DB_MIGRATION_DB_USER=postgres
export DB_MIGRATION_DB_PASSWORD=postgres
export DB_MIGRATION_DB_NAME=mydb
```

## Usage

### Initialize Migration Environment

```bash
./bin/db-migration-cli init
```

This creates:
- The migrations directory
- The `schema_migrations` table in your database

### Create a New Migration

```bash
./bin/db-migration-cli create create_users_table
```

This generates two files in the migrations directory:
- `20240101120000_create_users_table.up.sql`
- `20240101120000_create_users_table.down.sql`

Edit these files with your SQL statements.

### Apply Pending Migrations

```bash
./bin/db-migration-cli up
```

### Rollback Last Migration

```bash
./bin/db-migration-cli down
```

### Check Migration Status

```bash
./bin/db-migration-cli status
```

## Development

### Available Make Targets

```bash
make help          # Show available commands
make build         # Build the binary
make clean         # Clean build files
make install       # Install dependencies
make test          # Run tests
make fmt           # Format code
make vet           # Run go vet
make lint          # Run linter
```

### Running Tests

```bash
make test
```

### Code Formatting

```bash
make fmt
```

## Dependencies

- **[cobra](https://github.com/spf13/cobra)**: CLI framework
- **[viper](https://github.com/spf13/viper)**: Configuration management
- **[pgx](https://github.com/jackc/pgx)**: PostgreSQL driver and toolkit

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License
