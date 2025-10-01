# Migrations Directory

This directory contains database migration files.

## Migration File Naming Convention

Migration files should follow this naming pattern:
```
{timestamp}_{description}.{up|down}.sql
```

For example:
- `20240101120000_create_users_table.up.sql`
- `20240101120000_create_users_table.down.sql`

## Creating Migrations

Use the CLI to create new migrations:
```bash
db-migration-cli create <migration_name>
```

This will automatically generate both `.up.sql` and `.down.sql` files with the correct timestamp prefix.

## Migration Files

- **`.up.sql`**: Contains SQL statements to apply the migration
- **`.down.sql`**: Contains SQL statements to rollback the migration
