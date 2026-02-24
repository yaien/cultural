package assets

import "embed"

//go:generate npx esbuild *.js *.css --entry-names=[name].min --bundle --outdir=dist/  --minify --sourcemap
//go:generate sh -c "curl -fsSL -z dist/alpine.min.js https://cdn.jsdelivr.net/npm/alpinejs@3.15.8/dist/cdn.min.js -o dist/alpine.min.js"
//go:generate sh -c "curl -fsSL -z dist/coloris.min.js https://cdn.jsdelivr.net/gh/mdbassit/Coloris@latest/dist/coloris.min.js -o dist/coloris.min.js"
//go:generate sh -c "curl -fsSL -z dist/coloris.min.css https://cdn.jsdelivr.net/gh/mdbassit/Coloris@latest/dist/coloris.min.css -o dist/coloris.min.css"

//go:embed dist/* dist/**/*
var FS embed.FS
