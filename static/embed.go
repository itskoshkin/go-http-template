package static

import (
	"embed"
)

//go:embed templates/*.gohtml
var TemplatesFS embed.FS

//go:embed styles scripts icons fonts
var PublicFS embed.FS
