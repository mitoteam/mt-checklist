package web

import (
	"embed"
	"io/fs"
	"net/http"
)

// embedded web assets
//
//go:embed assets/script.min.js
//go:embed assets/vendor/*.css assets/vendor/*.js
//go:embed assets/webfonts/*
//go:embed assets/css/*.css
//go:embed assets/images/*.png
//go:embed assets/favicon.ico
//go:embed templates/*.html
var embedFS embed.FS

var webAssetsHttpFS http.FileSystem
var templatesFS fs.FS

func init() {
	//prepare FS for subdirectory "/assets"
	webAssetsFS, _ := fs.Sub(embedFS, "assets")
	webAssetsHttpFS = http.FS(webAssetsFS)

	//prepare FS for subdirectory "/templates"
	templatesFS, _ = fs.Sub(embedFS, "templates")
}
