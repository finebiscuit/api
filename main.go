/*
Copyright © 2021 Patryk Kalinowski <patryk@kalinowski.dev>
This file is part of the Biscuit API.
*/
package main

import "github.com/finebiscuit/api/cmd"

var (
	version   string
	buildTime string
)

func main() {
	cmd.Execute()
}
