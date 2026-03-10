## WebAssembly TODO demo

This folder contains a tiny TODO app written in Go and compiled to WebAssembly.

### Run it (the simplest way)

1. Install Go 1.22+ from https://go.dev/dl/.
2. Right-click this `misc/wasm` folder and choose **Open in Terminal/PowerShell**.
3. Run this single command:

   ```sh
   GOOS=js GOARCH=wasm go build -o main.wasm main.go && go run web/web.go
   ```

   On Windows PowerShell:

   ```pwsh
   $env:GOOS="js"; $env:GOARCH="wasm"; go build -o main.wasm main.go; go run web/web.go
   ```

4. Your browser opens to http://localhost:8080 where you can add, check off, and delete TODOs. Tasks are saved in `localStorage` so refreshes keep your list.

The `web/web.go` helper serves the current folder and opens the browser for you. If you prefer to serve the files yourself, any static file server works as long as `main.wasm` and `wasm_exec.js` are reachable from `index.html`.
