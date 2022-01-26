package assets

import (
	"embed"
	_ "embed"
)

//go:embed templates
var TemplatesFs embed.FS

//go:embed build.txt
var BuildInfo string

//go:embed buildver.txt
var BuildVer string
