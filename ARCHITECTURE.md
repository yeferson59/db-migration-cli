# Architecture Documentation

## Overview

The `db-migration-cli` follows clean architecture principles with clear separation of concerns, making the codebase maintainable, testable, and scalable.

## Architecture Layers

### 1. Presentation Layer (CLI)
- **Location**: `cmd/db-migration-cli/commands/`
- **Responsibility**: Handles user interaction through CLI commands
- **Components**:
  - `root.go`: Root command and global configuration
  - `init.go`: Initialize migration environment
  - `create.go`: Create new migration files
  - `up.go`: Apply pending migrations
  - `down.go`: Rollback migrations
  - `status.go`: Show migration status

### 2. Application Layer
- **Location**: `internal/migration/`
- **Responsibility**: Core business logic for migrations
- **Components**:
  - `service.go`: Migration service with all operations
  - `parser.go`: Migration file parsing utilities

### 3. Infrastructure Layer

#### Configuration
- **Location**: `internal/config/`
- **Responsibility**: Load and manage application configuration
- **Features**:
  - YAML file support
  - Environment variables
  - Command-line flags
  - Connection string generation

#### Database
- **Location**: `internal/database/`
- **Responsibility**: Database connection and operations
- **Features**:
  - PostgreSQL connection pooling (via pgx)
  - Query execution
  - Transaction management

### 4. Utility Layer
- **Location**: `pkg/utils/`
- **Responsibility**: Reusable utility functions
- **Current utilities**:
  - File system operations

## Design Patterns

### 1. Service Pattern
The migration service (`internal/migration/service.go`) encapsulates all business logic related to migrations, providing a clean API for the CLI commands.

### 2. Repository Pattern
The database layer (`internal/database/database.go`) abstracts database operations, making it easy to test and potentially swap implementations.

### 3. Configuration Management
Using Viper for flexible configuration management:
- File-based (YAML)
- Environment variables
- Command-line flags
- Configuration precedence

## Data Flow

```
User Input (CLI)
    ↓
Commands Layer (cmd/)
    ↓
Service Layer (internal/migration/)
    ↓
Database Layer (internal/database/)
    ↓
PostgreSQL Database
```

## Dependencies

### Core Dependencies
- **cobra** (v1.10.1): CLI framework
- **viper** (v1.21.0): Configuration management
- **pgx/v5** (v5.7.6): PostgreSQL driver and toolkit

### Why These Dependencies?

1. **Cobra**: Industry-standard CLI framework with:
   - Automatic help generation
   - Flag parsing
   - Subcommand support
   - Shell completion

2. **Viper**: Flexible configuration with:
   - Multiple format support
   - Environment variable binding
   - Live configuration reloading
   - Hierarchical configuration

3. **pgx**: Modern PostgreSQL driver with:
   - Better performance than database/sql
   - Native PostgreSQL types
   - Connection pooling
   - Context support

## Project Structure

```
db-migration-cli/
├── cmd/                        # CLI commands (Presentation)
│   └── db-migration-cli/
│       └── commands/
├── internal/                   # Private application code
│   ├── config/                # Configuration management
│   ├── database/              # Database layer
│   └── migration/             # Business logic
├── pkg/                       # Public utilities
│   └── utils/
├── migrations/                # Migration files
├── main.go                    # Application entry point
├── Makefile                   # Build automation
├── config.example.yaml        # Example configuration
└── README.md                  # Documentation
```

## Key Design Decisions

### 1. Clean Architecture
- **Benefit**: Easy to test, maintain, and extend
- **Implementation**: Clear layer separation with dependency injection

### 2. Internal vs Pkg
- **internal/**: Private packages, implementation details
- **pkg/**: Public, reusable utilities
- **Benefit**: Clear API boundaries

### 3. Configuration Flexibility
- Support for multiple configuration sources
- **Benefit**: Works in different environments (dev, staging, prod)

### 4. Migration File Naming
- Format: `{timestamp}_{description}.{up|down}.sql`
- **Benefit**: Chronological ordering, clear purpose

## Testing Strategy

The architecture supports multiple testing levels:

1. **Unit Tests**: Test individual functions in isolation
2. **Integration Tests**: Test database operations
3. **E2E Tests**: Test complete CLI commands

## Future Extensibility

The architecture is designed to easily add:

1. **Multiple Database Support**: Add drivers for MySQL, SQLite, etc.
2. **Migration Rollback Strategies**: Multiple rollback options
3. **Migration Versioning**: Advanced version control
4. **Cloud Database Support**: Specific cloud provider features
5. **Web UI**: Add a web interface layer
6. **API Server**: Expose functionality via REST/GraphQL API

## Performance Considerations

1. **Connection Pooling**: Using pgxpool for efficient connections
2. **Batch Operations**: Support for applying multiple migrations
3. **Lazy Loading**: Load configurations and connections only when needed

## Security

1. **Configuration**: Sensitive data (passwords) can be set via environment variables
2. **SQL Injection**: Using parameterized queries via pgx
3. **File Permissions**: Migration files have appropriate permissions (0644)

## Monitoring & Logging

Current implementation includes:
- Command execution feedback
- Error reporting with context
- Migration status tracking in database

Future enhancements:
- Structured logging (e.g., zap, zerolog)
- Metrics collection
- Audit trails
