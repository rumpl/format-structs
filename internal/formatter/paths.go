package formatter

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// sourcePaths resolves files, directories and directory/... recursive patterns.
func sourcePaths(root string, patterns []string) ([]string, error) {
	found := map[string]bool{}
	for _, pattern := range patterns {
		recursive := pattern == "..." || strings.HasSuffix(pattern, "/...")
		if recursive {
			pattern = strings.TrimSuffix(pattern, "...")
			if pattern == "" {
				pattern = "."
			}
		}
		name := pattern
		if !filepath.IsAbs(name) {
			name = filepath.Join(root, name)
		}
		name = filepath.Clean(name)
		info, err := os.Stat(name)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			if filepath.Ext(name) != ".go" {
				return nil, fmt.Errorf("not a Go source file: %s", name)
			}
			found[name] = true
			continue
		}
		err = filepath.WalkDir(name, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() && path != name {
				if !recursive || strings.HasPrefix(entry.Name(), ".") || entry.Name() == "vendor" || entry.Name() == "testdata" {
					return filepath.SkipDir
				}
			}
			if entry.Type().IsRegular() && filepath.Ext(path) == ".go" {
				found[path] = true
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	paths := make([]string, 0, len(found))
	for path := range found {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}
