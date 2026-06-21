package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func cmdClean() {
	dir := filepath.Join(svDir, "runtimes")
	total := 0

	for _, rt := range []string{"python", "php", "node", "go"} {
		rtDir := filepath.Join(dir, rt)
		entries, err := os.ReadDir(rtDir)
		if err != nil {
			continue
		}

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

		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			if e.Name() == active {
				continue
			}
			verDir := filepath.Join(rtDir, e.Name())
			if err := os.RemoveAll(verDir); err == nil {
				fmt.Printf("  %s %s\n", tr("removed"), runtimeLabel(rt)+" "+e.Name())
				total++
			}
		}
	}

	if total == 0 {
		fmt.Println("  " + tr("nothing_to_clean"))
	} else {
		fmt.Println("  " + trFmt("cleaned_versions", total))
	}
}
