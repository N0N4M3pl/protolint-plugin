# protolint-plugin

Plugin with custom rules for protolint

## With protolint

`protolint lint -plugin plugin.exe filename.proto`

### On Windows system

`protolint lint -plugin "plugin.exe -go_style=false" filename.proto`

## Usefull commands

`go mod tidy`

`go build -o bin/plugin.exe main.go`
