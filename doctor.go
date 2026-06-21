package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var doctorRuntimes = []string{"python", "php", "node", "go", "deno", "bun", "java", "rust"}

func activeVersion(rt string) string {
	switch rt {
	case "python":
		return cfg.Python
	case "php":
		return cfg.PHP
	case "node":
		return cfg.Node
	case "go":
		return cfg.Go
	case "deno":
		return cfg.Deno
	case "bun":
		return cfg.Bun
	case "java":
		return cfg.Java
	case "rust":
		return cfg.Rust
	}
	return ""
}

func cmdDoctor() {
	ok := true

	home, _ := os.UserHomeDir()
	binDir := filepath.Join(home, ".pivot", "bin")
	curPath := os.Getenv("PATH")
	if strings.Contains(curPath, binDir) {
		fmt.Printf("  \u2713 %s\n", trFmt("doctor_path_ok", binDir))
	} else {
		fmt.Printf("  \u2717 %s\n", trFmt("doctor_path_missing", binDir))
		ok = false
	}

	for _, rt := range doctorRuntimes {
		versions := listVersions(rt)
		active := activeVersion(rt)

		if len(versions) == 0 {
			fmt.Printf("  \u26a0 %s\n", trFmt("doctor_no_versions", runtimeLabel(rt)))
			continue
		}

		sysCount := 0
		for _, v := range versions {
			if v.source == "System" {
				sysCount++
			}
		}
		if sysCount > 1 {
			fmt.Printf("  \u26a0 %s\n", trFmt("doctor_conflict", runtimeLabel(rt)))
			for _, v := range versions {
				if v.source == "System" {
					fmt.Printf("         %s (%s)\n", v.version, v.path)
				}
			}
		}

		if active != "" {
			found := false
			for _, v := range versions {
				if v.version == active {
					found = true
					bin := filepath.Join(binDir, rt+exeSuffix())
					if _, err := os.Stat(bin); err != nil {
						fmt.Printf("  \u2717 %s\n", trFmt("doctor_bin_missing", runtimeLabel(rt), bin))
						ok = false
					}
					break
				}
			}
			if !found {
				fmt.Printf("  \u2717 %s\n", trFmt("doctor_active_missing", runtimeLabel(rt), active))
				ok = false
			} else {
				fmt.Printf("  \u2713 %s %s\n", runtimeLabel(rt), active)
			}
		} else {
			fmt.Printf("  \u26a0 %s\n", trFmt("doctor_not_active", runtimeLabel(rt)))
		}
	}

	if cfg.Go != "" {
		goRoot := os.Getenv("GOROOT")
		if goRoot == "" {
			fmt.Printf("  \u26a0 %s\n", tr("doctor_no_goroot"))
		}
	}
	if cfg.Java != "" {
		javaHome := os.Getenv("JAVA_HOME")
		if javaHome == "" {
			fmt.Printf("  \u26a0 %s\n", tr("doctor_no_javahome"))
		}
	}

	if ok {
		fmt.Println("\n  " + tr("doctor_all_good"))
	} else {
		fmt.Println("\n  " + tr("doctor_issues_found"))
	}
}

func cmdDoctorFix() {
	home, _ := os.UserHomeDir()
	binDir := filepath.Join(home, ".pivot", "bin")
	curPath := os.Getenv("PATH")

	fixed := false

	if !strings.Contains(curPath, binDir) {
		if runtime.GOOS == "windows" {
			curPath = binDir + ";" + curPath
		} else {
			curPath = binDir + ":" + curPath
		}
		os.Setenv("PATH", curPath)
		fmt.Println("  " + trFmt("doctor_fixed_path", binDir))
		fixed = true
	}

	for _, rt := range runtimesWithExtra() {
		active := activeVersion(rt)
		if active == "" {
			continue
		}
		bin := filepath.Join(binDir, rt+exeSuffix())
		if _, err := os.Stat(bin); err != nil {
			versions := listVersions(rt)
			v := findByPrefix(versions, active)
			if v != nil {
				activateVersion(rt, *v)
				fmt.Println("  " + trFmt("doctor_fixed_bin", runtimeLabel(rt)))
				fixed = true
			}
		}
	}

	if !fixed {
		fmt.Println("\n  " + tr("doctor_all_good"))
	}
}
