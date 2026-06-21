package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
	for _, name := range []string{"python", "php", "node", "go"} {
		if v, ok := versions[name]; ok {
			fmt.Printf("    %s: %s\n", name, v)
		}
	}
}

func fetchLatestVersions() map[string]string {
	result := map[string]string{}

	if v := scrapeVersion("https://www.python.org/downloads/", `Python\s+(\d+\.\d+\.\d+)`); v != "" {
		result["python"] = v
	}
	if v := scrapeVersion("https://windows.php.net/download/", `php-(\d+\.\d+\.\d+)-nts`); v != "" {
		result["php"] = v
	}
	if v := scrapeVersion("https://nodejs.org/dist/latest/", `node-v(\d+\.\d+\.\d+)`); v != "" {
		result["node"] = v
	}
	if v := scrapeVersion("https://go.dev/dl/", `go(\d+\.\d+\.\d+)\.windows`); v != "" {
		result["go"] = v
	}

	return result
}

func scrapeVersion(url, pattern string) string {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body := make([]byte, 512*1024)
	n, _ := resp.Body.Read(body)
	content := string(body[:n])

	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(content)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}
