package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDir ensures a directory exists, creating it if necessary
func EnsureDir(path string) error {
	if path == "" {
		return nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("create directory %s: %w", path, err)
	}

	return nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// WriteFile writes content to a file, creating parent directories if needed
func WriteFile(path string, content []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}

	if err := os.WriteFile(path, content, perm); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}

	return nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read source file: %w", err)
	}

	if err := WriteFile(dst, data, 0644); err != nil {
		return fmt.Errorf("write destination file: %w", err)
	}

	return nil
}
