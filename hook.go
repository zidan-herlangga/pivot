package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func cmdHook() {
	shell, _ := os.LookupEnv("SHELL")
	if shell == "" {
		shell = "bash"
	}

	name := filepath.Base(shell)
	binDir := filepath.Join(svDir, "bin")

	switch {
	case strings.Contains(name, "fish"):
		fmt.Print(`function __pivot_cd_hook --on-variable PWD
  if status --is-interactive
    if test -f .pivotrc
      while read -lat line
        set -l parts (string split "=" -- $line)
        test -z "$parts[1]"; and continue
        string match -q "#*" "$parts[1]"; and continue
        pivot use $parts[1] $parts[2] 2>/dev/null
      end < .pivotrc
    end
  end
end
`)
	case strings.Contains(name, "zsh"):
		fmt.Printf(`_pivot_hook() {
  if [[ -f .pivotrc ]]; then
    while IFS='=' read -r key val; do
      [[ -z "$key" || "$key" == \#* ]] && continue
      pivot use "$key" "$val" 2>/dev/null
    done < .pivotrc
  fi
}
[[ -z "${precmd_functions[(r)_pivot_hook]}" ]] && precmd_functions+=(_pivot_hook)

export PATH="%s:$PATH"
`, binDir)
	default:
		fmt.Printf(`_pivot_hook() {
  if [[ -f .pivotrc ]]; then
    while IFS='=' read -r key val; do
      [[ -z "$key" || "$key" == \#* ]] && continue
      pivot use "$key" "$val" 2>/dev/null
    done < .pivotrc
  fi
}
[[ "$(type -t __pivot_hook)" != "function" ]] && PROMPT_COMMAND="_pivot_hook;$PROMPT_COMMAND"

export PATH="%s:$PATH"
`, binDir)
	}
}
