package main

import (
	"os/exec"
	"path/filepath"
	"regexp"
)

var (
	denoVer = regexp.MustCompile(`deno (\d+\.\d+\.\d+)`)
	bunVer  = regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	javaVer = regexp.MustCompile(`(?:openjdk|java) (?:version )?"?(\d+\.\d+\.\d+)`)
	rustVer = regexp.MustCompile(`rustc (\d+\.\d+\.\d+)`)
)

var extraRuntimes = []string{"deno", "bun", "java", "rust"}

func detectExtra(key string) *versionInfo {
	exe := key + exeSuffix()
	switch key {
	case "deno":
		exe = "deno" + exeSuffix()
	case "rust":
		exe = "rustc" + exeSuffix()
	}

	path, err := exec.LookPath(exe)
	if err != nil {
		return nil
	}

	realPath, _ := filepath.EvalSymlinks(path)
	var args []string
	var re *regexp.Regexp

	switch key {
	case "deno":
		args = []string{"--version"}
		re = denoVer
	case "bun":
		args = []string{"--version"}
		re = bunVer
	case "java":
		args = []string{"--version"}
		re = javaVer
	case "rust":
		args = []string{"--version"}
		re = rustVer
	}

	cmd := exec.Command(realPath, args...)
	out, _ := cmd.Output()
	m := re.FindStringSubmatch(string(out))
	if len(m) < 2 {
		return nil
	}

	return &versionInfo{version: m[1], source: "System", path: filepath.Dir(realPath)}
}

func listExtraVersions(key string) []versionInfo {
	if v := detectExtra(key); v != nil {
		return []versionInfo{*v}
	}
	return nil
}

func runtimesWithExtra() []string {
	return append([]string{"python", "php", "node", "go"}, extraRuntimes...)
}
