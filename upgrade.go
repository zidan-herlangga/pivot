package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type ghRelease struct {
	TagName string      `json:"tag_name"`
	Assets  []ghAsset   `json:"assets"`
}

type ghAsset struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func cmdUpgrade() {
	fmt.Println("  " + tr("checking_upgrades"))

	// Fetch latest release
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/zidan-herlangga/pivot/releases/latest")
	if err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_fetch_failed"))
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var release ghRelease
	if err := json.Unmarshal(body, &release); err != nil || release.TagName == "" {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_fetch_failed"))
		return
	}

	// Compare versions
	latestVer := release.TagName
	if compareVersions(version, latestVer) >= 0 {
		fmt.Println("  " + trFmt("upgrade_latest", version))
		return
	}

	// Find matching asset
	arch := "amd64"
	if runtime.GOARCH == "386" {
		arch = "386"
	}
	var assetName string
	switch runtime.GOOS {
	case "windows":
		assetName = fmt.Sprintf("pivot-windows-%s.zip", arch)
	case "linux":
		assetName = fmt.Sprintf("pivot-linux-%s.tar.gz", arch)
	case "darwin":
		assetName = fmt.Sprintf("pivot-darwin-%s.tar.gz", arch)
	default:
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_unsupported"))
		return
	}

	var downloadURL string
	for _, a := range release.Assets {
		if a.Name == assetName {
			downloadURL = a.URL
			break
		}
	}
	if downloadURL == "" {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_no_asset"))
		return
	}

	fmt.Println("  " + trFmt("upgrade_downloading", latestVer))

	// Download to temp
	tmpDir, _ := os.MkdirTemp("", "pivot-upgrade")
	defer os.RemoveAll(tmpDir)
	archivePath := filepath.Join(tmpDir, assetName)

	req, _ := http.NewRequest("GET", downloadURL, nil)
	req.Header.Set("Accept", "application/octet-stream")
	resp2, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_download_failed"))
		return
	}
	defer resp2.Body.Close()

	out, _ := os.Create(archivePath)
	io.Copy(out, resp2.Body)
	out.Close()

	// Extract
	if err := extractArchive(archivePath, tmpDir); err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_extract_failed"))
		return
	}

	// Find the binary
	exeName := "pivot" + exeSuffix()
	newBinary := filepath.Join(tmpDir, exeName)
	if _, err := os.Stat(newBinary); err != nil {
		// Try subdirectory
		entries, _ := os.ReadDir(tmpDir)
		for _, e := range entries {
			if e.IsDir() {
				candidate := filepath.Join(tmpDir, e.Name(), exeName)
				if _, err2 := os.Stat(candidate); err2 == nil {
					newBinary = candidate
					break
				}
			}
		}
	}

	// Replace current binary
	self, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_failed"))
		return
	}

	data, err := os.ReadFile(newBinary)
	if err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_failed"))
		return
	}

	if err := os.WriteFile(self, data, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_failed"))
		return
	}

	fmt.Println("  " + trFmt("upgrade_done", latestVer))
}
