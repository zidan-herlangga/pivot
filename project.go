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
	// PHP
	{"Laravel", "laravel/laravel", "^11.0", "composer", ""},
	{"CodeIgniter 4", "codeigniter4/app-project", "^4.0", "composer", ""},
	{"Symfony", "symfony/skeleton", "^7.0", "composer", ""},
	{"WordPress (Bedrock)", "roots/bedrock", "^2.0", "composer", ""},

	// JavaScript / TypeScript
	{"React (Vite)", "create-vite", "latest", "npm", "react"},
	{"Next.js", "create-next-app", "latest", "npm", ""},
	{"Vue (Vite)", "create-vite", "latest", "npm", "vue"},
	{"AdonisJS", "create-adonisjs", "latest", "npm", ""},
	{"Svelte (Vite)", "create-vite", "latest", "npm", "svelte"},
	{"Nuxt", "create-nuxt-app", "latest", "npm", ""},
	{"Solid (Vite)", "create-vite", "latest", "npm", "solid"},

	// Python
	{"Django", "django", "", "pip", ""},
	{"Flask", "flask", "", "pip", ""},
	{"FastAPI", "fastapi", "", "pip", ""},

	// Go
	{"Gin", "gin-gonic/gin", "", "go", ""},
	{"Echo", "labstack/echo", "", "go", ""},
	{"Fiber", "gofiber/fiber", "", "go", ""},

	// Ruby
	{"Ruby on Rails", "rails", "", "rails", ""},

	// Java
	{"Spring Boot", "spring-boot", "", "spring", ""},
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
		switch fw.typ {
		case "pip":
			return fmt.Errorf("'python' not found in PATH — install Python first")
		case "go":
			return fmt.Errorf("'go' not found in PATH — install Go first")
		case "rails":
			return fmt.Errorf("'rails' not found in PATH — run: gem install rails")
		case "spring":
			return fmt.Errorf("'mvn' or 'gradle' not found in PATH — install Maven/Gradle first")
		}
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
		case "create-nuxt-app":
			return "npx", []string{"create-nuxt-app@latest", dir}
		}
	case "pip":
		switch fw.pkg {
		case "django":
			return "python", []string{"-m", "django", "startproject", filepath.Base(dir)}
		case "flask":
			return "python", []string{"-m", "flask", "new", dir}
		case "fastapi":
			return "python", []string{"-m", "fastapi", "new", dir}
		}
	case "go":
		switch fw.pkg {
		case "gin-gonic/gin":
			return "go", []string{"mod", "init", filepath.Base(dir)}
		case "labstack/echo":
			return "go", []string{"mod", "init", filepath.Base(dir)}
		case "gofiber/fiber":
			return "go", []string{"mod", "init", filepath.Base(dir)}
		}
	case "rails":
		return "rails", []string{"new", dir}
	case "spring":
		_, mvnErr := exec.LookPath("mvn")
		_, gradleErr := exec.LookPath("gradle")
		if mvnErr == nil {
			return "mvn", []string{"archetype:generate", "-DgroupId=com.example", "-DartifactId=" + filepath.Base(dir), "-DarchetypeArtifactId=maven-archetype-quickstart", "-DinteractiveMode=false"}
		}
		if gradleErr == nil {
			return "gradle", []string{"init", "--project-name", filepath.Base(dir), "--type", "java-application"}
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
