package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func cmdShell(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, tr("usage_shell"))
		os.Exit(1)
	}
	rt := args[0]
	ver := args[1]

	versions := listVersions(rt)
	v := findByPrefix(versions, ver)
	if v == nil {
		fmt.Fprintf(os.Stderr, tr("version_not_found"), ver, rt)
		os.Exit(1)
	}

	binDir := v.path
	if rt == "go" {
		alt := filepath.Join(v.path, "bin")
		if _, err := os.Stat(filepath.Join(alt, "go"+exeSuffix())); err == nil {
			binDir = alt
		}
	}

	shell, _ := os.LookupEnv("SHELL")
	if shell == "" {
		shell = "sh"
		if runtime.GOOS == "windows" {
			shell = "cmd"
		}
	}

	curPath := os.Getenv("PATH")
	newPath := binDir + string(os.PathListSeparator) + curPath

	env := os.Environ()
	for i, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			env[i] = "PATH=" + newPath
		}
	}

	if rt == "go" {
		goRoot := filepath.Join(svDir, "runtimes", "go", v.version)
		if v.source == "System" {
			goRoot = v.path
		}
		env = append(env, "GOROOT="+goRoot)
	}

	fmt.Fprintf(os.Stderr, "  pivot: %s %s — type 'exit' to return\n", runtimeLabel(rt), v.version)

	cmd := exec.Command(shell)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
