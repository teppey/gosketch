package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func newCommand(args []string) error {
	template := defaultTemplate
	var data []byte
	var err error
	if len(args) > 0 {
		if args[0] == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(args[0])
		}
		if err != nil {
			return err
		}
		template = string(data)
	}

	dir, err := baseDir()
	if err != nil {
		return err
	}

	if err := os.Mkdir(dir, 0700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create directory: %s: %w", dir, err)
	}

	if err := goModInit(dir); err != nil {
		return err
	}

	path := filepath.Join(dir, scribbleFile)
	if exists(path) {
		oldPath := filepath.Join(dir, scribbleOldFile)
		if err := os.Rename(path, oldPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	if err := os.WriteFile(path, []byte(template), 0600); err != nil {
		return fmt.Errorf("failed to write file: %s: %w", path, err)
	}

	if err := edit(path); err != nil {
		return err
	}

	if err := goModTidy(dir); err != nil {
		return err
	}

	if err := goimports(path); err != nil {
		return err
	}

	if err := run(path); err != nil {
		return err
	}

	return nil
}
