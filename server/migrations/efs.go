package migrations

import (
	"embed"
)

//go:embed *.sql
var EmbeddedFiles embed.FS
