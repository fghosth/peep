package main

import (
	"github.com/fghosth/peep/cmd/cmd"
)
var (
	Version   string
	BuildTime string
	GoVersion string
)

func main() {
	cmd.GoVersion = GoVersion
	cmd.BuildTime = BuildTime
	cmd.Version = Version
	cmd.Execute()
}
