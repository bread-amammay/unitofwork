package postgres

import "embed"

//go:embed migrations/*.sql
var EmbedMigrations embed.FS
