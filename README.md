# pengu-lang

A small Domain Specific Language (DSL) for defining microservices with a very simple, readable syntax. The DSL compiles into Go code that runs a production-ready HTTP server using Go's standard `net/http` package.

## Features
- Extremely minimal, clean syntax for defining HTTP services
- Compiles down to standard Go Code
- No complex setup required
- Built-in `log` and `respond` primitives

## Installation
Clone the repository and build the CLI tool:

```bash
git clone https://github.com/yourusername/pengu-lang
cd pengu-lang
go build -o pengu ./cli/main.go
```

## Usage commands

* `pengu init` : Creates starter `.ms` files in `./examples`
* `pengu generate <file.ms>` : Only generates the Go code into `./generated`
* `pengu build <file.ms>` : Generates the Go code and builds it into a binary executable
* `pengu run <file.ms>` : Generates the Go code and runs it immediately using `go run`

## Syntax Example
Save this as `service.ms`:

```
version 1

service payment

route POST "/pay"
    log "payment request"
    respond 200 "ok"

route GET "/health"
    respond 200 "healthy"
```

## DSL Rules
1. Every file must begin with `version 1`
2. At most one `service` per file. Example: `service user`
3. A service can contain multiple `route` definitions. Example: `route GET "/profile"`
4. Actions inside routes (`log`, `respond`) MUST be indented. They belong to the preceding route block.
5. Currently supported actions:
    - `log "message"`
    - `respond status_code "message"`
