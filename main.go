package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var version = "dev"

func main() {
	detectLang()

	home, _ := os.UserHomeDir()
	svDir = filepath.Join(home, ".pivot")
	ensureDir(svDir, "runtimes", "projects", "profiles", "bin")

	loadConfig(svDir)
	checkAutoApply()

	args := os.Args[1:]
	if len(args) == 0 {
		runInteractive()
		return
	}

	switch args[0] {
	case "list":
		cmdList()
	case "use":
		cmdUse(args[1:])
	case "install":
		cmdInstall(args[1:])
	case "create":
		cmdCreate(args[1:])
	case "profile":
		cmdProfile(args[1:])
	case "init":
		cmdInit()
	case "run":
		cmdRun(args[1:])
	case "doctor":
		cmdDoctor()
	case "upgrade":
		cmdUpgrade()
	case "clean":
		cmdClean()
	case "update":
		checkUpdates(svDir)
	case "env":
		printEnv()
	case "version", "--version", "-v":
		fmt.Println("pivot version", version)
	default:
		fmt.Fprintf(os.Stderr, "Usage: pivot <command> [args]\n\n"+tr("commands")+":\n")
		fmt.Fprintf(os.Stderr, "  list              %s\n", tr("show_installed"))
		fmt.Fprintf(os.Stderr, "  use <rt> <ver>    %s\n", tr("activate_version"))
		fmt.Fprintf(os.Stderr, "  install <rt>      %s\n", tr("download_runtime"))
		fmt.Fprintf(os.Stderr, "  create <fw> <name> %s\n", tr("create_framework"))
		fmt.Fprintf(os.Stderr, "  profile <op> <name> %s\n", tr("manage_profiles"))
		fmt.Fprintf(os.Stderr, "  init              %s\n", tr("create_pivotrc"))
		fmt.Fprintf(os.Stderr, "  run <rt> <ver> <cmd> %s\n", tr("run_with_version"))
		fmt.Fprintf(os.Stderr, "  doctor            %s\n", tr("diagnose_system"))
		fmt.Fprintf(os.Stderr, "  upgrade           %s\n", tr("upgrade_self"))
		fmt.Fprintf(os.Stderr, "  clean             %s\n", tr("clean_runtimes"))
		fmt.Fprintf(os.Stderr, "  update            %s\n", tr("check_new_versions"))
		fmt.Fprintf(os.Stderr, "  env               %s\n", tr("print_path_setup"))
		os.Exit(1)
	}
}

func ensureDir(base string, dirs ...string) {
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(base, d), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create %s: %v\n", d, err)
		}
	}
}

func checkAutoApply() {
	cwd, _ := os.Getwd()
	p := cwd
	for {
		rc := filepath.Join(p, ".pivotrc")
		if _, err := os.Stat(rc); err == nil {
			data, err := os.ReadFile(rc)
			if err != nil {
				return
			}
			for _, line := range parseEnvLines(string(data)) {
				switch line.key {
				case "python":
					if cfg.Python == "" {
						cfg.Python = line.val
					}
				case "php":
					if cfg.PHP == "" {
						cfg.PHP = line.val
					}
				case "node":
					if cfg.Node == "" {
						cfg.Node = line.val
					}
				case "go":
					if cfg.Go == "" {
						cfg.Go = line.val
					}
				}
			}
			return
		}
		parent := filepath.Dir(p)
		if parent == p {
			return
		}
		p = parent
	}
}

type envLine struct {
	key string
	val string
}

func parseEnvLines(content string) []envLine {
	var lines []envLine
	for _, raw := range strings.Split(content, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			lines = append(lines, envLine{key: strings.TrimSpace(parts[0]), val: strings.TrimSpace(parts[1])})
		}
	}
	return lines
}

func runInteractive() {
	for {
		items := []string{
			trFmt("python_label", orDash(activePython())),
			trFmt("php_label", orDash(activePHP())),
			trFmt("nodejs_label", orDash(activeNode())),
			trFmt("go_label", orDash(activeGo())),
			tr("create_project"),
			tr("profiles"),
			tr("check_updates"),
			tr("exit"),
		}
		sel := menu(tr("version_switcher"), items)
		if sel < 0 {
			break
		}
		switch sel {
		case 0:
			selectVersion("python")
		case 1:
			selectVersion("php")
		case 2:
			selectVersion("node")
		case 3:
			selectVersion("go")
		case 4:
			createProject()
		case 5:
			profileMenu()
		case 6:
			checkUpdates(svDir)
			pause()
		case 7:
			return
		}
	}
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func cmdList() {
	for _, r := range []string{"python", "php", "node", "go"} {
		fmt.Printf("\n  %s:\n", runtimeLabel(r))
		for _, v := range listVersions(r) {
			fmt.Printf("    %s  [%s] %s\n", v.version, v.source, v.path)
		}
	}
}

func cmdUse(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, tr("usage_use"))
		os.Exit(1)
	}
	vs := listVersions(args[0])
	v := findByPrefix(vs, args[1])
	if v == nil {
		fmt.Fprintf(os.Stderr, tr("version_not_found"), args[1], args[0])
		os.Exit(1)
	}
	activateVersion(args[0], *v)
}

func cmdInstall(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, tr("usage_install"))
		os.Exit(1)
	}
	ver := ""
	if len(args) > 1 {
		ver = args[1]
	}
	if err := downloadRuntime(args[0], ver); err != nil {
		fmt.Fprintf(os.Stderr, tr("install_failed"), err)
		os.Exit(1)
	}
}

func cmdCreate(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, tr("usage_create"))
		os.Exit(1)
	}
	if err := scaffoldProject(args[0], args[1]); err != nil {
		fmt.Fprintf(os.Stderr, tr("create_failed"), err)
		os.Exit(1)
	}
}

func cmdProfile(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, tr("usage_profile"))
		os.Exit(1)
	}
	switch args[0] {
	case "save":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: pivot profile save <name>")
			os.Exit(1)
		}
		saveProfile(args[1])
	case "load":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: pivot profile load <name>")
			os.Exit(1)
		}
		loadProfile(args[1])
	case "list":
		listProfiles()
	case "delete":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: pivot profile delete <name>")
			os.Exit(1)
		}
		deleteProfile(args[1])
	default:
		fmt.Fprintln(os.Stderr, tr("usage_profile"))
		os.Exit(1)
	}
}
