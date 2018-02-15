package workflow

import (
	"os"
	"os/exec"
)

func runTfCommand(buildPath string, args ...string) error {
	tfCommand := exec.Command("terraform", args...)
	tfCommand.Dir = buildPath
	tfCommand.Stdin = os.Stdin
	tfCommand.Stdout = os.Stdout
	tfCommand.Stderr = os.Stderr
	err := tfCommand.Run()
	if err != nil {
		return err
	}
	return nil
}

func tfInit(buildPath string, codePath string) error {
	return runTfCommand(buildPath, "init", codePath)
}

func tfDestroy(buildPath string, state string, codePath string) error {
	return runTfCommand(buildPath, "destroy", "-force", "-state="+state+".tfstate", codePath)
}

func tfApply(buildPath string, state string, codePath string) error {
	return runTfCommand(buildPath, "apply", "-state="+state+".tfstate", codePath)
}
