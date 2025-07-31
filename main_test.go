package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitGitDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gogit-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := initGitDir(tmpDir); err != nil {
		t.Fatalf("initGitDir failed: %v", err)
	}

	dirs := []string{".git", ".git/objects", ".git/refs"}
	for _, dir := range dirs {
		fullPath := filepath.Join(tmpDir, dir)
		info, err := os.Stat(fullPath)
		if err != nil {
			t.Errorf("expected directory %s to exist, got error: %v", fullPath, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", fullPath)
		}
	}

	headPath := filepath.Join(tmpDir, ".git", "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		t.Errorf("expected HEAD file to exist, got error: %v", err)
	}
	expected := "ref: refs/heads/main\n"
	if string(data) != expected {
		t.Errorf("HEAD file contents = %q, want %q", string(data), expected)
	}
}
