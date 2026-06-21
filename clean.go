package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func cmdClean() {
	dir := filepath.Join(svDir, "runtimes")
	total := 0

	for _, rt := range allRuntimes() {
		rtDir := filepath.Join(dir, rt)
		entries, err := os.ReadDir(rtDir)
		if err != nil {
			continue
		}

		active := activeVersion(rt)

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
