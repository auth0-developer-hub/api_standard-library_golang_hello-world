package embed

import (
	"embed"
	"io/fs"
)

//go:embed swagger
var publicHTML embed.FS

func PublicHTMLFS() fs.FS {
	publicHTMLfs, _ := fs.Sub(publicHTML, "swagger")
	return publicHTMLfs
}
