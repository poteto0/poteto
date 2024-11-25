package utils

import (
	"bytes"
	"os"
)

func PotetoPrint(msg string) {
	buf := bytes.NewBuffer([]byte(msg))
	buf.WriteTo(os.Stdout)
}
