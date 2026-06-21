package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func addToPath(dir string, key string) {
	home, _ := os.UserHomeDir()
	binDir := filepath.Join(home, ".pivot", "bin")
	os.MkdirAll(binDir, 0755)

	var targets []string
	switch key {
	case "python":
		targets = []string{"python", "python3"}
	case "php":
		targets = []string{"php"}
	case "node":
		targets = []string{"node", "npm", "npx"}
	case "go":
		targets = []string{"go"}
		if cfg.Go != "" {
			goRoot := filepath.Join(svDir, "runtimes", "go", cfg.Go)
			os.Setenv("GOROOT", goRoot)
		}
	case "deno":
		targets = []string{"deno"}
	case "bun":
		targets = []string{"bun"}
	case "java":
		targets = []string{"java"}
		// Set JAVA_HOME
		if cfg.Java != "" {
			javaHome := filepath.Join(svDir, "runtimes", "java", cfg.Java)
			os.Setenv("JAVA_HOME", javaHome)
		}
	case "rust":
		targets = []string{"rustc", "cargo"}
	}

	for _, name := range targets {
		src := filepath.Join(dir, name+exeSuffix())
		if _, err := os.Stat(src); err != nil {
			alt := filepath.Join(dir, "bin", name+exeSuffix())
			if _, err2 := os.Stat(alt); err2 != nil {
				continue
			}
			src = alt
		}
		dst := filepath.Join(binDir, name+exeSuffix())
		os.Remove(dst)
		if runtime.GOOS == "windows" {
			copyFile(src, dst)
		} else {
			os.Symlink(src, dst)
		}
	}
}

func copyFile(src, dest string) {
	data, err := os.ReadFile(src)
	if err != nil {
		return
	}
	os.WriteFile(dest, data, 0755)
}

func printEnv() {
	home, _ := os.UserHomeDir()
	binDir := filepath.Join(home, ".pivot", "bin")

	if runtime.GOOS == "windows" {
		fmt.Printf("SET PATH=%s;%%PATH%%\n", binDir)
	} else {
		fmt.Printf("export PATH=\"%s:$PATH\"\n", binDir)
	}
}
