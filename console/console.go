package console

import (
	"bufio"
	"fmt"
)

func CursorAbs(x int, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}
func CursorUp(n int)   { fmt.Printf("\033[%dA", n) }
func CursorDown(n int) { fmt.Printf("\033[%dB", n) }
func CursorFor(n int)  { fmt.Printf("\033[%dC", n) }
func CursorBack(n int) { fmt.Printf("\033[%dD", n) }

func EraseLine() { fmt.Print("\r\033[K") }

func Write(mes string)                          { fmt.Print(mes) }
func Writeln(mes string)                        { fmt.Println(mes) }
func Writef(format string, args ...interface{}) { fmt.Printf(format, args...) }

var sc *bufio.Scanner

func SetScanner(scanner *bufio.Scanner) {
	sc = scanner
}

func Prompt(mes string) string {
	Writef("%s", mes)
	return Scan()
}

func Scan() string {
	if sc.Scan() {
		return sc.Text()
	}
	return ""
}
