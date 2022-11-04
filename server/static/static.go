package static

import "embed"

// Assets represents the embedded files.
//
//go:embed *.tmpl pages/*.tmpl
var TemplateAssets embed.FS

//go:embed *.css *.svg
var StyleAssets embed.FS
