package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type versionInfo struct {
	version string
	source  string
	path    string
}

var (
	pyVer   = regexp.MustCompile(`Python (\d+\.\d+\.\d+)`)
	phpVer  = regexp.MustCompile(`PHP (\d+\.\d+\.\d+)`)
	nodeVer = regexp.MustCompile(`v(\d+\.\d+\.\d+)`)
	goVer   = regexp.MustCompile(`go(\d+\.\d+\.\d+)`)
)

type configData struct {
	Python string `json:"python"`
	PHP    string `json:"php"`
	Node   string `json:"node"`
	Go     string `json:"go"`
}

var svDir string
var cfg configData

func loadConfig(dir string) {
	svDir = dir
	data, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		return
	}
	json.Unmarshal(data, &cfg)
}

func saveConfig() {
	os.MkdirAll(svDir, 0755)
	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(filepath.Join(svDir, "config.json"), data, 0644)
}

func activePython() string { return cfg.Python }
func activePHP() string    { return cfg.PHP }
func activeNode() string   { return cfg.Node }
func activeGo() string     { return cfg.Go }

func runtimeLabel(key string) string {
	switch key {
	case "python":
		return "Python"
	case "php":
		return "PHP"
	case "node":
		return "Node.js"
	case "go":
		return "Go"
	}
	return key
}

// ---- Detection ----
func listVersions(key string) []versionInfo {
	dir := filepath.Join(svDir, "runtimes", key)
	versions := detectPortable(dir, key, exeName(key))
	if sys := detectSystem(key); sys != nil {
		versions = append(versions, *sys)
	}
	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i].version, versions[j].version) > 0
	})
	return versions
}

func exeName(key string) string {
	name := key
	if key == "node" {
		return "node" + exeSuffix()
	}
	return name + exeSuffix()
}

func detectPortable(dir, key, exe string) []versionInfo {
	var versions []versionInfo
	entries, err := os.ReadDir(dir)
	if err != nil {
		return versions
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		exePath := filepath.Join(dir, e.Name(), exe)
		if _, err := os.Stat(exePath); err == nil {
			versions = append(versions, versionInfo{version: e.Name(), source: "Portable", path: filepath.Dir(exePath)})
		}
	}
	// Also check nested bin/ for Go (go installs to bin/go.exe)
	if key == "go" {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			binPath := filepath.Join(dir, e.Name(), "bin", "go"+exeSuffix())
			if _, err := os.Stat(binPath); err == nil {
				versions = append(versions, versionInfo{version: e.Name(), source: "Portable", path: filepath.Dir(binPath)})
			}
		}
	}
	return versions
}

func detectSystem(key string) *versionInfo {
	exe := exeName(key)
	path, err := exec.LookPath(exe)
	if err != nil {
		// Try alternate names
		switch key {
		case "python":
			path, err = exec.LookPath("python3" + exeSuffix())
		case "php":
			path, err = exec.LookPath("php" + exeSuffix())
		case "node":
			path, err = exec.LookPath("node" + exeSuffix())
		case "go":
			path, err = exec.LookPath("go" + exeSuffix())
		}
		if err != nil {
			return nil
		}
	}
	realPath, _ := filepath.EvalSymlinks(path)
	v := getExeVersion(realPath, key)
	if v == "" {
		return nil
	}
	source := "System"
	return &versionInfo{version: v, source: source, path: filepath.Dir(realPath)}
}

func getExeVersion(exePath, key string) string {
	var args []string
	switch key {
	case "python":
		args = []string{"--version"}
	case "php":
		args = []string{"-v"}
	case "node":
		args = []string{"--version"}
	case "go":
		args = []string{"version"}
	}
	cmd := exec.Command(exePath, args...)
	out, _ := cmd.Output()
	return extractVersion(string(out), key)
}

func extractVersion(output, key string) string {
	var m []string
	switch key {
	case "python":
		m = pyVer.FindStringSubmatch(output)
	case "php":
		m = phpVer.FindStringSubmatch(output)
	case "node":
		m = nodeVer.FindStringSubmatch(output)
	case "go":
		m = goVer.FindStringSubmatch(output)
	}
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func exeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

// ---- Activation ----
func activateVersion(key string, v versionInfo) {
	dir := filepath.Join(svDir, "runtimes", key, v.version)
	if v.source == "System" {
		dir = v.path
	}
	addToPath(dir, key)

	switch key {
	case "python":
		cfg.Python = v.version
	case "php":
		cfg.PHP = v.version
	case "node":
		cfg.Node = v.version
	case "go":
		cfg.Go = v.version
	}
	saveConfig()
	os.Setenv("SV_"+strings.ToUpper(key), v.version)
	fmt.Println()
	fmt.Println("  " + trFmt("version_set", runtimeLabel(key), v.version))
}

func selectVersion(key string) {
	for {
		versions := listVersions(key)
		if len(versions) == 0 {
			fmt.Println("\n  " + tr("no_versions"))
			pause()
			return
		}
		items := make([]string, len(versions))
		for i, v := range versions {
			label := v.version + "  [" + v.source + "]"
			if len(v.path) > 40 {
				label += " .." + v.path[len(v.path)-38:]
			} else {
				label += " " + v.path
			}
			items[i] = label
		}
		sel := menu(trFmt("select_version", strings.ToUpper(runtimeLabel(key))), items)
		if sel < 0 {
			return
		}
		activateVersion(key, versions[sel])
		pause()
		return
	}
}

func findByPrefix(versions []versionInfo, prefix string) *versionInfo {
	prefix = strings.ToLower(prefix)
	for _, v := range versions {
		if strings.HasPrefix(strings.ToLower(v.version), prefix) {
			return &v
		}
	}
	return nil
}

func compareVersions(a, b string) int {
	pa := parseVersion(a)
	pb := parseVersion(b)
	for i := 0; i < 3; i++ {
		if i >= len(pa) && i >= len(pb) {
			return 0
		}
		if i >= len(pa) {
			return -1
		}
		if i >= len(pb) {
			return 1
		}
		if pa[i] != pb[i] {
			return pa[i] - pb[i]
		}
	}
	return 0
}

func parseVersion(v string) []int {
	v = strings.TrimLeft(v, "vV")
	parts := strings.Split(v, ".")
	var nums []int
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			break
		}
		nums = append(nums, n)
	}
	return nums
}

// ---- Profiles ----
type profileData struct {
	Python string `json:"python"`
	PHP    string `json:"php"`
	Node   string `json:"node"`
	Go     string `json:"go"`
}

func profileDir() string {
	return filepath.Join(svDir, "profiles")
}

func saveProfile(name string) {
	p := profileData{
		Python: cfg.Python,
		PHP:    cfg.PHP,
		Node:   cfg.Node,
		Go:     cfg.Go,
	}
	os.MkdirAll(profileDir(), 0755)
	data, _ := json.MarshalIndent(p, "", "  ")
	os.WriteFile(filepath.Join(profileDir(), name+".json"), data, 0644)
	fmt.Println("  " + trFmt("profile_saved", name))
}

func loadProfile(name string) {
	data, err := os.ReadFile(filepath.Join(profileDir(), name+".json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %s\n", trFmt("profile_not_found", name))
		return
	}
	var p profileData
	json.Unmarshal(data, &p)
	if p.Python != "" {
		cfg.Python = p.Python
	}
	if p.PHP != "" {
		cfg.PHP = p.PHP
	}
	if p.Node != "" {
		cfg.Node = p.Node
	}
	if p.Go != "" {
		cfg.Go = p.Go
	}
	saveConfig()
	fmt.Println("  " + trFmt("profile_loaded", name))
}

func listProfiles() {
	entries, err := os.ReadDir(profileDir())
	if err != nil || len(entries) == 0 {
		fmt.Println("  " + tr("no_profiles"))
		return
	}
	fmt.Println("  " + tr("profiles_list"))
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			fmt.Printf("    %s\n", strings.TrimSuffix(e.Name(), ".json"))
		}
	}
}

func deleteProfile(name string) {
	p := filepath.Join(profileDir(), name+".json")
	if err := os.Remove(p); err != nil {
		fmt.Fprintf(os.Stderr, "  %s\n", trFmt("profile_not_found", name))
		return
	}
	fmt.Println("  " + trFmt("profile_deleted", name))
}

func profileMenu() {
	items := []string{tr("save"), tr("load"), tr("list"), tr("delete"), tr("back")}
	for {
		sel := menu(tr("profiles"), items)
		if sel < 0 || sel == 4 {
			return
		}
		switch sel {
		case 0:
			name := input(tr("profile_name"))
			if name != "" {
				saveProfile(name)
				pause()
			}
		case 1:
			name := input(tr("profile_name_load"))
			if name != "" {
				loadProfile(name)
				pause()
			}
		case 2:
			listProfiles()
			pause()
		case 3:
			name := input(tr("profile_name_delete"))
			if name != "" {
				deleteProfile(name)
				pause()
			}
		}
	}
}

// ---- .pivotrc ----
func cmdInit() {
	cwd, _ := os.Getwd()
	rcPath := filepath.Join(cwd, ".pivotrc")
	if _, err := os.Stat(rcPath); err == nil {
		fmt.Println("  " + trFmt("pivotrc_exists", cwd))
		return
	}
	content := "# pivot runtime config\n"
	if cfg.Python != "" {
		content += "python=" + cfg.Python + "\n"
	}
	if cfg.PHP != "" {
		content += "php=" + cfg.PHP + "\n"
	}
	if cfg.Node != "" {
		content += "node=" + cfg.Node + "\n"
	}
	if cfg.Go != "" {
		content += "go=" + cfg.Go + "\n"
	}
	os.WriteFile(rcPath, []byte(content), 0644)
	fmt.Println("  " + trFmt("pivotrc_created", cwd))
}
