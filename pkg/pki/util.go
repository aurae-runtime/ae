package pki

import (
	"fmt"
	"os"
	"path/filepath"
)

func createFile(path, filename, content string) error {
	path = filepath.Clean(path)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fp := filepath.Join(path, filename)

	err = writeStringToFile(fp, content)
	if err != nil {
		return err
	}

	return nil
}

func writeStringToFile(p string, s string) error {
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", p, err)
	}
	defer f.Close()

	_, err = f.WriteString(s)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", p, err)
	}

	return nil
}
