//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

var Default = Build

func Build() error {
	fmt.Println("building...")
	if err := os.MkdirAll("bin", 0755); err != nil {
		return err
	}
	if err := sh.Run("go", "build", "-o", "bin/srs-sip.exe", "main/main.go"); err != nil {
		return err
	}
	fmt.Println("build done")
	return nil
}
