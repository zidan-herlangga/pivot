package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/term"
)

func menu(title string, items []string) int {
	sel := 0
	w := 48

	for {
		clear()
		fmt.Println()
		fmt.Print("  " + strings.Repeat("=", w) + "\n")
		pad := (w - len(title) - 2) / 2
		if pad < 0 {
			pad = 0
		}
		rpad := w - pad - len(title) - 2
		if rpad < 0 {
			rpad = 0
		}
		fmt.Print("  " + strings.Repeat(" ", pad) + "  " + title + strings.Repeat(" ", rpad) + "  \n")
		fmt.Print("  " + strings.Repeat("=", w) + "\n")
		for i, item := range items {
			if i == sel {
				n := w - len(item) - 4
				if n < 0 {
					n = 0
				}
				fmt.Print("\033[7m  > " + item + strings.Repeat(" ", n) + "\033[0m\n")
			} else {
				n := w - len(item) - 5
				if n < 0 {
					n = 0
				}
				fmt.Print("    " + item + strings.Repeat(" ", n) + "\n")
			}
		}
		fmt.Print("  " + strings.Repeat("=", w) + "\n")
		fmt.Print("    ↓/↑  Enter  Esc\n")

		key := readKey()
		switch key {
		case "up":
			sel = (sel - 1 + len(items)) % len(items)
		case "down":
			sel = (sel + 1) % len(items)
		case "enter":
			return sel
		case "esc":
			return -1
		}
	}
}

func input(prompt string) string {
	fmt.Print("\n  " + prompt + "\n  > ")
	var s string
	fmt.Scanln(&s)
	return strings.TrimSpace(s)
}

func pause() {
	fmt.Print("\n  " + tr("press_any_key"))
	readKey()
}

func clear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[2J\033[H")
	}
}

func readKey() string {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		var b [1]byte
		os.Stdin.Read(b[:])
		return string(b[:])
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		var b [1]byte
		os.Stdin.Read(b[:])
		return string(b[:])
	}
	defer term.Restore(fd, oldState)

	var buf [3]byte
	n, _ := os.Stdin.Read(buf[:])
	if n == 1 {
		switch buf[0] {
		case 13:
			return "enter"
		case 27:
			return "esc"
		case 3:
			os.Exit(0)
		}
		return string(buf[0])
	}
	if n == 3 && buf[0] == 27 && buf[1] == 91 {
		switch buf[2] {
		case 65:
			return "up"
		case 66:
			return "down"
		case 67:
			return "right"
		case 68:
			return "left"
		}
	}
	return ""
}
