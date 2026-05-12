BIN := "bin"
DATA := "./data"

gomod := "izanr.com/chat"
appname := "chat"

debug := "false"

version := `git rev-parse HEAD | head -c8`
goversion := `go version | awk '{print $3}' | sed 's/go//'`
gotags := "generated"

d_ldflags := if debug == "false" { "-s -w" } else { "" }
ldflags := d_ldflags + " -X " + gomod + "/config.Version=" + version + " -X " + gomod + "/config.GoVersion=" + goversion

default: test build

run: build
    {{ BIN / appname }} --config={{ DATA }}/config.json

build arch=`go env GOARCH` os=`go env GOOS`: bin generate
    CGO_ENABLED=0 \
    GOARCH={{ arch }} \
    GOOS={{ os }} \
    go build \
    -tags "{{ gotags }}" \
    -ldflags "{{ ldflags }}" \
    -o {{ BIN / appname }} \
    {{ gomod }}

test force="false" short="false": generate
    flags_1="{{ if force == "true" { "-count=1" } else { "" } }}"; \
    flags_2="{{ if short == "true" { "-short" } else { "" } }}"; \
    CGO_ENABLED=1 \
    go test \
    -tags "{{ gotags }}" \
    -ldflags "{{ ldflags }} -X {{ gomod }}/internal/utils.HashPasswordCost=8" \
    -v -race $flags_1 $flags_2 \
    ./...

generate: deps
    go run ./scripts/sql_generate/main.go
    govalid ./...

deps:
    go install github.com/sivchari/govalid/cmd/govalid@latest

bin:
    mkdir -p {{ BIN }}

update:
    go mod tidy
    go get -u ./...
    go mod tidy

data:
    mkdir -p {{ DATA }}

clear:
    #!/bin/bash
    shopt -s globstar
    rm -f ./**/*_validator.go
    rm -f ./**/*_generate.go
    rm -rf {{ BIN / "*" }}
