package formatter

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func writeFixture(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0640); err != nil {
		t.Fatal(err)
	}
}

type failingOutput struct{}

func (failingOutput) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestOutputFailureStopsBeforeWriting(t *testing.T) {
	root := t.TempDir()
	original := "package sample\nvar a=[]int{1,2}\n"
	writeFixture(t, root, "sample.go", original)
	if err := Run(root, false, nil, failingOutput{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("output error: %v", err)
	}
	raw, _ := os.ReadFile(filepath.Join(root, "sample.go"))
	if string(raw) != original {
		t.Fatal("source changed after reporting failed")
	}
}

func TestRunPreservesPermissions(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "sample.go", "package sample\nvar x=[]int{1,2}\n")
	if err := Run(root, false, nil, io.Discard); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(root, "sample.go"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0640 {
		t.Fatalf("permissions changed: %v", info.Mode())
	}
}
