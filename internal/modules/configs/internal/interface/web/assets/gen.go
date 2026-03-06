package assets

import "embed"

//go:generate npx esbuild dashboard.ts --bundle --outfile=dist/dashboard.min.js --minify --sourcemap

//go:embed dist/*
var FS embed.FS
