package pkg

import "os/exec"

type Package struct {
	OS   SystemType
	Path string
}

type Installer struct {
	pack *Package
	cmd  *exec.Cmd
}
