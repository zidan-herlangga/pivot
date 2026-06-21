package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 30 * time.Minute}

func downloadRuntime(key, version string) error {
	if version == "" {
		version = defaultVersion(key)
	}
	dir := filepath.Join(svDir, "runtimes", key)
	url, zipName := downloadURL(key, version)
	if url == "" {
		return fmt.Errorf(tr("no_download_for_platform"), runtimeLabel(key), runtime.GOOS)
	}
	dest := filepath.Join(dir, version)
	if _, err := os.Stat(dest); err == nil {
		fmt.Println("  " + trFmt("already_downloaded", runtimeLabel(key), version))
		return nil
	}
	fmt.Println("  " + trFmt("downloading", runtimeLabel(key), version))
	zipPath := filepath.Join(os.TempDir(), zipName)
	if err := downloadFile(url, zipPath); err != nil {
		return fmt.Errorf(tr("download_failed"), err)
	}
	fmt.Println("  " + tr("extracting"))
	os.MkdirAll(dest, 0755)
	if err := extractArchive(zipPath, dest); err != nil {
		return fmt.Errorf(tr("extract_failed"), err)
	}
	// Java & Deno extract with top-level dir — flatten
	fixExtractedLayout(key, dest)
	os.Remove(zipPath)
	return nil
}

func fixExtractedLayout(key, dest string) {
	switch key {
	case "java", "deno", "go":
		entries, _ := os.ReadDir(dest)
		if len(entries) == 1 && entries[0].IsDir() {
			sub := filepath.Join(dest, entries[0].Name())
			copyDirContents(sub, dest)
			os.RemoveAll(sub)
		}
	}
}

func copyDirContents(src, dest string) {
	entries, _ := os.ReadDir(src)
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dest, e.Name())
		if e.IsDir() {
			os.MkdirAll(d, 0755)
			copyDirContents(s, d)
		} else {
			data, _ := os.ReadFile(s)
			os.WriteFile(d, data, 0755)
		}
	}
}

func defaultVersion(key string) string {
	switch key {
	case "python":
		return "3.12.4"
	case "php":
		return "8.3.8"
	case "node":
		return "20.15.0"
	case "go":
		return "1.22.5"
	case "java":
		return "21"
	case "deno":
		return "1.44.4"
	case "bun":
		return "1.1.17"
	}
	return ""
}

func downloadURL(key, version string) (string, string) {
	goos := runtime.GOOS
	switch key {
	case "python":
		if goos == "windows" {
			return fmt.Sprintf("https://www.python.org/ftp/python/%s/python-%s-embed-amd64.zip", version, version), fmt.Sprintf("python-%s.zip", version)
		}
		return fmt.Sprintf("https://www.python.org/ftp/python/%s/Python-%s.tgz", version, version), fmt.Sprintf("Python-%s.tgz", version)
	case "php":
		if goos != "windows" {
			return "", ""
		}
		vs := "vs16"
		if compareVersions(version, "8.5.0") >= 0 {
			vs = "vs17"
		}
		return fmt.Sprintf("https://windows.php.net/downloads/releases/php-%s-nts-Win32-%s-x64.zip", version, vs), fmt.Sprintf("php-%s.zip", version)
	case "node":
		if goos == "windows" {
			return fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%s-win-x64.zip", version, version), fmt.Sprintf("node-v%s.zip", version)
		}
		return fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%s-%s-x64.tar.gz", version, version, goos), fmt.Sprintf("node-v%s.tar.gz", version)
	case "go":
		if goos == "windows" {
			return fmt.Sprintf("https://go.dev/dl/go%s.windows-amd64.zip", version), fmt.Sprintf("go%s.zip", version)
		}
		return fmt.Sprintf("https://go.dev/dl/go%s.%s-amd64.tar.gz", version, goos), fmt.Sprintf("go%s.tar.gz", version)
	case "java":
		osMap := map[string]string{"windows": "windows", "linux": "linux", "darwin": "mac"}
		osStr := osMap[goos]
		if osStr == "" {
			return "", ""
		}
		ext := "tar.gz"
		if goos == "windows" {
			ext = "zip"
		}
		return fmt.Sprintf("https://api.adoptium.net/v3/binary/latest/%s/ga/%s/x64/jdk/hotspot/normal/eclipse?project=jdk", version, osStr), fmt.Sprintf("java-%s.%s", version, ext)
	case "deno":
		archMap := map[string]string{"windows": "x86_64-pc-windows-msvc", "linux": "x86_64-unknown-linux-gnu", "darwin": "x86_64-apple-darwin"}
		arch := archMap[goos]
		if arch == "" {
			return "", ""
		}
		return fmt.Sprintf("https://github.com/denoland/deno/releases/download/v%s/deno-%s.zip", version, arch), fmt.Sprintf("deno-%s.zip", version)
	case "bun":
		osMap := map[string]string{"windows": "windows", "linux": "linux", "darwin": "darwin"}
		osStr := osMap[goos]
		if osStr == "" {
			return "", ""
		}
		return fmt.Sprintf("https://github.com/oven-sh/bun/releases/download/bun-v%s/bun-%s-x64.zip", version, osStr), fmt.Sprintf("bun-%s.zip", version)
	}
	return "", ""
}

func downloadFile(url, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "pivot/2.0")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	_, err = io.Copy(out, resp.Body)
	return err
}

func extractArchive(src, dest string) error {
	if strings.HasSuffix(src, ".zip") {
		return unzip(src, dest)
	}
	if strings.HasSuffix(src, ".tar.gz") || strings.HasSuffix(src, ".tgz") {
		return untargz(src, dest)
	}
	return fmt.Errorf("unknown archive format: %s", src)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), 0755)
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.Create(fpath)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func untargz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fpath := filepath.Join(dest, hdr.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}
		if hdr.Typeflag == tar.TypeDir {
			os.MkdirAll(fpath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), 0755)
		out, err := os.Create(fpath)
		if err != nil {
			return err
		}
		_, err = io.Copy(out, tr)
		out.Close()
		if err != nil {
			return err
		}
	}
	return nil
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
