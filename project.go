package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	exe, args := scaffoldCmd(dir, fw)
	if exe == "" {
		return fmt.Errorf("unsupported project type: %s", fw.typ)
	}

	if _, err := exec.LookPath(exe); err != nil {
		return fmt.Errorf("'%s' not found in PATH — install it first", exe)
	}

	fmt.Printf("  %s %s...\n", fw.name, tr("creating"))
	os.MkdirAll(dir, 0755)

	cmd := exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return fmt.Errorf("scaffold failed: %w", err)
	}

	fmt.Println("  " + trFmt("project_created", filepath.Base(dir)))
	return nil
}

func scaffoldCmd(dir string, fw framework) (string, []string) {
	switch fw.typ {
	case "composer":
		return "composer", []string{"create-project", fw.pkg + ":" + fw.ver, dir, "--no-interaction", "--prefer-dist"}
	case "npm":
		switch fw.pkg {
		case "create-vite":
			if fw.tmpl != "" {
				return "npx", []string{"create-vite@latest", dir, "--template", fw.tmpl}
			}
			return "npx", []string{"create-vite@latest", dir}
		case "create-next-app":
			return "npx", []string{"create-next-app@latest", dir}
		case "create-adonisjs":
			return "npm", []string{"init", "adonisjs@latest", dir}
		}
	}
	return "", nil
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
