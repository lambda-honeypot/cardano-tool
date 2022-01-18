package cmdrunner

import (
	"os/exec"
)

//CmdRunner is the interface implementation that runs commands on the command line
type CmdRunner struct {
	Command string
}

//ExecuteCommand execute a command on the command line and prints the output
func (r CmdRunner) ExecuteCommand(args ...string) ([]byte, error) {
	trimmedArgs := deleteEmpty(args)
	aCmd := exec.Command(r.Command, trimmedArgs...)
	stdout, err := aCmd.CombinedOutput()
	return stdout, err
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
