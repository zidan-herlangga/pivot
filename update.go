package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type updateCache struct {
	Timestamp time.Time         `json:"timestamp"`
	Versions  map[string]string `json:"versions"`
}

func checkUpdates(dir string) {
	cachePath := filepath.Join(dir, "update-cache.json")
	cache := loadUpdateCache(cachePath)
	if cache != nil && time.Since(cache.Timestamp).Hours() < 24 {
		printUpdateVersions(cache.Versions)
		return
	}

	fmt.Println("  " + tr("checking_updates"))
	versions := fetchLatestVersions()
	if versions != nil {
		saveUpdateCache(cachePath, versions)
		printUpdateVersions(versions)
	} else {
		fmt.Println("  " + tr("update_failed"))
	}
}

func loadUpdateCache(path string) *updateCache {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var c updateCache
	json.Unmarshal(data, &c)
	return &c
}

func saveUpdateCache(path string, versions map[string]string) {
	c := updateCache{Timestamp: time.Now(), Versions: versions}
	data, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(path, data, 0644)
}

func printUpdateVersions(versions map[string]string) {
	fmt.Println("\n  " + tr("latest_versions"))
	for _, name := range []string{"python", "php", "node", "go", "java", "deno", "bun"} {
		if v, ok := versions[name]; ok {
			fmt.Printf("    %s: %s\n", name, v)
		}
	}
}

func fetchLatestVersions() map[string]string {
	type result struct {
		key string
		val string
	}
	ch := make(chan result, 7)

	fetch := func(url, pattern, key string) {
		v := scrapeVersion(url, pattern)
		ch <- result{key, v}
	}

	go fetch("https://www.python.org/downloads/", `Python\s+(\d+\.\d+\.\d+)`, "python")
	go fetch("https://windows.php.net/download/", `php-(\d+\.\d+\.\d+)-nts`, "php")
	go fetch("https://nodejs.org/dist/latest/", `node-v(\d+\.\d+\.\d+)`, "node")
	go fetch("https://go.dev/dl/", `go(\d+\.\d+\.\d+)\.windows`, "go")
	go fetch("https://github.com/adoptium/temurin21-binaries/releases", `jdk-(\d+\.\d+\.\d+\+\d+)`, "java")
	go fetch("https://github.com/denoland/deno/releases/latest", `deno\s+v?(\d+\.\d+\.\d+)`, "deno")
	go fetch("https://github.com/oven-sh/bun/releases/latest", `bun-v?(\d+\.\d+\.\d+)`, "bun")

	m := make(map[string]string, 7)
	for i := 0; i < 7; i++ {
		r := <-ch
		if r.val != "" {
			m[r.key] = r.val
		}
	}
	return m
}
