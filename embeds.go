package mapboxdemo

import (
	"embed"
)

//go:embed templates/*.gohtml
var TemplateFS embed.FS
