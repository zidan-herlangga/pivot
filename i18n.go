package main

import (
	"fmt"
	"os"
	"strings"
)

var lang string

func detectLang() {
	lang = "en"
	for _, e := range []string{"LANG", "LC_ALL", "LC_MESSAGES"} {
		if v := os.Getenv(e); v != "" {
			if strings.HasPrefix(strings.ToLower(v), "id") || strings.HasPrefix(strings.ToLower(v), "in") {
				lang = "id"
				return
			}
		}
	}
}

func tr(key string, args ...interface{}) string {
	m, ok := i18n[key]
	if !ok {
		return key
	}
	s := m[0]
	if lang == "id" && len(m) > 1 {
		s = m[1]
	}
	if len(args) > 0 {
		return fmt.Sprintf(s, args...)
	}
	return s
}

func trFmt(key string, args ...interface{}) string {
	return tr(key, args...)
}

var i18n = map[string][2]string{
	// Menu
	"version_switcher":  {"VERSION SWITCHER", "PEMILIH VERSI"},
	"python":            {"Python", "Python"},
	"php":               {"PHP", "PHP"},
	"nodejs":            {"Node.js", "Node.js"},
	"go":                {"Go", "Go"},
	"create_project":    {"Create Project", "Buat Project"},
	"profiles":          {"Profiles", "Profil"},
	"check_updates":     {"Check Updates", "Cek Update"},
	"exit":              {"Exit", "Keluar"},
	"back":              {"Back", "Kembali"},

	// Version selection
	"no_versions":       {"No versions found.", "Tidak ada versi ditemukan."},
	"version_set":       {"%s %s set as default.", "%s %s ditetapkan sebagai bawaan."},
	"select_version":    {"%s VERSIONS", "%s VERSI"},
	"usage_use":         {"Usage: pivot use <runtime> <version>", "Penggunaan: pivot use <runtime> <version>"},
	"usage_install":     {"Usage: pivot install <runtime> [version]", "Penggunaan: pivot install <runtime> [version]"},
	"usage_create":      {"Usage: pivot create <framework> <name>", "Penggunaan: pivot create <framework> <nama>"},
	"usage_profile":     {"Usage: pivot profile save|load|list|delete <name>", "Penggunaan: pivot profile simpan|muat|daftar|hapus <nama>"},
	"version_not_found": {"Version '%s' not found for %s", "Versi '%s' tidak ditemukan untuk %s"},
	"unknown_runtime":   {"Unknown runtime: %s (use python, php, node, go)", "Runtime tidak dikenal: %s (gunakan python, php, node, go)"},

	// Install
	"already_downloaded": {"%s %s already downloaded.", "%s %s sudah diunduh."},
	"downloading":        {"Downloading %s %s ...", "Mengunduh %s %s ..."},
	"extracting":         {"Extracting ...", "Mengekstrak ..."},
	"download_failed":    {"download failed: %s", "unduh gagal: %s"},
	"extract_failed":     {"extract failed: %s", "ekstrak gagal: %s"},
	"install_failed":     {"Install failed: %v", "Instalasi gagal: %v"},

	// Project creation
	"project_created":     {"Project '%s' created!", "Project '%s' dibuat!"},
	"creating":            {"creating project...", "membuat project..."},
	"create_failed":       {"Create failed: %v", "Pembuatan gagal: %v"},
	"folder_exists":       {"Folder '%s' already exists.", "Folder '%s' sudah ada."},
	"project_name_for":    {"Project name for %s", "Nama project untuk %s"},
	"profile_name":        {"Profile name", "Nama profil"},
	"profile_name_load":   {"Profile name to load", "Nama profil yang akan dimuat"},
	"profile_name_delete": {"Profile name to delete", "Nama profil yang akan dihapus"},
	"framework_unknown":   {"Unknown framework: %s", "Framework tidak dikenal: %s"},

	// Profiles
	"profile_saved":      {"Profile '%s' saved.", "Profil '%s' disimpan."},
	"profile_loaded":     {"Profile '%s' loaded.", "Profil '%s' dimuat."},
	"profile_not_found":  {"Profile '%s' not found.", "Profil '%s' tidak ditemukan."},
	"profile_deleted":    {"Profile '%s' deleted.", "Profil '%s' dihapus."},
	"no_profiles":        {"No profiles saved.", "Belum ada profil disimpan."},
	"profiles_list":      {"Profiles:", "Profil:"},

	// .pivotrc
	"pivotrc_exists":   {".pivotrc already exists in %s", ".pivotrc sudah ada di %s"},
	"pivotrc_created":  {".pivotrc created in %s", ".pivotrc dibuat di %s"},

	// Update
	"checking_updates":   {"Checking for updates...", "Memeriksa update..."},
	"update_failed":      {"Failed to check updates.", "Gagal memeriksa update."},
	"latest_versions":    {"Latest available versions:", "Versi terbaru yang tersedia:"},

	// List
	"installed":          {"Installed versions:", "Versi terpasang:"},

	// Misc
	"press_any_key":      {"Press any key...", "Tekan tombol apa saja..."},
	"save":               {"Save Profile", "Simpan Profil"},
	"load":               {"Load Profile", "Muat Profil"},
	"list":               {"List Profiles", "Daftar Profil"},
	"delete":             {"Delete Profile", "Hapus Profil"},

	// Interactive menu labels
	"python_label":       {"Python  [%s]", "Python  [%s]"},
	"php_label":          {"PHP     [%s]", "PHP     [%s]"},
	"nodejs_label":       {"Node.js [%s]", "Node.js [%s]"},
	"go_label":           {"Go [%s]", "Go [%s]"},

	// Help text
	"commands":           {"Commands", "Perintah"},
	"show_installed":     {"Show installed versions", "Tampilkan versi terpasang"},
	"activate_version":   {"Activate a runtime version", "Aktifkan versi runtime"},
	"download_runtime":   {"Download a runtime", "Unduh runtime"},
	"create_framework":   {"Create a framework project", "Buat project framework"},
	"manage_profiles":    {"Manage profiles", "Kelola profil"},
	"create_pivotrc":     {"Create .pivotrc in current dir", "Buat .pivotrc di direktori ini"},
	"check_new_versions": {"Check for newer versions", "Periksa versi terbaru"},
	"print_path_setup":   {"Print PATH setup for shell", "Tampilkan pengaturan PATH untuk shell"},
}
