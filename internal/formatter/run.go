package formatter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// Run formats selected files or reports required changes when check is true.
// File names are written to output in stable order.
func Run(root string, check bool, patterns []string, output io.Writer) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	if len(patterns) == 0 {
		patterns = []string{
			"./...",
		}
	}
	paths, err := sourcePaths(root, patterns)
	if err != nil {
		return err
	}
	files := map[string]*sourceFile{}
	for _, name := range paths {
		raw, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, name, raw, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		if ast.IsGenerated(file) {
			continue
		}
		source := &sourceFile{
			original: raw,
			inserts:  map[int]string{},
		}
		collect(fset, file, source)
		files[name] = source
	}
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	changed := 0
	for _, name := range names {
		source := files[name]
		result, err := source.format()
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		if bytes.Equal(source.original, result) {
			continue
		}
		changed++
		rel, _ := filepath.Rel(root, name)
		if _, err := fmt.Fprintln(output, rel); err != nil {
			return err
		}
		if !check {
			info, err := os.Stat(name)
			if err != nil {
				return err
			}
			if err := os.WriteFile(name, result, info.Mode().Perm()); err != nil {
				return err
			}
		}
	}
	if check && changed > 0 {
		return fmt.Errorf("%d files need composite literal formatting", changed)
	}
	return nil
}
