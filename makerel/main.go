package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	targets := []struct {
		os   string
		arch []string
	}{
		{"linux",
			[]string{
				"386",
				"amd64",
				"arm",
				"arm64",
				"loong64",
				"mips",
				"mips64",
				"mips64le",
				"mipsle",
				"ppc64",
				"ppc64le",
				"riscv64",
				"s390x",
			},
		},
		{
			"android",
			[]string{"arm64"},
		},
	}

	for _, t := range targets {
		for _, arch := range t.arch {
			build := exec.Command("go", "build")
			build.Stderr = os.Stderr
			build.Stdout = os.Stdout
			build.Env = append(os.Environ(), "GOOS="+t.os, "GOARCH="+arch)
			if err := build.Run(); err != nil {
				panic(err)
			}
			zip := exec.Command("zip", "")
			zip.Stderr = os.Stderr
			zip.Stdout = os.Stdout
			zip.Args = []string{
				"-1",
				fmt.Sprintf("sblast_%s_%s.zip", t.os, arch),
				"sblast",
				"LICENSE",
				"README.md",
			}
			if err := zip.Run(); err != nil {
				panic(err)
			}
			if err := os.Remove("sblast"); err != nil {
				panic(err)
			}
		}
	}
}
