// Package formatter expands single-line Go composite literals without type checking.
package formatter

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/scanner"
	"go/token"
	"sort"
)

type sourceFile struct {
	original []byte
	inserts  map[int]string
}

func collect(fset *token.FileSet, file *ast.File, source *sourceFile) {
	ast.Inspect(file, func(node ast.Node) bool {
		literal, ok := node.(*ast.CompositeLit)
		if !ok || len(literal.Elts) == 0 {
			return true
		}
		if fset.Position(literal.Lbrace).Line != fset.Position(literal.Rbrace).Line {
			return true
		}
		offset := func(pos token.Pos) int { return fset.Position(pos).Offset }
		source.inserts[offset(literal.Lbrace)+1] = "\n"
		for i, element := range literal.Elts {
			end := offset(element.End())
			next := offset(literal.Rbrace)
			if i+1 < len(literal.Elts) {
				next = offset(literal.Elts[i+1].Pos())
			}
			// Locate a trailing comma lexically, not inside a string or comment.
			tokens := token.NewFileSet()
			fragment := tokens.AddFile("", -1, next-end)
			var scan scanner.Scanner
			scan.Init(fragment, source.original[end:next], nil, scanner.ScanComments)
			comma := -1
			for {
				pos, tok, _ := scan.Scan()
				if tok == token.EOF {
					break
				}
				if tok == token.COMMA {
					comma = end + fragment.Offset(pos)
					break
				}
			}
			if comma >= 0 {
				source.inserts[comma+1] = "\n"
			} else {
				source.inserts[end] = ",\n"
			}
		}
		return true
	})
}

func (s *sourceFile) format() ([]byte, error) {
	positions := make([]int, 0, len(s.inserts))
	for pos := range s.inserts {
		positions = append(positions, pos)
	}
	sort.Ints(positions)
	var out bytes.Buffer
	previous := 0
	for _, pos := range positions {
		out.Write(s.original[previous:pos])
		out.WriteString(s.inserts[pos])
		previous = pos
	}
	out.Write(s.original[previous:])
	return format.Source(out.Bytes())
}
