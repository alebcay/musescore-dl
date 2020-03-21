//go:generate go run build/generate.go

package main

import (
	"github.com/alebcay/musescore-dl/internal/cmd"
)

func main() {
	cmd.Execute()
}
