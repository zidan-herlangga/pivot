package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type framework struct {
	name string
	pkg  string
	ver  string
	typ  string
	tmpl string
}

var frameworks = []framework{
	{"Laravel", "laravel/laravel", "^11.0", "composer", ""},
	{"CodeIgniter 4", "codeigniter4/app-project", "^4.0", "composer", ""},
	{"Symfony", "symfony/skeleton", "^7.0", "composer", ""},
	{"WordPress (Bedrock)", "roots/bedrock", "^2.0", "composer", ""},
	{"React (Vite)", "create-vite", "latest", "npm", "react"},
	{"Next.js", "create-next-app", "latest", "npm", ""},
	{"Vue (Vite)", "create-vite", "latest", "npm", "vue"},
	{"AdonisJS", "create-adonisjs", "latest", "npm", ""},
}

func scaffoldProject(slug, name string) error {
	for _, fw := range frameworks {
		slugLower := strings.ToLower(slug)
		fwLower := strings.ToLower(strings.Split(fw.name, " ")[0])
		if slugLower == fwLower || slugLower == strings.ToLower(fw.name) {
			cwd, _ := os.Getwd()
			projectDir := filepath.Join(cwd, name)
			return doScaffold(projectDir, fw)
		}
	}
	return fmt.Errorf(tr("framework_unknown"), slug)
}

func doScaffold(dir string, fw framework) error {
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf(tr("folder_exists"), dir)
	}
	fmt.Printf("  %s %s...\n", fw.name, tr("creating"))
	os.MkdirAll(dir, 0755)

	var cmd *exec.Cmd
	switch fw.typ {
	case "composer":
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", "composer", "create-project", fw.pkg+":"+fw.ver, dir, "--no-interaction", "--prefer-dist")
		} else {
			cmd = exec.Command("composer", "create-project", fw.pkg+":"+fw.ver, dir, "--no-interaction", "--prefer-dist")
		}
	case "npm":
		switch fw.pkg {
		case "create-vite":
			if fw.tmpl != "" {
				cmd = exec.Command("npx", "create-vite@latest", dir, "--template", fw.tmpl)
			} else {
				cmd = exec.Command("npx", "create-vite@latest", dir)
			}
		case "create-next-app":
			cmd = exec.Command("npx", "create-next-app@latest", dir)
		case "create-adonisjs":
			cmd = exec.Command("npm", "init", "adonisjs@latest", dir)
		}
	}

	if cmd == nil {
		return fmt.Errorf("unsupported project type: %s", fw.typ)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scaffold failed: %w", err)
	}

	fmt.Println("  " + trFmt("project_created", filepath.Base(dir)))
	return nil
}

func createProject() {
	items := make([]string, len(frameworks))
	for i, f := range frameworks {
		items[i] = f.name
	}
	sel := menu(tr("create_project"), items)
	if sel < 0 {
		return
	}
	fw := frameworks[sel]
	name := input(trFmt("project_name_for", fw.name))
	if name == "" {
		return
	}
	cwd, _ := os.Getwd()
	projectDir := filepath.Join(cwd, name)
	if err := doScaffold(projectDir, fw); err != nil {
		fmt.Fprintf(os.Stderr, tr("create_failed"), err)
	}
	pause()
}
