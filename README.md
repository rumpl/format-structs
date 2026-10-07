# format-structs

A Go formatting tool that expands nonempty, single-line composite literals
(structs, maps, arrays, and slices).

```go
// Before
config := Config{Port: 8080, Name: "example"}

// After
config := Config{
    Port: 8080,
    Name: "example",
}
```

## Install

Requires Go 1.27 or newer.

```sh
go install github.com/rumpl/format-structs@latest
```

For a local checkout:

```sh
go install .
```

## Usage

Run from your project directory (a Go module is not required):

```sh
format-structs                     # format ./...
format-structs ./cmd/... ./pkg/... # select directory trees
format-structs -check ./...        # report changes without writing
format-structs -root /path/to/project ./...
```
