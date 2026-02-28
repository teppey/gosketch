package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

const (
	scribbleFile = "scribble.go"
	scribbleOldFile = "scribble.go.old"
	defaultTemplate = `package main

import "fmt"

func main() {
	fmt.Println("Hello, world")
}
`
)

func baseDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", errors.New("failed to get current user")
	}

	return filepath.Join(os.TempDir(), "goscribble_"+u.Username), nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
		return false
	}

	return true
}

func edit(path string) error {
	stdin, err := os.Open("/dev/tty")
	if err != nil {
		return err
	}
	defer stdin.Close()

	cmd := exec.Command("vim", path)
	cmd.Stdin = stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to edit file: %s: %w", cmd, err)
	}

	return nil
}

func goimports(path string) error {
	cmd := exec.Command("goimports", "-w", path)
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Print(string(out))
	}

	return err
}

func goModInit(modDir string) error {
	if exists(filepath.Join(modDir, "go.mod")) {
		return nil
	}

	cmd := exec.Command("go", "mod", "init", filepath.Base(modDir))
	cmd.Dir = modDir
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Print(string(out))
	}

	return err
}

func goModTidy(modDir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = modDir
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Printf("go mod tidy: %s", string(out))
	}

	return err
}

func run(path string) error {
	cmd := exec.Command("go", "run", path)
	cmd.Dir = filepath.Dir(path)
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Print(string(out))
	}

	return err
}
