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
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name               string `json:"name"`
	URL                string `json:"url"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func cmdUpgrade() {
	fmt.Println("  " + tr("checking_upgrades"))

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/zidan-herlangga/pivot/releases/latest")
	if err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_fetch_failed")+": "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var release ghRelease
	if err := json.Unmarshal(body, &release); err != nil || release.TagName == "" {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_fetch_failed"))
		return
	}

	latestVer := release.TagName
	if compareVersions(version, latestVer) >= 0 {
		fmt.Println("  " + trFmt("upgrade_latest", version))
		return
	}

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
			downloadURL = a.BrowserDownloadURL
			if downloadURL == "" {
				downloadURL = a.URL
			}
			break
		}
	}
	if downloadURL == "" {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_no_asset"))
		return
	}

	fmt.Println("  " + trFmt("upgrade_downloading", latestVer))

	tmpDir, _ := os.MkdirTemp("", "pivot-upgrade")
	defer os.RemoveAll(tmpDir)
	archivePath := filepath.Join(tmpDir, assetName)

	req, _ := http.NewRequest("GET", downloadURL, nil)
	req.Header.Set("Accept", "application/octet-stream")
	resp2, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_download_failed")+": "+err.Error())
		return
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		msg, _ := io.ReadAll(resp2.Body)
		fmt.Fprintf(os.Stderr, "  "+tr("upgrade_download_failed")+": HTTP %d\n  %s\n", resp2.StatusCode, string(msg))
		return
	}

	out, err := os.Create(archivePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_failed")+": create temp: "+err.Error())
		return
	}
	_, err = io.Copy(out, resp2.Body)
	out.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_failed")+": download: "+err.Error())
		return
	}

	if err := extractArchive(archivePath, tmpDir); err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_extract_failed")+": "+err.Error())
		return
	}

	exeName := "pivot" + exeSuffix()
	newBinary := filepath.Join(tmpDir, exeName)
	if _, err := os.Stat(newBinary); err != nil {
		entries, _ := os.ReadDir(tmpDir)
		found := false
		for _, e := range entries {
			if e.IsDir() {
				candidate := filepath.Join(tmpDir, e.Name(), exeName)
				if _, err2 := os.Stat(candidate); err2 == nil {
					newBinary = candidate
					found = true
					break
				}
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "  "+tr("upgrade_failed")+": binary not found in archive\n")
			return
		}
	}

	self, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_failed")+": self: "+err.Error())
		return
	}

	data, err := os.ReadFile(newBinary)
	if err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_failed")+": read: "+err.Error())
		return
	}

	if runtime.GOOS == "windows" {
		old := self + ".old"
		os.Remove(old)
		if err := os.Rename(self, old); err != nil {
			fmt.Fprintln(os.Stderr, "  "+tr("upgrade_failed")+": rename: "+err.Error())
			return
		}
		defer os.Remove(old)
	}

	if err := os.WriteFile(self, data, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "  "+tr("upgrade_failed")+": write: "+err.Error())
		return
	}

	fmt.Println("  " + trFmt("upgrade_done", latestVer))
}
