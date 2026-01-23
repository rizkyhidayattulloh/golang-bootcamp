package docs

import "embed"

//go:embed swagger.json swagger-ui.html
var EmbedAssets embed.FS
