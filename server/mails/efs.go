package mails

import (
	"embed"
)

//go:embed *.tmpl
var EmbeddedFiles embed.FS
