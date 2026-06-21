package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func cmdDoctor() {
	ok := true

	// Check .pivot/bin in PATH
	home, _ := os.UserHomeDir()
	binDir := filepath.Join(home, ".pivot", "bin")
	curPath := os.Getenv("PATH")
	if strings.Contains(curPath, binDir) {
		fmt.Printf("  \u2713 %s\n", trFmt("doctor_path_ok", binDir))
	} else {
		fmt.Printf("  \u2717 %s\n", trFmt("doctor_path_missing", binDir))
		ok = false
	}

	// Check each runtime
	for _, rt := range []string{"python", "php", "node", "go"} {
		versions := listVersions(rt)

		var active string
		switch rt {
		case "python":
			active = cfg.Python
		case "php":
			active = cfg.PHP
		case "node":
			active = cfg.Node
		case "go":
			active = cfg.Go
		}

		if len(versions) == 0 {
			fmt.Printf("  \u26a0 %s\n", trFmt("doctor_no_versions", runtimeLabel(rt)))
			continue
		}

		// Check for conflicts (multiple system versions)
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

		// Check if active version is still valid
		if active != "" {
			found := false
			for _, v := range versions {
				if v.version == active {
					found = true
					// Verify binary exists
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

	// Check Go GOROOT
	if cfg.Go != "" {
		goRoot := os.Getenv("GOROOT")
		if goRoot == "" {
			fmt.Printf("  \u26a0 %s\n", tr("doctor_no_goroot"))
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

	// Fix PATH
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

	// Fix active versions — re-link binaries
	for _, rt := range runtimesWithExtra() {
		var active string
		switch rt {
		case "python":
			active = cfg.Python
		case "php":
			active = cfg.PHP
		case "node":
			active = cfg.Node
		case "go":
			active = cfg.Go
		}
		if active == "" {
			continue
		}
		bin := filepath.Join(binDir, rt+exeSuffix())
		if _, err := os.Stat(bin); err != nil {
			// Re-activate
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
