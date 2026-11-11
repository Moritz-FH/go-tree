package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateTreeJSON(t *testing.T) {
	tempDir := t.TempDir()

	subDir1 := filepath.Join(tempDir, "subdir1")
	subDir2 := filepath.Join(subDir1, "subdir2")
	emptyDir := filepath.Join(tempDir, "emptydir")

	if err := os.Mkdir(subDir1, 0o755); err != nil {
		t.Fatalf("failed to create subdir1: %v", err)
	}
	if err := os.Mkdir(subDir2, 0o755); err != nil {
		t.Fatalf("failed to create subdir2: %v", err)
	}
	if err := os.Mkdir(emptyDir, 0o755); err != nil {
		t.Fatalf("failed to create emptydir: %v", err)
	}

	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(subDir1, "file2.txt")
	if err := os.WriteFile(file1, []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("world"), 0o644); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	symlinkPath := filepath.Join(tempDir, "symlink_file")
	if err := os.Symlink(file1, symlinkPath); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	t.Run("full tree output", func(t *testing.T) {
		data, err := generateTreeJSON(tempDir, false, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var result []map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if len(result) != 2 {
			t.Fatalf("expected 2 elements in result (tree and report), got %d", len(result))
		}

		treeNode := result[0]
		if treeNode["type"] != "directory" {
			t.Errorf("expected tree type 'directory', got %v", treeNode["type"])
		}
		if treeNode["name"] != tempDir {
			t.Errorf("expected root name %q, got %v", tempDir, treeNode["name"])
		}

		report := result[1]
		if report["type"] != "report" {
			t.Errorf("expected report type 'report', got %v", report["type"])
		}
		if report["directories"] != float64(4) {
			t.Errorf("expected 4 directories, got %v", report["directories"])
		}
		if report["files"] != float64(3) {
			t.Errorf("expected 3 files, got %v", report["files"])
		}
	})

	t.Run("directory only flag", func(t *testing.T) {
		data, err := generateTreeJSON(tempDir, true, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var result []map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		report := result[1]
		if report["directories"] != float64(4) {
			t.Errorf("expected 4 directories, got %v", report["directories"])
		}
		if _, exists := report["files"]; exists {
			t.Errorf("files should not be present in report when dirOnly=true, got %v", report["files"])
		}
	})

	t.Run("max depth limit", func(t *testing.T) {
		data, err := generateTreeJSON(tempDir, false, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var result []map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		report := result[1]
		if report["directories"] != float64(3) {
			t.Errorf("expected 3 directories at depth 1, got %v", report["directories"])
		}
		if report["files"] != float64(2) {
			t.Errorf("expected 2 files at depth 1, got %v", report["files"])
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		emptyTestDir := t.TempDir()
		data, err := generateTreeJSON(emptyTestDir, false, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var result []map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		treeNode := result[0]
		if _, exists := treeNode["contents"]; exists {
			t.Errorf("empty dir should not have contents key, got %v", treeNode["contents"])
		}

		report := result[1]
		if report["directories"] != float64(0) {
			t.Errorf("expected 0 directories, got %v", report["directories"])
		}
		if report["files"] != float64(0) {
			t.Errorf("expected 0 files, got %v", report["files"])
		}
	})

	t.Run("non-existent directory error", func(t *testing.T) {
		_, err := generateTreeJSON(filepath.Join(tempDir, "does-not-exist"), false, -1)
		if err == nil {
			t.Fatal("expected error for non-existent directory, got nil")
		}
	})
}

func TestPrintTree(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "sub")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatalf("failed to create sub dir: %v", err)
	}
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("data"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err := printTree(tempDir, "", 0, false, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
