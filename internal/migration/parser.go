package migration

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var migrationFileRegex = regexp.MustCompile(`^(\d{14})_(.+)\.(up|down)\.sql$`)

// parseMigrationFilename parses a migration filename
// Returns: version, name, isUp, error
func parseMigrationFilename(filename string) (string, string, bool, error) {
	base := filepath.Base(filename)

	matches := migrationFileRegex.FindStringSubmatch(base)
	if matches == nil {
		return "", "", false, fmt.Errorf("invalid migration filename: %s", filename)
	}

	version := matches[1]
	name := matches[2]
	direction := matches[3]

	// Clean up name (replace underscores with spaces for display)
	name = strings.ReplaceAll(name, "_", " ")

	isUp := direction == "up"

	return version, name, isUp, nil
}
