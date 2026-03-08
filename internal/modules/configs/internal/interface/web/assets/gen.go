package assets

import "embed"

//go:generate npx esbuild dashboard.ts --bundle --outfile=dist/dashboard.min.js --minify --sourcemap  --format=esm

//go:embed dist/*
var FS embed.FS
