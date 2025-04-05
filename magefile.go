//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

var (
	//targets = []string{"app1"}
	targets = map[string]string{
		"app1": "0.0.1",
	}
)

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	for k, v := range targets {
		fmt.Printf("ビルド開始。ターゲット名: %s\n", k)
		cmd1 := exec.Command("go", "build", "-o",
			fmt.Sprintf("bin/%s", k), fmt.Sprintf("cmd/%s/main.go", k))
		cmd1.Stdout = os.Stdout
		cmd1.Stderr = os.Stderr

		if err := cmd1.Run(); err != nil {
			return err
		}

		fmt.Printf("イメージ作成開始。ターゲット名: %s\n", k)
		cmd2 := exec.Command("podman", "build", "-t",
			fmt.Sprintf("localhost/%s:%s", k, v), "-f", fmt.Sprintf("cmd/%s/Dockerfile", k), ".")
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr

		if err := cmd2.Run(); err != nil {
			return err
		}
	}

	return nil
}

// Manage your deps, or running package managers.
// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("bin")
}
