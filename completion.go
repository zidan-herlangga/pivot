package main

import (
	"fmt"
	"os"
)

func cmdCompletion(args []string) {
	switch {
	case len(args) > 0 && args[0] == "bash":
		printBashCompletion()
	case len(args) > 0 && args[0] == "zsh":
		printZshCompletion()
	case len(args) > 0 && args[0] == "fish":
		printFishCompletion()
	default:
		fmt.Fprintln(os.Stderr, tr("usage_completion"))
		os.Exit(1)
	}
}

func printBashCompletion() {
	fmt.Print(`_pivot() {
  local cur="${COMP_WORDS[COMP_CWORD]}"
  local prev="${COMP_WORDS[COMP_CWORD-1]}"
  local cmds="list use install run create profile init doctor upgrade clean update env hook completion shell pin"
  local runtimes="python php node go deno bun java rust"

  if [[ $COMP_CWORD -eq 1 ]]; then
    COMPREPLY=($(compgen -W "$cmds" -- "$cur"))
  elif [[ $COMP_CWORD -eq 2 ]]; then
    case "$prev" in
      use|install|run) COMPREPLY=($(compgen -W "$runtimes" -- "$cur")) ;;
      create) COMPREPLY=($(compgen -W "laravel codeigniter symfony wordpress react nextjs vue adonisjs" -- "$cur")) ;;
      profile) COMPREPLY=($(compgen -W "save load list delete" -- "$cur")) ;;
      completion) COMPREPLY=($(compgen -W "bash zsh fish" -- "$cur")) ;;
    esac
  elif [[ $COMP_CWORD -eq 3 ]]; then
    case "$prev" in
      python|php|node|go|deno|bun|java|rust) COMPREPLY=($(compgen -W "latest system lts" -- "$cur")) ;;
    esac
  fi
}
complete -F _pivot pivot
`)
}

func printZshCompletion() {
	fmt.Print(`#compdef pivot
_pivot() {
  local line
  _arguments -C \
    "1:command:(list use install run create profile init doctor upgrade clean update env hook completion shell pin)" \
    "*::arg:->args"
  case $line[1] in
    use|install|run) _arguments "2:runtime:(python php node go deno bun java rust)" ;;
    create) _arguments "2:framework:(laravel codeigniter symfony wordpress react nextjs vue adonisjs)" ;;
    profile) _arguments "2:action:(save load list delete)" ;;
    completion) _arguments "2:shell:(bash zsh fish)" ;;
  esac
}
_pivot "$@"
`)
}

func printFishCompletion() {
	fmt.Print(`complete -c pivot -f -a "list use install run create profile init doctor upgrade clean update env hook completion shell pin" -d "Commands"
complete -c pivot -n "__fish_seen_subcommand_from use install run" -a "python php node go deno bun java rust" -d "Runtime"
complete -c pivot -n "__fish_seen_subcommand_from create" -a "laravel codeigniter symfony wordpress react nextjs vue adonisjs" -d "Framework"
complete -c pivot -n "__fish_seen_subcommand_from profile" -a "save load list delete" -d "Actions"
complete -c pivot -n "__fish_seen_subcommand_from completion" -a "bash zsh fish" -d "Shell"
`)
}
