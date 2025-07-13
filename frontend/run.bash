cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./static/


(trap 'kill 0' SIGINT; 
air --build.cmd "GOOS=js GOARCH=wasm go build -o static/checkers.wasm go-src/*.go" &
elm-watch hot & 
cd ./static/ && python3 -m http.server)

