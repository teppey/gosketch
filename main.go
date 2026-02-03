package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

type exitCode int

const (
	exitOK    exitCode = 0
	exitError exitCode = 1
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("error: no command")
		os.Exit(int(exitError))
	}

	command := os.Args[1]
	subArgs := os.Args[2:]
	code := exitOK
	switch command {
	case "new":
		code = commandNew(subArgs)
	default:
		fmt.Printf("error: unknown command: %q\n", command)
		code = exitError
	}

	os.Exit(int(code))
}

func commandNew(args []string) exitCode {
	u, err := user.Current()
	if err != nil {
		fmt.Println("error: failed to get current user")
		return exitError
	}

	dir := filepath.Join(os.TempDir(), "goplay_"+u.Username)
	if err := os.Mkdir(dir, 0700); err != nil && !os.IsExist(err) {
		fmt.Printf("error: faild to create directory: %s: %s\n", dir, err)
		return exitError
	}

	curmax := 1
	if paths, err := filepath.Glob(filepath.Join(dir, "[1-9].go")); err == nil {
		for _, path := range paths {
			stem, _ := strings.CutSuffix(filepath.Base(path), ".go")
			n, err := strconv.Atoi(stem)
			if err == nil && n > curmax {
				curmax = n
			}
		}
	}

	fmt.Println(dir)
	fmt.Println(curmax)

	return exitOK
}
