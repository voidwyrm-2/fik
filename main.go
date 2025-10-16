/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/voidwyrm-2/fik/cmd"
)

//go:embed version.txt
var version string

func main() {
	if err := cmd.Execute(strings.TrimSpace(version)); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
