# go-dlv-manager

go-delve DAP remote debugger process manager

Make sure that you have same file path or it wont debug/compile.

You can use rsync, or similar software.

## Install

```sh
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/ii64/go-dlv-manager@latest
```

## Example

File `launch.json`

```json
{
    "configurations": [
    {
        "name": "Launch file",
        "type": "go",
        "request": "launch",
        "mode": "exec",        // exec, debug, see more at `dlv dap --help`
        "program": "${file}",
        "showRegisters": true,  // unecessary
        "logOutput": "rpc",     // unecessary
        "port": 7456,
        "host": "127.0.0.1"
    }
    ]
}
```

File `.dlvman.toml`

```toml
[conf]
  Addr = "0.0.0.0:7456"
  Program = "dlv"
  ProcessArgs = ["dap", "--check-go-version=false", "--log-dest=3"]
  Debug = false

```

## Figure

Both running on the same computer, arm64 running on top of QEMU emulation.
It took some time to build the program inside the emulator, but this can be avoided by compiling it on the host machine. 

`GOARCH=arm64 go build -gcflags "all=-N -l" main.go`

On `launch.json`, you just need to change these line:
```json
    "mode": "exec",
    "program": "${workspaceFolder}/main",
```

| amd64 | arm64 (Cortex-A53) |
| ----- | ----- |
| ![amd64][fig0] | ![arm64][fig1] |

## License

The `go-dlv-manager` is under the [MIT license](LICENSE).


[fig0]: ./assets//x86_64.png
[fig1]: ./assets/arm64.png