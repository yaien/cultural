package assets

import "embed"

//go:generate npx esbuild pages.js --bundle --outfile=dist/pages.min.js --minify --sourcemap

//go:embed *.css *.js dist/*
var FS embed.FS
