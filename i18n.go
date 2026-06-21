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
	"unknown_runtime":   {"Unknown runtime: %s", "Runtime tidak dikenal: %s"},
	"valid_runtimes":    {"Valid runtimes: python, php, node, go, deno, bun, java, rust", "Runtime valid: python, php, node, go, deno, bun, java, rust"},

	// Install
	"already_downloaded": {"%s %s already downloaded.", "%s %s sudah diunduh."},
	"downloading":        {"Downloading %s %s ...", "Mengunduh %s %s ..."},
	"extracting":         {"Extracting ...", "Mengekstrak ..."},
	"download_failed":    {"download failed: %s", "unduh gagal: %s"},
	"extract_failed":     {"extract failed: %s", "ekstrak gagal: %s"},
	"install_failed":     {"Install failed: %v", "Instalasi gagal: %v"},
	"no_download_for_platform": {"%s download is not available for %s", "Unduhan %s tidak tersedia untuk %s"},

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

	// Run
	"usage_run":          {"Usage: pivot run <runtime> <version> <command> [args...]", "Penggunaan: pivot run <runtime> <version> <command> [argumen...]"},

	// Doctor
	"doctor_path_ok":      {"%s is in PATH", "%s ada di PATH"},
	"doctor_path_missing": {"%s is NOT in PATH — add it with: pivot env", "%s TIDAK ada di PATH — tambahkan dengan: pivot env"},
	"doctor_no_versions":  {"No %s versions found", "Tidak ada versi %s ditemukan"},
	"doctor_conflict":     {"Multiple system %s detected —可能會导致冲突", "Multiple system %s terdeteksi — mungkin konflik"},
	"doctor_bin_missing":  {"%s binary not found at %s — re-run: pivot use", "Binary %s tidak ditemukan di %s — jalankan ulang: pivot use"},
	"doctor_active_missing": {"%s %s is active but version not found", "%s %s aktif tetapi versi tidak ditemukan"},
	"doctor_not_active":   {"No active %s version — run: pivot use", "Tidak ada versi %s aktif — jalankan: pivot use"},
	"doctor_no_goroot":    {"GOROOT not set — Go may not work correctly", "GOROOT tidak diatur — Go mungkin tidak berfungsi dengan benar"},
	"doctor_all_good":     {"All checks passed!", "Semua pemeriksaan berhasil!"},
	"doctor_issues_found": {"Issues found — see above.", "Masalah ditemukan — lihat di atas."},

	// Interactive
	"install_runtime":  {"Install Runtime", "Pasang Runtime"},
	"version_for":      {"Version for %s (empty = default)", "Versi untuk %s (kosong = bawaan)"},

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
	"run_with_version":   {"Run a command with a specific version", "Jalankan perintah dengan versi tertentu"},
	"diagnose_system":    {"Diagnose system and PATH issues", "Diagnosa sistem dan masalah PATH"},
	"upgrade_self":       {"Upgrade pivot to latest version", "Perbarui pivot ke versi terbaru"},
	"upgrade_fetch_failed": {"Failed to fetch latest version", "Gagal mengambil versi terbaru"},
	"upgrade_latest":     {"Already at latest version (%s)", "Sudah versi terbaru (%s)"},
	"upgrade_unsupported": {"Unsupported platform", "Platform tidak didukung"},
	"upgrade_no_asset":   {"No matching binary found for this platform", "Binary tidak ditemukan untuk platform ini"},
	"upgrade_downloading": {"Downloading %s ...", "Mengunduh %s ..."},
	"upgrade_download_failed": {"Download failed", "Unduh gagal"},
	"upgrade_extract_failed": {"Extract failed", "Ekstrak gagal"},
	"upgrade_failed":     {"Upgrade failed", "Pembaruan gagal"},
	"upgrade_done":       {"Upgraded to %s!", "Diperbarui ke %s!"},
	"clean_runtimes":     {"Remove unused runtime versions", "Hapus versi runtime yang tidak dipakai"},
	"removed":            {"Removed", "Dihapus"},
	"nothing_to_clean":   {"Nothing to clean.", "Tidak ada yang perlu dibersihkan."},
	"cleaned_versions":   {"Removed %d version(s).", "%d versi dihapus."},
	"checking_upgrades":  {"Checking for upgrades...", "Memeriksa pembaruan..."},

	// Shell
	"usage_shell":        {"Usage: pivot shell <runtime> <version>", "Penggunaan: pivot shell <runtime> <version>"},
	"shell_version":      {"Spawn a shell with specific version active", "Buka shell dengan versi tertentu"},

	// Hook
	"hook_info":          {"Print shell hook for auto-apply .pivotrc", "Cetak hook shell untuk auto-apply .pivotrc"},

	// Completion
	"usage_completion":   {"Usage: pivot completion bash|zsh|fish", "Penggunaan: pivot completion bash|zsh|fish"},
	"completion_info":    {"Generate shell completion script", "Hasilkan skrip completion shell"},

	// Doctor --fix
	"doctor_fixed_path":  {"Fixed: added %s to PATH", "Diperbaiki: %s ditambahkan ke PATH"},
	"doctor_fixed_bin":   {"Fixed: re-linked %s binary", "Diperbaiki: binary %s ditautkan ulang"},

	// Extras
	"deno_label":         {"Deno    [%s]", "Deno    [%s]"},
	"bun_label":          {"Bun     [%s]", "Bun     [%s]"},
	"java_label":         {"Java    [%s]", "Java    [%s]"},
	"rust_label":         {"Rust    [%s]", "Rust    [%s]"},
}
