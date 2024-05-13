//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/sh"
)

var Default = Build

func Build() error {
	path := "bin"
	name := "srs-sip"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	name = filepath.Join(path, name)

	if err := sh.Run("go", "build", "-o", name, "main/main.go"); err != nil {
		return err
	}
	fmt.Println("build done")
	return nil
}
