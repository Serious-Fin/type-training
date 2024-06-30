package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"


	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

func getCursorPos() (int, int, error) {
	// Switch terminal to raw mode to capture input without pressing Enter
	fd := int(os.Stdin.Fd())
	oldState := syscall.Termios{}
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&oldState)), 0, 0, 0); err != 0 {
		return 0, 0, err
	}
	newState := oldState
	newState.Lflag &^= syscall.ICANON | syscall.ECHO
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		return 0, 0, err
	}
	defer syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)

	// Send ANSI escape code to query the cursor position
	fmt.Print("\033[6n")

	// Read the response: ESC [ rows ; cols R
	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 32)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		if b == 'R' {
			break
		}
		buf = append(buf, b)
	}

	// Parse the response
	response := string(buf)
	if response[0] != '\033' || response[1] != '[' {
		return 0, 0, fmt.Errorf("unexpected response format")
	}
	parts := strings.Split(response[2:], ";")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected response format")
	}
	rows, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	cols, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return rows, cols, nil
}

// Move the cursor to the specified position
func moveCursor(x, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

// Clear the line from the cursor to the end
func clearLine() {
	fmt.Print("\033[K")
}

func main() {
	text := []rune("The quick brown fox jumps over the lazy dog")
	fmt.Println("Type the following text:")
	fmt.Println(string(text))
	fmt.Println()
  row, col, _ := getCursorPos()
	// Set terminal to raw mode to capture input without pressing Enter
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	reader := bufio.NewReader(os.Stdin)
	input := make([]rune, 0, len(text))
	cursorX, cursorY := col, row

	//moveCursor(cursorX, cursorY)

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}

		if r == '\r' || r == '\n' {
			break
		}

		input = append(input, r)
		printComparison(text, input, cursorX, cursorY)

		cursorX += runewidth.RuneWidth(r)
		moveCursor(cursorX, cursorY)
	}

	fmt.Println()
}

func printComparison(expected []rune, actual []rune, cursorX, cursorY int) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	moveCursor(cursorX, cursorY)
	clearLine()

	for i := 0; i < len(expected); i++ {
		if i < len(actual) && expected[i] == actual[i] {
			fmt.Print(green(string(expected[i])))
		} else if i < len(actual) && expected[i] != actual[i] {
			fmt.Print(red(string(actual[i])))
		} else {
			fmt.Print(string(expected[i]))
		}
	}
}

func moveCursorUp(lines int) {
	fmt.Printf("\033[%dA", lines)
}

func moveCursorToStart() {
	fmt.Print("\r")
}

