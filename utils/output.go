package utils

import (
	"bytes"
	"os"
)

// Buffered to optimize
func PotetoPrint(msg string) {
	buf := bytes.NewBuffer([]byte(msg))
	buf.WriteTo(os.Stdout)
}
