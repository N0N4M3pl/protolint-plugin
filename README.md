# protolint-plugin

Plugin with custom rules for protolint

## With protolint

`protolint lint -plugin plugin.exe filename.proto`

### On Windows system

`protolint lint -plugin "plugin.exe -go_style=false" filename.proto`

## Test

1. `. build.sh`
2. `. test.sh`
3. `protoc --go_out=test/out test/proto/v1/service/snapshot_sevice.proto`

## Usefull commands

`go mod tidy`

`go build -o bin/plugin.exe main.go`

## inspiration

- buf.build
  - [rules](https://buf.build/docs/lint/rules#basic)
- prototool
  - [lint.md](https://github.com/uber/prototool/blob/dev/docs/lint.md)
  - [lint code](https://github.com/uber/prototool/tree/dev/internal/lint)
