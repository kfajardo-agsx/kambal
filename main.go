package main

import "gitlab.com/amihan/core/base.git/cmd"

var version = "dev"

func main() {
	cmd.Execute(version)
}
