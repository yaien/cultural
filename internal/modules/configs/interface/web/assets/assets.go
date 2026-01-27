package assets

import "embed"

//go:generate npx esbuild *.js *.css --entry-names=[name].min --bundle --outdir=dist/  --minify --sourcemap

//go:embed dist/*
var FS embed.FS
