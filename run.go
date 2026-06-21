package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func cmdRun(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, tr("usage_run"))
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

	// Build command from remaining args
	var cmdArgs []string
	if len(args) > 2 {
		cmdArgs = args[2:]
	} else {
		cmdArgs = []string{exeName(rt), "--version"}
	}

	// If first arg looks like a flag, prepend the runtime binary
	// e.g. "pivot run python 3 --version" -> python --version
	if len(cmdArgs) > 0 && len(cmdArgs[0]) > 0 && cmdArgs[0][0] == '-' {
		cmdArgs = append([]string{exeName(rt)}, cmdArgs...)
	}

	curPath := os.Getenv("PATH")
	newPath := binDir + string(os.PathListSeparator) + curPath

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = append(os.Environ(), "PATH="+newPath)
	if rt == "go" {
		goRoot := filepath.Join(svDir, "runtimes", "go", v.version)
		if v.source == "System" {
			goRoot = v.path
		}
		cmd.Env = append(cmd.Env, "GOROOT="+goRoot)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "run failed: %v\n", err)
		os.Exit(1)
	}
}
