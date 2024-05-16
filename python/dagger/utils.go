package main

import (
	"strings"
)

func (p *Python) Container() *Container {
	return p.Ctr
}

// Directory Return the current working directory
func (p *Python) Directory() *Directory {
	return p.Ctr.Directory(workdir)
}

// ShDashC returns a shell single command that runs the provided commands
func ShDashC(commands []string) []string {
	commands = append([]string{"set -o xtrace"}, commands...)
	return []string{"sh", "-c", strings.Join(commands, " && ")}
}
